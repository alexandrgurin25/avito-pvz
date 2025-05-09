// Code generated by MockGen. DO NOT EDIT.
// Source: product.go
//
// Generated by this command:
//
//	mockgen -source=product.go -destination=mocks/product_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	entity "avito-pvz/internal/entity"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockService) AddProduct(ctx context.Context, categoryId, pvzId string) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, categoryId, pvzId)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockServiceMockRecorder) AddProduct(ctx, categoryId, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockService)(nil).AddProduct), ctx, categoryId, pvzId)
}

// DeleteLastProduct mocks base method.
func (m *MockService) DeleteLastProduct(ctx context.Context, productId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, productId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockServiceMockRecorder) DeleteLastProduct(ctx, productId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockService)(nil).DeleteLastProduct), ctx, productId)
}
