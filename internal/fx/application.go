package fx

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/usecase"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/persistence/repository"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/persistence/repository/decorator"
	"go.uber.org/fx"
)

var ApplicationModule = fx.Module("application",
	fx.Provide(
		// usecases
		usecase.NewCreateLead,

		// services
		service.NewAuthService,

		// repositories
		func(cfg *config.Config, obs observability.Observability, dispatcher *messaging.AggregateRootEventDispatcher, dynamoDBClient *aws.DynamoDBClient) (leadRepo lead.LeadRepository) {
			leadRepo = repository.NewInMemoryLeadRepository(dispatcher)
			if !cfg.App.InMemoryDB {
				leadRepo = repository.NewDynamoDBLeadRepository(cfg.AWS.DynamoDB.LeadTableName, dispatcher, dynamoDBClient)
			}
			return decorator.NewMonitoringLeadRepository(leadRepo, obs)
		},
	),
)
