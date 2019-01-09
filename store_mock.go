package main

import "github.com/stretchr/testify/mock"

// mock store contains additional methods for inspection
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateBird(bird *Bird) error {
	/*
		when this method is called, `m.Called` records the call, and also returns the result that we pass to it.
	*/
	rets := m.Called(bird)
	return rets.Error(0)
}

func (m *MockStore) GetBirds() ([]*Bird, error) {
	rets := m.Called()
	// rets.Get() returns whatever we pass to it, we need to typecast it to the
	// type we expect, []*Bird
	return rets.Get(0).([]*Bird), rets.Error(1)
}

func InitMockStore() *MockStore {
	// initializes the store variable and assigns a new MockStore instance to
	// it
	s := new(MockStore)
	store = s
	return s
}
