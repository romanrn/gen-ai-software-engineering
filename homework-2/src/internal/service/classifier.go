package service

import (
	"strings"
	"support-tickets/internal/domain"
)

type Classifier struct{}

func NewClassifier() *Classifier {
	return &Classifier{}
}

var categoryKeywords = map[domain.Category][]string{
	domain.CategoryAccountAccess: {
		"login", "password", "2fa", "access", "account", "locked out", "sign in",
	},
	domain.CategoryTechnicalIssue: {
		"error", "crash", "bug", "broken", "not working", "500", "exception", "fail",
	},
	domain.CategoryBillingQuestion: {
		"payment", "invoice", "refund", "charge", "billing", "subscription", "price",
	},
	domain.CategoryFeatureRequest: {
		"feature", "enhancement", "suggestion", "request", "would like", "add support",
	},
	domain.CategoryBugReport: {
		"reproduce", "steps to", "expected", "actual", "defect", "regression",
	},
}

var priorityKeywords = map[domain.Priority][]string{
	domain.PriorityUrgent: {
		"can't access", "critical", "production down", "security", "outage", "breach",
	},
	domain.PriorityHigh: {
		"important", "blocking", "asap", "urgent", "immediately",
	},
	domain.PriorityLow: {
		"minor", "cosmetic", "nice to have", "when possible",
	},
}

func (c *Classifier) Classify(subject, description string) *domain.Classification {
	text := strings.ToLower(subject + " " + description)

	// Find best matching category
	bestCategory := domain.CategoryOther
	bestCategoryScore := 0.0
	var foundKeywords []string

	for category, keywords := range categoryKeywords {
		matched := 0
		var categoryMatched []string
		for _, kw := range keywords {
			if strings.Contains(text, strings.ToLower(kw)) {
				matched++
				categoryMatched = append(categoryMatched, kw)
			}
		}
		score := float64(matched) / float64(len(keywords))
		if score > bestCategoryScore {
			bestCategoryScore = score
			bestCategory = category
			foundKeywords = categoryMatched
		}
	}

	// Find priority
	bestPriority := domain.PriorityMedium
	bestPriorityScore := 0.0

	for priority, keywords := range priorityKeywords {
		matched := 0
		for _, kw := range keywords {
			if strings.Contains(text, strings.ToLower(kw)) {
				matched++
			}
		}
		score := float64(matched) / float64(len(keywords))
		if score > bestPriorityScore {
			bestPriorityScore = score
			bestPriority = priority
		}
	}

	// Cap confidence at 1.0
	confidence := bestCategoryScore
	if confidence > 1.0 {
		confidence = 1.0
	}

	reasoning := "Matched " + string(bestCategory)
	if len(foundKeywords) > 0 {
		reasoning += " keywords"
	}

	return &domain.Classification{
		Category:   bestCategory,
		Priority:   bestPriority,
		Confidence: confidence,
		Reasoning:  reasoning,
		Keywords:   foundKeywords,
	}
}
