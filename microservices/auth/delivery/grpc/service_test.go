package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockAuthUsecase struct {
	mock.Mock
}

func (m *mockAuthUsecase) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestAuthService_GetUserById_Success(t *testing.T) {
	mockUC := new(mockAuthUsecase)
	svc := service.NewAuthService(mockUC, logger.NewStub())

	expectedUser := domain.User{
		ID:        1,
		Email:     "user@mail.com",
		FirstName: "Ivan",
		LastName:  "Petrov",
		Image:     "avatar.png",
		CreatedAt: time.Now(),
	}

	mockUC.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)

	resp, err := svc.GetUserById(context.Background(), &authpb.GetUserRequest{Id: 1})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.User.Id)
	assert.Equal(t, "Ivan", resp.User.FirstName)
	mockUC.AssertExpectations(t)
}

func TestAuthService_GetUserById_NotFound(t *testing.T) {
	mockUC := new(mockAuthUsecase)
	svc := service.NewAuthService(mockUC, logger.NewStub())

	mockUC.On("GetUserByID", mock.Anything, 2).Return(domain.User{}, errors.New("not found"))

	resp, err := svc.GetUserById(context.Background(), &authpb.GetUserRequest{Id: 2})
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
	mockUC.AssertExpectations(t)
}
