package mock_rank

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	db "github.com/project-ippl-dev/tanding-api/internal/db"
	tools "github.com/project-ippl-dev/tanding-api/internal/tools"
	gomock "go.uber.org/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// Summary mocks base method.
func (m *MockUsecase) Summary(ctx context.Context, eventID uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Summary", ctx, eventID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Summary indicates an expected call of Summary.
func (mr *MockUsecaseMockRecorder) Summary(ctx, eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Summary", reflect.TypeOf((*MockUsecase)(nil).Summary), ctx, eventID)
}

// SetRewardCertificateName mocks base method.
func (m *MockUsecase) SetRewardCertificateName(rank int16, className string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRewardCertificateName", rank, className)
	ret0, _ := ret[0].(string)
	return ret0
}

// SetRewardCertificateName indicates an expected call of SetRewardCertificateName.
func (mr *MockUsecaseMockRecorder) SetRewardCertificateName(rank, className any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRewardCertificateName", reflect.TypeOf((*MockUsecase)(nil).SetRewardCertificateName), rank, className)
}

// SetRewardCertificateCommitteeName mocks base method.
func (m *MockUsecase) SetRewardCertificateCommitteeName(role db.EventRole) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRewardCertificateCommitteeName", role)
	ret0, _ := ret[0].(string)
	return ret0
}

// SetRewardCertificateCommitteeName indicates an expected call of SetRewardCertificateCommitteeName.
func (mr *MockUsecaseMockRecorder) SetRewardCertificateCommitteeName(role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRewardCertificateCommitteeName", reflect.TypeOf((*MockUsecase)(nil).SetRewardCertificateCommitteeName), role)
}

// SetRankPoint mocks base method.
func (m *MockUsecase) SetRankPoint(rank int16) int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRankPoint", rank)
	ret0, _ := ret[0].(int32)
	return ret0
}

// SetRankPoint indicates an expected call of SetRankPoint.
func (mr *MockUsecaseMockRecorder) SetRankPoint(rank any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRankPoint", reflect.TypeOf((*MockUsecase)(nil).SetRankPoint), rank)
}

// StoreCertificateByRegistrationID mocks base method.
func (m *MockUsecase) StoreCertificateByRegistrationID(ctx context.Context, tx interface{}, txQuery interface{}, arg interface{}) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCertificateByRegistrationID", ctx, tx, txQuery, arg)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreCertificateByRegistrationID indicates an expected call of StoreCertificateByRegistrationID.
func (mr *MockUsecaseMockRecorder) StoreCertificateByRegistrationID(ctx, tx, txQuery, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCertificateByRegistrationID", reflect.TypeOf((*MockUsecase)(nil).StoreCertificateByRegistrationID), ctx, tx, txQuery, arg)
}

// StoreCertificateExcludeRegistrationID mocks base method.
func (m *MockUsecase) StoreCertificateExcludeRegistrationID(ctx context.Context, tx interface{}, txQuery interface{}, arg interface{}) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCertificateExcludeRegistrationID", ctx, tx, txQuery, arg)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreCertificateExcludeRegistrationID indicates an expected call of StoreCertificateExcludeRegistrationID.
func (mr *MockUsecaseMockRecorder) StoreCertificateExcludeRegistrationID(ctx, tx, txQuery, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCertificateExcludeRegistrationID", reflect.TypeOf((*MockUsecase)(nil).StoreCertificateExcludeRegistrationID), ctx, tx, txQuery, arg)
}

// StoreCertificateEventCommittee mocks base method.
func (m *MockUsecase) StoreCertificateEventCommittee(ctx context.Context, tx interface{}, txQuery interface{}, eventID uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCertificateEventCommittee", ctx, tx, txQuery, eventID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreCertificateEventCommittee indicates an expected call of StoreCertificateEventCommittee.
func (mr *MockUsecaseMockRecorder) StoreCertificateEventCommittee(ctx, tx, txQuery, eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCertificateEventCommittee", reflect.TypeOf((*MockUsecase)(nil).StoreCertificateEventCommittee), ctx, tx, txQuery, eventID)
}

// StoreCertificateClubs mocks base method.
func (m *MockUsecase) StoreCertificateClubs(ctx context.Context, tx interface{}, txQuery interface{}, eventID uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCertificateClubs", ctx, tx, txQuery, eventID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreCertificateClubs indicates an expected call of StoreCertificateClubs.
func (mr *MockUsecaseMockRecorder) StoreCertificateClubs(ctx, tx, txQuery, eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCertificateClubs", reflect.TypeOf((*MockUsecase)(nil).StoreCertificateClubs), ctx, tx, txQuery, eventID)
}

// SetCertificateClubName mocks base method.
func (m *MockUsecase) SetCertificateClubName(rank int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCertificateClubName", rank)
	ret0, _ := ret[0].(string)
	return ret0
}

// SetCertificateClubName indicates an expected call of SetCertificateClubName.
func (mr *MockUsecaseMockRecorder) SetCertificateClubName(rank any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCertificateClubName", reflect.TypeOf((*MockUsecase)(nil).SetCertificateClubName), rank)
}

// FetchOwnPoint mocks base method.
func (m *MockUsecase) FetchOwnPoint(ctx context.Context, userID string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOwnPoint", ctx, userID)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOwnPoint indicates an expected call of FetchOwnPoint.
func (mr *MockUsecaseMockRecorder) FetchOwnPoint(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOwnPoint", reflect.TypeOf((*MockUsecase)(nil).FetchOwnPoint), ctx, userID)
}

// FetchByClubID mocks base method.
func (m *MockUsecase) FetchByClubID(ctx context.Context, clubID interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByClubID", ctx, clubID)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchByClubID indicates an expected call of FetchByClubID.
func (mr *MockUsecaseMockRecorder) FetchByClubID(ctx, clubID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByClubID", reflect.TypeOf((*MockUsecase)(nil).FetchByClubID), ctx, clubID)
}

// Rank mocks base method.
func (m *MockUsecase) Rank(ctx context.Context, page, pageSize int32, arg interface{}) (tools.Pagination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rank", ctx, page, pageSize, arg)
	ret0, _ := ret[0].(tools.Pagination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Rank indicates an expected call of Rank.
func (mr *MockUsecaseMockRecorder) Rank(ctx, page, pageSize, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rank", reflect.TypeOf((*MockUsecase)(nil).Rank), ctx, page, pageSize, arg)
}

// UserRank mocks base method.
func (m *MockUsecase) UserRank(ctx context.Context, page, pageSize int32, arg interface{}) (tools.Pagination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserRank", ctx, page, pageSize, arg)
	ret0, _ := ret[0].(tools.Pagination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserRank indicates an expected call of UserRank.
func (mr *MockUsecaseMockRecorder) UserRank(ctx, page, pageSize, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserRank", reflect.TypeOf((*MockUsecase)(nil).UserRank), ctx, page, pageSize, arg)
}
