package decorator

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
)

type MonitoringLeadRepository struct {
	repo      lead.LeadRepository
	telemetry telemetry.Telemetry
}

func NewMonitoringLeadRepository(repo lead.LeadRepository, telemetry telemetry.Telemetry) *MonitoringLeadRepository {
	return &MonitoringLeadRepository{
		repo:      repo,
		telemetry: telemetry,
	}
}

func (m *MonitoringLeadRepository) Save(ctx context.Context, lead *model.Lead) error {
	ctx, span := m.telemetry.StartSpan(ctx, "repository.lead.save")
	defer span.End()

	err := m.repo.Save(ctx, lead)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func (m *MonitoringLeadRepository) GetByID(ctx context.Context, leadID string) (*model.Lead, error) {
	ctx, span := m.telemetry.StartSpan(ctx, "repository.lead.get_by_id")
	defer span.End()

	lead, err := m.repo.GetByID(ctx, leadID)
	if err != nil {
		span.RecordError(err)
	}
	return lead, err
}
