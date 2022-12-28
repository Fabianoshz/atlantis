// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/core/runtime/models (interfaces: Exec)

package mocks

import (
	pegomock "github.com/petergtz/pegomock"
	"reflect"
	"time"
)

type MockExec struct {
	fail func(message string, callerSkip ...int)
}

func NewMockExec(options ...pegomock.Option) *MockExec {
	mock := &MockExec{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockExec) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockExec) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockExec) LookPath(file string) (string, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockExec().")
	}
	params := []pegomock.Param{file}
	result := pegomock.GetGenericMockFrom(mock).Invoke("LookPath", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockExec) CombinedOutput(args []string, envs map[string]string, workdir string) (string, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockExec().")
	}
	params := []pegomock.Param{args, envs, workdir}
	result := pegomock.GetGenericMockFrom(mock).Invoke("CombinedOutput", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockExec) VerifyWasCalledOnce() *VerifierMockExec {
	return &VerifierMockExec{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockExec) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockExec {
	return &VerifierMockExec{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockExec) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockExec {
	return &VerifierMockExec{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockExec) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockExec {
	return &VerifierMockExec{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockExec struct {
	mock                   *MockExec
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockExec) LookPath(file string) *MockExec_LookPath_OngoingVerification {
	params := []pegomock.Param{file}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "LookPath", params, verifier.timeout)
	return &MockExec_LookPath_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockExec_LookPath_OngoingVerification struct {
	mock              *MockExec
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockExec_LookPath_OngoingVerification) GetCapturedArguments() string {
	file := c.GetAllCapturedArguments()
	return file[len(file)-1]
}

func (c *MockExec_LookPath_OngoingVerification) GetAllCapturedArguments() (_param0 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
	}
	return
}

func (verifier *VerifierMockExec) CombinedOutput(args []string, envs map[string]string, workdir string) *MockExec_CombinedOutput_OngoingVerification {
	params := []pegomock.Param{args, envs, workdir}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "CombinedOutput", params, verifier.timeout)
	return &MockExec_CombinedOutput_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockExec_CombinedOutput_OngoingVerification struct {
	mock              *MockExec
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockExec_CombinedOutput_OngoingVerification) GetCapturedArguments() ([]string, map[string]string, string) {
	args, envs, workdir := c.GetAllCapturedArguments()
	return args[len(args)-1], envs[len(envs)-1], workdir[len(workdir)-1]
}

func (c *MockExec_CombinedOutput_OngoingVerification) GetAllCapturedArguments() (_param0 [][]string, _param1 []map[string]string, _param2 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([][]string, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.([]string)
		}
		_param1 = make([]map[string]string, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.(map[string]string)
		}
		_param2 = make([]string, len(c.methodInvocations))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
	}
	return
}
