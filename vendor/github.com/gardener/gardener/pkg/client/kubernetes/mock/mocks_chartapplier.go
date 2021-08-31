// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/gardener/pkg/client/kubernetes (interfaces: ChartApplier)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	chartrenderer "github.com/gardener/gardener/pkg/chartrenderer"
	kubernetes "github.com/gardener/gardener/pkg/client/kubernetes"
	gomock "github.com/golang/mock/gomock"
)

// MockChartApplier is a mock of ChartApplier interface.
type MockChartApplier struct {
	ctrl     *gomock.Controller
	recorder *MockChartApplierMockRecorder
}

// MockChartApplierMockRecorder is the mock recorder for MockChartApplier.
type MockChartApplierMockRecorder struct {
	mock *MockChartApplier
}

// NewMockChartApplier creates a new mock instance.
func NewMockChartApplier(ctrl *gomock.Controller) *MockChartApplier {
	mock := &MockChartApplier{ctrl: ctrl}
	mock.recorder = &MockChartApplierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChartApplier) EXPECT() *MockChartApplierMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockChartApplier) Apply(arg0 context.Context, arg1, arg2, arg3 string, arg4 ...kubernetes.ApplyOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Apply", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *MockChartApplierMockRecorder) Apply(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockChartApplier)(nil).Apply), varargs...)
}

// Delete mocks base method.
func (m *MockChartApplier) Delete(arg0 context.Context, arg1, arg2, arg3 string, arg4 ...kubernetes.DeleteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockChartApplierMockRecorder) Delete(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockChartApplier)(nil).Delete), varargs...)
}

// Render mocks base method.
func (m *MockChartApplier) Render(arg0, arg1, arg2 string, arg3 interface{}) (*chartrenderer.RenderedChart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Render", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*chartrenderer.RenderedChart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Render indicates an expected call of Render.
func (mr *MockChartApplierMockRecorder) Render(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*MockChartApplier)(nil).Render), arg0, arg1, arg2, arg3)
}

// RenderArchive mocks base method.
func (m *MockChartApplier) RenderArchive(arg0 []byte, arg1, arg2 string, arg3 interface{}) (*chartrenderer.RenderedChart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderArchive", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*chartrenderer.RenderedChart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RenderArchive indicates an expected call of RenderArchive.
func (mr *MockChartApplierMockRecorder) RenderArchive(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderArchive", reflect.TypeOf((*MockChartApplier)(nil).RenderArchive), arg0, arg1, arg2, arg3)
}
