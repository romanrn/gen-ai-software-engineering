package service

import (
	"testing"
	"support-tickets/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassify_AccountAccess(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Cannot login to account",
		"I cannot login. Getting access denied error.",
	)

	assert.Equal(t, domain.CategoryAccountAccess, classification.Category)
	assert.True(t, len(classification.Keywords) > 0)
	assert.True(t, classification.Confidence > 0)
}

func TestClassify_TechnicalIssue(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Application crashes on startup",
		"The app crashes with error 500 every time.",
	)

	assert.Equal(t, domain.CategoryTechnicalIssue, classification.Category)
	assert.True(t, len(classification.Keywords) > 0)
}

func TestClassify_BillingQuestion(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Double billing charge",
		"I was charged twice for my subscription. Need refund.",
	)

	assert.Equal(t, domain.CategoryBillingQuestion, classification.Category)
	assert.True(t, len(classification.Keywords) > 0)
}

func TestClassify_FeatureRequest(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Feature request for dark mode",
		"Would like to have a dark mode feature in the app.",
	)

	assert.Equal(t, domain.CategoryFeatureRequest, classification.Category)
	assert.True(t, len(classification.Keywords) > 0)
}

func TestClassify_BugReport(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Form validation broken",
		"Steps to reproduce: click submit. Expected: error message. Actual: blank page.",
	)

	assert.Equal(t, domain.CategoryBugReport, classification.Category)
	assert.True(t, len(classification.Keywords) > 0)
}

func TestClassify_PriorityUrgent(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Critical production issue",
		"Production is down. Security breach detected. Immediate action required.",
	)

	assert.Equal(t, domain.PriorityUrgent, classification.Priority)
}

func TestClassify_PriorityHigh(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Important feature broken",
		"This is a blocking issue that needs urgent attention asap.",
	)

	assert.Equal(t, domain.PriorityHigh, classification.Priority)
}

func TestClassify_PriorityLow(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Minor cosmetic issue",
		"There's a minor cosmetic bug. Nice to have when possible.",
	)

	assert.Equal(t, domain.PriorityLow, classification.Priority)
}

func TestClassify_ConfidenceBounds(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Test subject",
		"Test description",
	)

	assert.True(t, classification.Confidence >= 0.0)
	assert.True(t, classification.Confidence <= 1.0)
}

func TestClassify_DefaultMediumPriority(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Some random issue",
		"This is a ticket with no priority keywords.",
	)

	assert.Equal(t, domain.PriorityMedium, classification.Priority)
}

func TestClassify_CaseSensitivity(t *testing.T) {
	c := NewClassifier()
	classification1 := c.Classify(
		"Cannot LOGIN to account",
		"LOGIN issues",
	)

	classification2 := c.Classify(
		"Cannot login to account",
		"login issues",
	)

	// Both should classify the same way
	assert.Equal(t, classification1.Category, classification2.Category)
}

func TestClassify_HasReasoning(t *testing.T) {
	c := NewClassifier()
	classification := c.Classify(
		"Login issue",
		"Cannot access account",
	)

	require.NotNil(t, classification)
	assert.NotEmpty(t, classification.Reasoning)
}
