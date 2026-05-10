package tests

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"support-tickets/internal/domain"
	"support-tickets/internal/service"
	"support-tickets/pkg/importer"
	"testing"
)

// Benchmark 1: validator with valid input
func BenchmarkPerf_ValidatorValid(b *testing.B) {
	v := service.NewValidator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateCreateTicket("CUST1", "user@example.com", "Name", "Subject",
			"This is a valid description with enough characters", "account_access", "high")
	}
}

// Benchmark 2: classifier on a realistic ticket text
func BenchmarkPerf_Classifier(b *testing.B) {
	c := service.NewClassifier()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Classify(
			"Cannot login — production down",
			"Critical security issue. Account locked out. Getting password error. Can't access anything.",
		)
	}
}

// Benchmark 3: CSV importer on fixture file
func BenchmarkPerf_CSVImporter(b *testing.B) {
	data, _ := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.csv"))
	imp := importer.NewCSVImporter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		imp.Parse(data)
	}
}

// Benchmark 4: JSON importer on fixture file
func BenchmarkPerf_JSONImporter(b *testing.B) {
	data, _ := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.json"))
	imp := importer.NewJSONImporter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		imp.Parse(data)
	}
}

// Benchmark 5: full HTTP round-trip — POST /tickets
func BenchmarkPerf_HTTPCreateTicket(b *testing.B) {
	app := newApp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := map[string]any{
			"customer_id": "C1", "customer_email": "u@e.com", "customer_name": "User",
			"subject":     "Bench subject",
			"description": "This is a benchmark description with enough chars.",
			"category":    string(domain.CategoryAccountAccess),
			"priority":    string(domain.PriorityHigh),
		}
		app.Test(jsonRequest(http.MethodPost, "/tickets", body))
	}
}

// Benchmark 6: parallel POST /tickets (concurrency stress)
func BenchmarkPerf_ConcurrentHTTPCreate(b *testing.B) {
	app := newApp()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			body := map[string]any{
				"customer_id":    fmt.Sprintf("C%d", i),
				"customer_email": fmt.Sprintf("u%d@bench.com", i),
				"customer_name":  "BenchUser",
				"subject":        "Parallel bench ticket",
				"description":    "Concurrent benchmark description with enough characters.",
				"category":       string(domain.CategoryTechnicalIssue),
				"priority":       string(domain.PriorityMedium),
			}
			app.Test(jsonRequest(http.MethodPost, "/tickets", body))
			i++
		}
	})
}

// Benchmark 7: GET /tickets with category+priority filter
func BenchmarkPerf_ListWithFilter(b *testing.B) {
	app := newApp()
	// pre-seed 50 tickets
	for i := 0; i < 50; i++ {
		body := map[string]any{
			"customer_id":    fmt.Sprintf("SEED%d", i),
			"customer_email": fmt.Sprintf("seed%d@bench.com", i),
			"customer_name":  "Seed User",
			"subject":        "Seed ticket",
			"description":    "Seeding tickets for list benchmark with enough characters.",
			"category":       string(domain.CategoryTechnicalIssue),
			"priority":       string(domain.PriorityHigh),
		}
		app.Test(jsonRequest(http.MethodPost, "/tickets", body))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(jsonRequest(http.MethodGet, "/tickets?category=technical_issue&priority=high", nil))
	}
}
