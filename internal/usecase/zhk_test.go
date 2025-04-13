package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	mockRepo "github.com/go-park-mail-ru/2025_1_404/internal/repository/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetZhkByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockRepository(ctrl)
	mockLogger, _ := logger.NewZapLogger()

	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	zhkID := int64(1)

	expectedZhk := domain.Zhk{
		ID:      zhkID,
		Name:    "ЖК Тестовый",
		Address: "ул. Тестовая, 1",
	}

	t.Run("GetZhkByID ok", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhkID).Return(expectedZhk, nil)

		result, err := zhkUsecase.GetZhkByID(ctx, zhkID)

		assert.NoError(t, err)
		assert.Equal(t, expectedZhk, result)
	})

	t.Run("zhk not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("zhk not found")

		mockRepo.EXPECT().GetZhkByID(ctx, zhkID).Return(domain.Zhk{}, expectedErr)

		result, err := zhkUsecase.GetZhkByID(ctx, zhkID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.Zhk{}, result)
	})
}
