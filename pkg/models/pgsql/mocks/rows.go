// Code generated by mockery (devel). DO NOT EDIT.

package pgsql_mock

import (
	pgconn "github.com/jackc/pgconn"
	pgproto3 "github.com/jackc/pgproto3/v2"
)

// Close provides a mock function with given fields:
func (_m *MockDBConn) Close() {
	_m.Called()
}

// CommandTag provides a mock function with given fields:
func (_m *MockDBConn) CommandTag() pgconn.CommandTag {
	ret := _m.Called()

	var r0 pgconn.CommandTag
	if rf, ok := ret.Get(0).(func() pgconn.CommandTag); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgconn.CommandTag)
		}
	}

	return r0
}

// Err provides a mock function with given fields:
func (_m *MockDBConn) Err() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FieldDescriptions provides a mock function with given fields:
func (_m *MockDBConn) FieldDescriptions() []pgproto3.FieldDescription {
	ret := _m.Called()

	var r0 []pgproto3.FieldDescription
	if rf, ok := ret.Get(0).(func() []pgproto3.FieldDescription); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pgproto3.FieldDescription)
		}
	}

	return r0
}

// Next provides a mock function with given fields:
func (_m *MockDBConn) Next() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RawValues provides a mock function with given fields:
func (_m *MockDBConn) RawValues() [][]byte {
	ret := _m.Called()

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func() [][]byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	return r0
}

// Scan provides a mock function with given fields: dest
func (_m *MockDBConn) Scan(dest ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Values provides a mock function with given fields:
func (_m *MockDBConn) Values() ([]interface{}, error) {
	ret := _m.Called()

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func() []interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
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
