package mocks

import "errors"

type SubscriptionMockError struct {
	SubscriptionMock
}

func (sm *SubscriptionMockError) Err() <-chan error {
	sm.errChan = make(chan error, 1)
	go func() { sm.errChan <- errors.New("error on subscribing") }()
	return sm.errChan
}

type SubscriptionMock struct {
	errChan chan error
}

func (sm *SubscriptionMock) Err() <-chan error {
	sm.errChan = make(chan error, 1)
	return sm.errChan
}

func (sm *SubscriptionMock) Unsubscribe() {
	close(sm.errChan)
}
