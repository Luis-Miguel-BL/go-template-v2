package decorator

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
)

type MonitoringLeadRepository struct {
	repo lead.LeadRepository
}

func NewMonitoringLeadRepository(repo lead.LeadRepository) *MonitoringLeadRepository {
	return &MonitoringLeadRepository{
		repo: repo,
	}
}

func (m *MonitoringLeadRepository) Save(ctx context.Context, lead *model.Lead) error {
	ctx, span := observability.GetObservability().StartSpan(ctx, "repository.lead.save")
	defer span.End()

	err := m.repo.Save(ctx, lead)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func (m *MonitoringLeadRepository) GetByID(ctx context.Context, leadID string) (*model.Lead, error) {
	ctx, span := observability.GetObservability().StartSpan(ctx, "repository.lead.get_by_id")
	defer span.End()

	lead, err := m.repo.GetByID(ctx, leadID)
	if err != nil {
		span.RecordError(err)
	}
	return lead, err
}
