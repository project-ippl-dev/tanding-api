package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tanding-api/internal/user"
	userMocks "tanding-api/internal/user/mocks"
	"tanding-api/internal/tools/pagination"
)

func TestUserUsecase_FindAll(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	paginationQuery := pagination.PaginationQuery{Page: 1, Limit: 10}

	expectedUsers := []user.User{
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
	}

	repoMock.On("FindAll", ctx, paginationQuery).Return(expectedUsers, nil)

	users, err := usecase.FindAll(ctx, paginationQuery)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_FindAll_Error(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	paginationQuery := pagination.PaginationQuery{Page: 1, Limit: 10}

	expectedError := errors.New("database error")
	repoMock.On("FindAll", ctx, paginationQuery).Return(nil, expectedError)

	users, err := usecase.FindAll(ctx, paginationQuery)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_FindByID(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	userID := int64(1)

	expectedUser := &user.User{ID: userID, Username: "testuser"}

	repoMock.On("FindByID", ctx, userID).Return(expectedUser, nil)

	foundUser, err := usecase.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, foundUser)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_FindByID_NotFound(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	userID := int64(1)

	repoMock.On("FindByID", ctx, userID).Return(nil, nil)

	foundUser, err := usecase.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.Nil(t, foundUser)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_FindByID_Error(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	userID := int64(1)

	expectedError := errors.New("database error")
	repoMock.On("FindByID", ctx, userID).Return(nil, expectedError)

	foundUser, err := usecase.FindByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Create(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	newUser := user.User{Username: "newuser", Email: "newuser@example.com"}

	expectedID := int64(1)
	repoMock.On("Create", ctx, newUser).Return(expectedID, nil)

	createdID, err := usecase.Create(ctx, newUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, createdID)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Create_Error(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	newUser := user.User{Username: "newuser", Email: "newuser@example.com"}

	expectedError := errors.New("database error")
	repoMock.On("Create", ctx, newUser).Return(int64(0), expectedError)

	createdID, err := usecase.Create(ctx, newUser)

	assert.Error(t, err)
	assert.Equal(t, int64(0), createdID)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Update(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	updatedUser := user.User{ID: 1, Username: "updateduser"}

	repoMock.On("Update", ctx, updatedUser).Return(nil)

	err := usecase.Update(ctx, updatedUser)

	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Update_Error(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	updatedUser := user.User{ID: 1, Username: "updateduser"}

	expectedError := errors.New("database error")
	repoMock.On("Update", ctx, updatedUser).Return(expectedError)

	err := usecase.Update(ctx, updatedUser)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Delete(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	userID := int64(1)

	repoMock.On("Delete", ctx, userID).Return(nil)

	err := usecase.Delete(ctx, userID)

	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestUserUsecase_Delete_Error(t *testing.T) {
	repoMock := new(userMocks.Repository)
	usecase := user.NewUsecase(repoMock)

	ctx := context.Background()
	userID := int64(1)

	expectedError := errors.New("database error")
	repoMock.On("Delete", ctx, userID).Return(expectedError)

	err := usecase.Delete(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}