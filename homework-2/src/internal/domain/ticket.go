package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryAccountAccess  Category = "account_access"
	CategoryTechnicalIssue Category = "technical_issue"
	CategoryBillingQuestion Category = "billing_question"
	CategoryFeatureRequest  Category = "feature_request"
	CategoryBugReport       Category = "bug_report"
	CategoryOther           Category = "other"
)

type Priority string

const (
	PriorityUrgent Priority = "urgent"
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Status string

const (
	StatusNew            Status = "new"
	StatusInProgress     Status = "in_progress"
	StatusWaitingCustomer Status = "waiting_customer"
	StatusResolved       Status = "resolved"
	StatusClosed         Status = "closed"
)

type Source string

const (
	SourceWebForm Source = "web_form"
	SourceEmail   Source = "email"
	SourceAPI     Source = "api"
	SourceChat    Source = "chat"
	SourcePhone   Source = "phone"
)

type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeMobile  DeviceType = "mobile"
	DeviceTypeTablet  DeviceType = "tablet"
)

type TicketMetadata struct {
	Source     Source     `json:"source"`
	Browser    string     `json:"browser"`
	DeviceType DeviceType `json:"device_type"`
}

type Classification struct {
	Category   Category `json:"category"`
	Priority   Priority `json:"priority"`
	Confidence float64  `json:"confidence"`
	Reasoning  string   `json:"reasoning"`
	Keywords   []string `json:"keywords_found"`
}

type Ticket struct {
	ID                 uuid.UUID       `json:"id"`
	CustomerID         string          `json:"customer_id"`
	CustomerEmail      string          `json:"customer_email"`
	CustomerName       string          `json:"customer_name"`
	Subject            string          `json:"subject"`
	Description        string          `json:"description"`
	Category           Category        `json:"category"`
	Priority           Priority        `json:"priority"`
	Status             Status          `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	ResolvedAt         *time.Time      `json:"resolved_at"`
	AssignedTo         *string         `json:"assigned_to"`
	Tags               []string        `json:"tags"`
	Metadata           TicketMetadata  `json:"metadata"`
	Classification     *Classification `json:"classification,omitempty"`
}

func NewTicket(
	customerID string,
	customerEmail string,
	customerName string,
	subject string,
	description string,
	category Category,
	priority Priority,
	metadata TicketMetadata,
) *Ticket {
	now := time.Now().UTC()
	return &Ticket{
		ID:          uuid.New(),
		CustomerID:  customerID,
		CustomerEmail: customerEmail,
		CustomerName: customerName,
		Subject:     subject,
		Description: description,
		Category:    category,
		Priority:    priority,
		Status:      StatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        []string{},
		Metadata:    metadata,
	}
}

func IsValidCategory(c string) bool {
	switch Category(c) {
	case CategoryAccountAccess, CategoryTechnicalIssue, CategoryBillingQuestion,
		CategoryFeatureRequest, CategoryBugReport, CategoryOther:
		return true
	}
	return false
}

func IsValidPriority(p string) bool {
	switch Priority(p) {
	case PriorityUrgent, PriorityHigh, PriorityMedium, PriorityLow:
		return true
	}
	return false
}

func IsValidStatus(s string) bool {
	switch Status(s) {
	case StatusNew, StatusInProgress, StatusWaitingCustomer, StatusResolved, StatusClosed:
		return true
	}
	return false
}

func IsValidSource(s string) bool {
	switch Source(s) {
	case SourceWebForm, SourceEmail, SourceAPI, SourceChat, SourcePhone:
		return true
	}
	return false
}

func IsValidDeviceType(d string) bool {
	switch DeviceType(d) {
	case DeviceTypeDesktop, DeviceTypeMobile, DeviceTypeTablet:
		return true
	}
	return false
}
