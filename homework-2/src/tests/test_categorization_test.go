package tests

import (
	"support-tickets/internal/domain"
	"support-tickets/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func classifier() *service.Classifier {
	return service.NewClassifier()
}

// Test 1: account_access keywords detected
func TestCategorization_AccountAccess(t *testing.T) {
	c := classifier().Classify("Cannot login", "Password is not working. Account access denied.")
	assert.Equal(t, domain.CategoryAccountAccess, c.Category)
	assert.NotEmpty(t, c.Keywords)
}

// Test 2: technical_issue keywords detected
func TestCategorization_TechnicalIssue(t *testing.T) {
	c := classifier().Classify("App crashes", "Getting 500 error. The bug causes a crash.")
	assert.Equal(t, domain.CategoryTechnicalIssue, c.Category)
}

// Test 3: billing_question keywords detected
func TestCategorization_BillingQuestion(t *testing.T) {
	c := classifier().Classify("Invoice problem", "Charged twice. Need refund for duplicate payment.")
	assert.Equal(t, domain.CategoryBillingQuestion, c.Category)
}

// Test 4: feature_request keywords detected
func TestCategorization_FeatureRequest(t *testing.T) {
	c := classifier().Classify("Feature request", "Would like to add support for dark mode enhancement.")
	assert.Equal(t, domain.CategoryFeatureRequest, c.Category)
}

// Test 5: bug_report keywords detected
func TestCategorization_BugReport(t *testing.T) {
	c := classifier().Classify("Defect found", "Steps to reproduce: click save. Expected result, actual defect.")
	assert.Equal(t, domain.CategoryBugReport, c.Category)
}

// Test 6: priority urgent detected
func TestCategorization_PriorityUrgent(t *testing.T) {
	c := classifier().Classify("Critical outage", "Production is down. Security breach detected.")
	assert.Equal(t, domain.PriorityUrgent, c.Priority)
}

// Test 7: priority high detected
func TestCategorization_PriorityHigh(t *testing.T) {
	c := classifier().Classify("Blocking issue", "This is important and blocking our team asap.")
	assert.Equal(t, domain.PriorityHigh, c.Priority)
}

// Test 8: priority low detected
func TestCategorization_PriorityLow(t *testing.T) {
	c := classifier().Classify("Minor bug", "Cosmetic issue only, nice to have when possible.")
	assert.Equal(t, domain.PriorityLow, c.Priority)
}

// Test 9: defaults to medium when no priority keywords
func TestCategorization_PriorityMediumDefault(t *testing.T) {
	c := classifier().Classify("General question", "I have a question about the system.")
	assert.Equal(t, domain.PriorityMedium, c.Priority)
}

// Test 10: confidence score is in [0, 1] and reasoning is non-empty
func TestCategorization_ConfidenceAndReasoning(t *testing.T) {
	c := classifier().Classify("Login error", "Cannot access account due to password issue.")
	require.NotNil(t, c)
	assert.GreaterOrEqual(t, c.Confidence, 0.0)
	assert.LessOrEqual(t, c.Confidence, 1.0)
	assert.NotEmpty(t, c.Reasoning)
}
