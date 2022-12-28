// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/events (interfaces: StepRunner)

package mocks

import (
	pegomock "github.com/petergtz/pegomock"
	command "github.com/runatlantis/atlantis/server/events/command"
	"reflect"
	"time"
)

type MockStepRunner struct {
	fail func(message string, callerSkip ...int)
}

func NewMockStepRunner(options ...pegomock.Option) *MockStepRunner {
	mock := &MockStepRunner{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockStepRunner) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockStepRunner) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockStepRunner) Run(ctx command.ProjectContext, extraArgs []string, path string, envs map[string]string) (string, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockStepRunner().")
	}
	params := []pegomock.Param{ctx, extraArgs, path, envs}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Run", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
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

func (mock *MockStepRunner) VerifyWasCalledOnce() *VerifierMockStepRunner {
	return &VerifierMockStepRunner{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockStepRunner) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockStepRunner {
	return &VerifierMockStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockStepRunner) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockStepRunner {
	return &VerifierMockStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockStepRunner) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockStepRunner {
	return &VerifierMockStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockStepRunner struct {
	mock                   *MockStepRunner
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockStepRunner) Run(ctx command.ProjectContext, extraArgs []string, path string, envs map[string]string) *MockStepRunner_Run_OngoingVerification {
	params := []pegomock.Param{ctx, extraArgs, path, envs}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Run", params, verifier.timeout)
	return &MockStepRunner_Run_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockStepRunner_Run_OngoingVerification struct {
	mock              *MockStepRunner
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockStepRunner_Run_OngoingVerification) GetCapturedArguments() (command.ProjectContext, []string, string, map[string]string) {
	ctx, extraArgs, path, envs := c.GetAllCapturedArguments()
	return ctx[len(ctx)-1], extraArgs[len(extraArgs)-1], path[len(path)-1], envs[len(envs)-1]
}

func (c *MockStepRunner_Run_OngoingVerification) GetAllCapturedArguments() (_param0 []command.ProjectContext, _param1 [][]string, _param2 []string, _param3 []map[string]string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]command.ProjectContext, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(command.ProjectContext)
		}
		_param1 = make([][]string, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.([]string)
		}
		_param2 = make([]string, len(c.methodInvocations))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
		_param3 = make([]map[string]string, len(c.methodInvocations))
		for u, param := range params[3] {
			_param3[u] = param.(map[string]string)
		}
	}
	return
}
