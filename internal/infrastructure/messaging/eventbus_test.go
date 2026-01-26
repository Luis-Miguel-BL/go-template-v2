package messaging

import (
	"context"
	"sync"
	"testing"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry/mocks"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fakeEvent struct {
	*domain.EventBase
}

func (f fakeEvent) EventName() domain.EventName { return "lead.created" }
func (f fakeEvent) EventID() domain.EventID     { return "event-123" }

type handlerMock struct {
	mock.Mock
}

func (h *handlerMock) Handle(ctx context.Context, e domain.Event) {
	h.Called(ctx, e)
}

type subscriberMock struct {
	sync  eventbus.EventHandlersMap
	async eventbus.EventHandlersMap
}

func (s subscriberMock) SubscribedEvents() (
	eventbus.EventHandlersMap,
	eventbus.EventHandlersMap,
) {
	return s.sync, s.async
}

type DomainEventBusTestSuite struct {
	suite.Suite
	wg        *sync.WaitGroup
	telemetry *mocks.Telemetry
	handler   *handlerMock
	bus       eventbus.EventBus
}

func (s *DomainEventBusTestSuite) SetupTest() {
	s.wg = &sync.WaitGroup{}
	s.telemetry = new(mocks.Telemetry)
	s.handler = new(handlerMock)
}

func (s *DomainEventBusTestSuite) TearDownTest() {
	s.telemetry.AssertExpectations(s.T())
	s.handler.AssertExpectations(s.T())
}

func (s *DomainEventBusTestSuite) TestPublish() {

	s.Run("Should execute sync handler", func() {

		subscriber := subscriberMock{
			sync: map[domain.EventName]eventbus.EventHandler{
				"lead.created": s.handler.Handle,
			},
		}

		span := new(mocks.Span)
		span.On("SetAttributes", mock.Anything).Return()
		span.On("End").Return()
		s.telemetry.
			On("StartSpan", mock.Anything, "EventBus.Publish").
			Return(context.Background(), span).
			Once()

		s.handler.
			On("Handle", mock.Anything, mock.Anything).
			Once()

		s.bus = NewDomainEventBus(s.wg, s.telemetry, subscriber)

		s.bus.Publish(context.Background(), fakeEvent{})
		s.wg.Wait()
	})

	s.Run("Should execute async handler", func() {

		subscriber := subscriberMock{
			async: map[domain.EventName]eventbus.EventHandler{
				"lead.created": s.handler.Handle,
			},
		}

		span := new(mocks.Span)
		span.On("SetAttributes", mock.Anything).Return()
		span.On("End").Return()
		s.telemetry.
			On("StartSpan", mock.Anything, "EventBus.Publish").
			Return(context.Background(), span).
			Once()

		s.handler.
			On("Handle", mock.Anything, mock.Anything).
			Once()

		s.bus = NewDomainEventBus(s.wg, s.telemetry, subscriber)

		s.bus.Publish(context.Background(), fakeEvent{})
		s.wg.Wait()
	})

	s.Run("Should execute both sync and async handlers", func() {

		subscriber := subscriberMock{
			sync: map[domain.EventName]eventbus.EventHandler{
				"lead.created": s.handler.Handle,
			},
			async: map[domain.EventName]eventbus.EventHandler{
				"lead.created": s.handler.Handle,
			},
		}

		span := new(mocks.Span)
		span.On("SetAttributes", mock.Anything).Return()
		span.On("End").Return()
		s.telemetry.
			On("StartSpan", mock.Anything, "EventBus.Publish").
			Return(context.Background(), span).
			Twice()

		s.handler.
			On("Handle", mock.Anything, mock.Anything).
			Twice()

		s.bus = NewDomainEventBus(s.wg, s.telemetry, subscriber)

		s.bus.Publish(context.Background(), fakeEvent{})
		s.wg.Wait()
	})

	s.Run("Should not panic when no handlers exist", func() {

		s.bus = NewDomainEventBus(s.wg, s.telemetry)

		assert.NotPanics(s.T(), func() {
			s.bus.Publish(context.Background(), fakeEvent{})
		})
	})
}

func TestDomainEventBusTestSuite(t *testing.T) {
	suite.Run(t, new(DomainEventBusTestSuite))
}
