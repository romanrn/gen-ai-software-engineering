package service

import (
	"fmt"
	"strings"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/in"
	"support-tickets/internal/ports/out"
	"support-tickets/pkg/importer"
)

type TicketServiceImpl struct {
	repo       out.TicketRepository
	validator  *Validator
	classifier *Classifier
	importers  map[string]importer.TicketImporter
}

func NewTicketService(repo out.TicketRepository) *TicketServiceImpl {
	return &TicketServiceImpl{
		repo:       repo,
		validator:  NewValidator(),
		classifier: NewClassifier(),
		importers: map[string]importer.TicketImporter{
			"csv":  importer.NewCSVImporter(),
			"json": importer.NewJSONImporter(),
			"xml":  importer.NewXMLImporter(),
		},
	}
}

func (s *TicketServiceImpl) Create(
	customerID string,
	email string,
	name string,
	subject string,
	description string,
	category domain.Category,
	priority domain.Priority,
	metadata domain.TicketMetadata,
) (*domain.Ticket, error) {
	if err := s.validator.ValidateCreateTicket(customerID, email, name, subject, description, string(category), string(priority)); err != nil {
		return nil, err
	}

	ticket := domain.NewTicket(customerID, email, name, subject, description, category, priority, metadata)
	if err := s.repo.Save(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketServiceImpl) BulkImport(format string, data []byte) (*in.BulkImportResult, error) {
	imp, ok := s.importers[strings.ToLower(format)]
	if !ok {
		return nil, &domain.ErrInvalidInput{Message: fmt.Sprintf("unsupported format: %s", format)}
	}

	records, err := imp.Parse(data)
	if err != nil {
		return nil, err
	}

	result := &in.BulkImportResult{
		Total: len(records),
	}

	for i, record := range records {
		if err := s.validator.ValidateCreateTicket(
			record.CustomerID,
			record.CustomerEmail,
			record.CustomerName,
			record.Subject,
			record.Description,
			record.Category,
			record.Priority,
		); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, in.ImportError{
				Row:     i + 1,
				Message: err.Error(),
			})
			continue
		}

		ticket := domain.NewTicket(
			record.CustomerID,
			record.CustomerEmail,
			record.CustomerName,
			record.Subject,
			record.Description,
			domain.Category(record.Category),
			domain.Priority(record.Priority),
			domain.TicketMetadata{
				Source:     domain.Source(record.Metadata.Source),
				Browser:    record.Metadata.Browser,
				DeviceType: domain.DeviceType(record.Metadata.DeviceType),
			},
		)

		if err := s.repo.Save(ticket); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, in.ImportError{
				Row:     i + 1,
				Message: fmt.Sprintf("failed to save: %v", err),
			})
			continue
		}

		result.Successful++
	}

	return result, nil
}

func (s *TicketServiceImpl) List(req *in.TicketListRequest) ([]*domain.Ticket, error) {
	filter := &out.TicketFilter{
		Status:     req.Status,
		Category:   req.Category,
		Priority:   req.Priority,
		CustomerID: req.CustomerID,
	}

	return s.repo.FindAll(filter)
}

func (s *TicketServiceImpl) GetByID(id string) (*domain.Ticket, error) {
	return s.repo.FindByID(id)
}

func (s *TicketServiceImpl) Update(id string, updates map[string]interface{}) (*domain.Ticket, error) {
	return s.repo.Update(id, updates)
}

func (s *TicketServiceImpl) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *TicketServiceImpl) AutoClassify(id string) (*domain.Classification, error) {
	ticket, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, &domain.ErrNotFound{ID: id}
	}

	classification := s.classifier.Classify(ticket.Subject, ticket.Description)

	updates := map[string]interface{}{
		"classification": classification,
	}
	if _, err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return classification, nil
}
