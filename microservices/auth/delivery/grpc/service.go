package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authService struct {
	UC     auth.AuthUsecase
	logger logger.Logger
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService(usecase auth.AuthUsecase, logger logger.Logger) *authService {
	return &authService{UC: usecase, logger: logger, UnimplementedAuthServiceServer: authpb.UnimplementedAuthServiceServer{}}
}

func (s *authService) GetUserById(ctx context.Context, r *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	id := int(r.GetId())
	user, err := s.UC.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Warn("failed to find user")
		return nil, status.Errorf(codes.NotFound, "cannot find user by id: %v", err)
	}
	return &authpb.GetUserResponse{
		User: &authpb.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Image:     user.Image,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}
