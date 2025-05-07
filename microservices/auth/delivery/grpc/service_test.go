package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthService_GetUserById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAuthUsecase(ctrl)
	svc := service.NewAuthService(mockUC, logger.NewStub())

	expectedUser := domain.User{
		ID:        1,
		Email:     "user@mail.com",
		FirstName: "Ivan",
		LastName:  "Petrov",
		Image:     "avatar.png",
		CreatedAt: time.Now(),
	}

	mockUC.EXPECT().GetUserByID(gomock.Any(), 1).Return(expectedUser, nil)

	resp, err := svc.GetUserById(context.Background(), &authpb.GetUserRequest{Id: 1})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.User.Id)
	assert.Equal(t, "Ivan", resp.User.FirstName)
	assert.Equal(t, "Petrov", resp.User.LastName)
	assert.Equal(t, "user@mail.com", resp.User.Email)
	assert.Equal(t, "avatar.png", resp.User.Image)
	assert.NotNil(t, resp.User.CreatedAt)
}

func TestAuthService_GetUserById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAuthUsecase(ctrl)
	svc := service.NewAuthService(mockUC, logger.NewStub())

	mockUC.EXPECT().GetUserByID(gomock.Any(), 2).Return(domain.User{}, errors.New("not found"))

	resp, err := svc.GetUserById(context.Background(), &authpb.GetUserRequest{Id: 2})
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
