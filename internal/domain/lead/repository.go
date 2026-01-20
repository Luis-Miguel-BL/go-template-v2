package lead

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
)

type LeadRepository interface {
	Save(ctx context.Context, lead *model.Lead) error
	GetByID(ctx context.Context, leadID string) (*model.Lead, error)
}
