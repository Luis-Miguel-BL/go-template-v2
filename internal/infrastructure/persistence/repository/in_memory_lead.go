package repository

import (
	"context"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging"
)

type InMemoryLeadRepository struct {
	mu         *sync.RWMutex
	dispatcher *messaging.AggregateRootEventDispatcher
	leads      map[string]*model.Lead
}

func NewInMemoryLeadRepository(dispatcher *messaging.AggregateRootEventDispatcher) *InMemoryLeadRepository {
	return &InMemoryLeadRepository{
		mu:         &sync.RWMutex{},
		dispatcher: dispatcher,
		leads:      make(map[string]*model.Lead),
	}
}

func (r *InMemoryLeadRepository) Save(ctx context.Context, lead *model.Lead) (err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.leads[string(lead.LeadUUID)] = lead
	r.dispatcher.PublishUncommitedEvents(ctx, lead)
	return nil
}

func (r *InMemoryLeadRepository) GetByID(ctx context.Context, leadID string) (lead *model.Lead, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lead = r.leads[leadID]
	return lead, nil
}
