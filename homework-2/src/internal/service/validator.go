package service

import (
	"regexp"
	"support-tickets/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreateTicket(
	customerID string,
	email string,
	name string,
	subject string,
	description string,
	category string,
	priority string,
) *domain.ErrValidation {
	var details []domain.ValidationDetail

	if customerID == "" {
		details = append(details, domain.ValidationDetail{
			Field:   "customer_id",
			Message: "customer_id is required",
		})
	}

	if email == "" {
		details = append(details, domain.ValidationDetail{
			Field:   "customer_email",
			Message: "customer_email is required",
		})
	} else if !emailRegex.MatchString(email) {
		details = append(details, domain.ValidationDetail{
			Field:   "customer_email",
			Message: "customer_email is not a valid email address",
		})
	}

	if name == "" {
		details = append(details, domain.ValidationDetail{
			Field:   "customer_name",
			Message: "customer_name is required",
		})
	}

	if subject == "" {
		details = append(details, domain.ValidationDetail{
			Field:   "subject",
			Message: "subject is required",
		})
	} else if len(subject) < 1 || len(subject) > 200 {
		details = append(details, domain.ValidationDetail{
			Field:   "subject",
			Message: "subject must be between 1 and 200 characters",
		})
	}

	if description == "" {
		details = append(details, domain.ValidationDetail{
			Field:   "description",
			Message: "description is required",
		})
	} else if len(description) < 10 || len(description) > 2000 {
		details = append(details, domain.ValidationDetail{
			Field:   "description",
			Message: "description must be between 10 and 2000 characters",
		})
	}

	if !domain.IsValidCategory(category) {
		details = append(details, domain.ValidationDetail{
			Field:   "category",
			Message: "category must be one of: account_access, technical_issue, billing_question, feature_request, bug_report, other",
		})
	}

	if !domain.IsValidPriority(priority) {
		details = append(details, domain.ValidationDetail{
			Field:   "priority",
			Message: "priority must be one of: urgent, high, medium, low",
		})
	}

	if len(details) > 0 {
		return &domain.ErrValidation{Details: details}
	}

	return nil
}
