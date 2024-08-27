package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface) {

}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "test", Payload: "payload"}
	suite.event2 = TestEvent{Name: "test1", Payload: "payload2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	req := suite.Require()

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 2)

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler2)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event2.GetName()], 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	req.ErrorIs(err, ErrHandlerAlreadyRegistered)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 2)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	req := suite.Require()

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	req.NoError(err)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 2)

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event2.GetName()], 1)

	suite.eventDispatcher.Clear()
	suite.Empty(suite.eventDispatcher.handlers)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	req := suite.Require()

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	req.NoError(err)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	req.NoError(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 2)

	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	suite.False(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)
	suite.eventDispatcher.Register(suite.event.GetName(), eh)
	suite.eventDispatcher.Dispatch(&suite.event)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
