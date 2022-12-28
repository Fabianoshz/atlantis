// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/core/runtime (interfaces: PullApprovedChecker)

package mocks

import (
	pegomock "github.com/petergtz/pegomock"
	models "github.com/runatlantis/atlantis/server/events/models"
	"reflect"
	"time"
)

type MockPullApprovedChecker struct {
	fail func(message string, callerSkip ...int)
}

func NewMockPullApprovedChecker(options ...pegomock.Option) *MockPullApprovedChecker {
	mock := &MockPullApprovedChecker{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockPullApprovedChecker) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockPullApprovedChecker) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockPullApprovedChecker) PullIsApproved(_param0 models.Repo, _param1 models.PullRequest) (models.ApprovalStatus, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockPullApprovedChecker().")
	}
	params := []pegomock.Param{_param0, _param1}
	result := pegomock.GetGenericMockFrom(mock).Invoke("PullIsApproved", params, []reflect.Type{reflect.TypeOf((*models.ApprovalStatus)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 models.ApprovalStatus
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(models.ApprovalStatus)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockPullApprovedChecker) VerifyWasCalledOnce() *VerifierMockPullApprovedChecker {
	return &VerifierMockPullApprovedChecker{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockPullApprovedChecker) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockPullApprovedChecker {
	return &VerifierMockPullApprovedChecker{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockPullApprovedChecker) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockPullApprovedChecker {
	return &VerifierMockPullApprovedChecker{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockPullApprovedChecker) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockPullApprovedChecker {
	return &VerifierMockPullApprovedChecker{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockPullApprovedChecker struct {
	mock                   *MockPullApprovedChecker
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockPullApprovedChecker) PullIsApproved(_param0 models.Repo, _param1 models.PullRequest) *MockPullApprovedChecker_PullIsApproved_OngoingVerification {
	params := []pegomock.Param{_param0, _param1}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "PullIsApproved", params, verifier.timeout)
	return &MockPullApprovedChecker_PullIsApproved_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockPullApprovedChecker_PullIsApproved_OngoingVerification struct {
	mock              *MockPullApprovedChecker
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockPullApprovedChecker_PullIsApproved_OngoingVerification) GetCapturedArguments() (models.Repo, models.PullRequest) {
	_param0, _param1 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1], _param1[len(_param1)-1]
}

func (c *MockPullApprovedChecker_PullIsApproved_OngoingVerification) GetAllCapturedArguments() (_param0 []models.Repo, _param1 []models.PullRequest) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]models.Repo, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(models.Repo)
		}
		_param1 = make([]models.PullRequest, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.(models.PullRequest)
		}
	}
	return
}
