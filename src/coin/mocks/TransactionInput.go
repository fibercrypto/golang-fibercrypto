// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/fibercrypto/FiberCryptoWallet/src/core"
import mock "github.com/stretchr/testify/mock"

// TransactionInput is an autogenerated mock type for the TransactionInput type
type TransactionInput struct {
	mock.Mock
}

// Clone provides a mock function with given fields:
func (_m *TransactionInput) Clone() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCoins provides a mock function with given fields: ticker
func (_m *TransactionInput) GetCoins(ticker string) (uint64, error) {
	ret := _m.Called(ticker)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(string) uint64); ok {
		r0 = rf(ticker)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ticker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetId provides a mock function with given fields:
func (_m *TransactionInput) GetId() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSpentOutput provides a mock function with given fields:
func (_m *TransactionInput) GetSpentOutput() core.TransactionOutput {
	ret := _m.Called()

	var r0 core.TransactionOutput
	if rf, ok := ret.Get(0).(func() core.TransactionOutput); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TransactionOutput)
		}
	}

	return r0
}

// SupportedAssets provides a mock function with given fields:
func (_m *TransactionInput) SupportedAssets() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
