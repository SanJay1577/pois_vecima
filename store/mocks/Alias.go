// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	models "pois/models"

	mock "github.com/stretchr/testify/mock"
)

// Alias is an autogenerated mock type for the Alias type
type Alias struct {
	mock.Mock
}

// CreateAlias provides a mock function with given fields: _a0
func (_m *Alias) CreateAlias(_a0 models.Alias) (models.Alias, error) {
	ret := _m.Called(_a0)

	var r0 models.Alias
	if rf, ok := ret.Get(0).(func(models.Alias) models.Alias); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(models.Alias)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Alias) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlias provides a mock function with given fields: _a0, _a1
func (_m *Alias) DeleteAlias(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAlias provides a mock function with given fields: _a0
func (_m *Alias) FindAlias(_a0 string) ([]models.Alias, error) {
	ret := _m.Called(_a0)

	var r0 []models.Alias
	if rf, ok := ret.Get(0).(func(string) []models.Alias); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Alias)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAlias interface {
	mock.TestingT
	Cleanup(func())
}

// NewAlias creates a new instance of Alias. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAlias(t mockConstructorTestingTNewAlias) *Alias {
	mock := &Alias{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
