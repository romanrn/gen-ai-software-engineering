package service

import (
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/in"
	"support-tickets/pkg/importer"
	"testing"
)

func BenchmarkValidator_ValidInput(b *testing.B) {
	v := NewValidator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateCreateTicket(
			"CUST001",
			"test@example.com",
			"John Doe",
			"Test Subject",
			"This is a valid description with enough characters",
			"account_access",
			"high",
		)
	}
}

func BenchmarkValidator_InvalidInput(b *testing.B) {
	v := NewValidator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateCreateTicket("", "bad-email", "", "", "x", "invalid", "bad")
	}
}

func BenchmarkClassifier_Classify(b *testing.B) {
	c := NewClassifier()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Classify(
			"Cannot login to my account",
			"I cannot access my account. Getting password error. Critical production issue.",
		)
	}
}

func BenchmarkCSVImporter_Parse(b *testing.B) {
	imp := importer.NewCSVImporter()
	data := []byte("customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n" +
		"C1,a@example.com,Alice,Subject,Valid description with enough chars,account_access,high,web_form,Chrome,desktop\n" +
		"C2,b@example.com,Bob,Subject 2,Another valid description here,billing_question,low,email,Firefox,mobile\n" +
		"C3,c@example.com,Carol,Subject 3,Yet another valid description ok,technical_issue,urgent,api,Safari,desktop")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		imp.Parse(data)
	}
}

func BenchmarkTicketService_Create(b *testing.B) {
	repo := memory.New()
	svc := NewTicketService(repo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.Create(
			"CUST001",
			"test@example.com",
			"John Doe",
			"Test Subject",
			"This is a valid description with enough characters",
			domain.CategoryAccountAccess,
			domain.PriorityHigh,
			domain.TicketMetadata{},
		)
	}
}

func BenchmarkTicketService_List(b *testing.B) {
	repo := memory.New()
	svc := NewTicketService(repo)
	for i := 0; i < 100; i++ {
		svc.Create("C", "t@e.com", "N", "Subject", "Valid description with enough chars",
			domain.CategoryAccountAccess, domain.PriorityHigh, domain.TicketMetadata{})
	}
	req := &in.TicketListRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.List(req)
	}
}
