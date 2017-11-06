// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/ory/fosite (interfaces: AuthorizeEndpointHandler)

// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	context "context"

	gomock "github.com/golang/mock/gomock"
	fosite "github.com/ory/fosite"
)

// Mock of AuthorizeEndpointHandler interface
type MockAuthorizeEndpointHandler struct {
	ctrl     *gomock.Controller
	recorder *_MockAuthorizeEndpointHandlerRecorder
}

// Recorder for MockAuthorizeEndpointHandler (not exported)
type _MockAuthorizeEndpointHandlerRecorder struct {
	mock *MockAuthorizeEndpointHandler
}

func NewMockAuthorizeEndpointHandler(ctrl *gomock.Controller) *MockAuthorizeEndpointHandler {
	mock := &MockAuthorizeEndpointHandler{ctrl: ctrl}
	mock.recorder = &_MockAuthorizeEndpointHandlerRecorder{mock}
	return mock
}

func (_m *MockAuthorizeEndpointHandler) EXPECT() *_MockAuthorizeEndpointHandlerRecorder {
	return _m.recorder
}

func (_m *MockAuthorizeEndpointHandler) HandleAuthorizeEndpointRequest(_param0 context.Context, _param1 fosite.AuthorizeRequester, _param2 fosite.AuthorizeResponder) error {
	ret := _m.ctrl.Call(_m, "HandleAuthorizeEndpointRequest", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAuthorizeEndpointHandlerRecorder) HandleAuthorizeEndpointRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HandleAuthorizeEndpointRequest", arg0, arg1, arg2)
}
