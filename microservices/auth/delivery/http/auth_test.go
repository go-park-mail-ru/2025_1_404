package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Тест на регистрацию
func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)
	t.Run("registration ok", func(t *testing.T) {
		req := domain.RegisterRequest{
			Email:     "email@mail.ru",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Password:  "GoodPassword123",
		}

		user := domain.User{
			ID:        0,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Password:  "hashedPassword",
			Image:     "",
		}

		mockUS.EXPECT().IsEmailTaken(gomock.Any(), req.Email).Return(false)
		mockUS.EXPECT().CreateUser(gomock.Any(), req.Email, req.Password, req.FirstName, req.LastName).Return(user, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		userHandlers.Register(response, request)

		res := response.Result()

		assert.Equal(t, http.StatusCreated, res.StatusCode)

		cookie := res.Cookies()[0]
		assert.Equal(t, "token", cookie.Name)
		assert.Equal(t, cookie.SameSite, http.SameSiteStrictMode)
		assert.True(t, cookie.HttpOnly)
		assert.False(t, cookie.Secure)

		user.Password = ""
		var responseBody domain.User
		err := json.NewDecoder(response.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, user, responseBody)
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("bad json"))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		userHandlers.Register(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("invalid request fields", func(t *testing.T) {
		req := domain.RegisterRequest{
			Email:     "badEmail",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Password:  "GoodPassword123",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		userHandlers.Register(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("email is taken", func(t *testing.T) {
		req := domain.RegisterRequest{
			Email:     "email@taken.ru",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Password:  "GoodPassword123",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().IsEmailTaken(gomock.Any(), req.Email).Return(true)

		userHandlers.Register(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("usecase CreateUser failed", func(t *testing.T) {
		req := domain.RegisterRequest{
			Email:     "email@mail.ru",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Password:  "GoodPassword123",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().IsEmailTaken(gomock.Any(), req.Email).Return(false)
		mockUS.EXPECT().CreateUser(gomock.Any(), req.Email, req.Password, req.FirstName, req.LastName).Return(domain.User{}, fmt.Errorf("create user fail"))

		userHandlers.Register(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)

	t.Run("login ok", func(t *testing.T) {
		req := domain.LoginRequest{
			Email:    "email@mail.ru",
			Password: "GoodPassword123",
		}

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user := domain.User{
			ID:        0,
			Email:     req.Email,
			Password:  string(hashPassword),
			FirstName: "Name",
			LastName:  "LastName",
			Image:     "",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(user, nil)

		userHandlers.Login(response, request)

		res := response.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		cookie := res.Cookies()[0]
		assert.Equal(t, "token", cookie.Name)
		assert.Equal(t, cookie.SameSite, http.SameSiteStrictMode)
		assert.True(t, cookie.HttpOnly)
		assert.False(t, cookie.Secure)

		user.Password = ""
		var responseBody domain.User
		err := json.NewDecoder(response.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, user, responseBody)
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("bad json"))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		userHandlers.Login(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("invalid request fields", func(t *testing.T) {
		req := domain.LoginRequest{
			Email:    "",
			Password: "pass",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		userHandlers.Login(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("usecase GetUserByEmail failed", func(t *testing.T) {
		req := domain.LoginRequest{
			Email:    "email@mail.ru",
			Password: "Password123",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(domain.User{}, fmt.Errorf("usecase GetUserByEmail failed"))

		userHandlers.Login(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})

	t.Run("mismatched passwords", func(t *testing.T) {
		req := domain.LoginRequest{
			Email:    "email@mail.ru",
			Password: "Password123",
		}

		user := domain.User{
			Email:    req.Email,
			Password: "differentHash",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(user, nil)

		userHandlers.Login(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})
}

func TestMeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)

	t.Run("me ok", func(t *testing.T) {
		user := domain.User{
			Email:     "email@mail.ru",
			FirstName: "Ivan",
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		requestWithCtx := request.WithContext(ctx)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)

		userHandlers.Me(response, requestWithCtx)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("userID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		response := httptest.NewRecorder()

		userHandlers.Me(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("usecase GetUserByID failed", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		requestWithCtx := request.WithContext(ctx)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByID(gomock.Any(), 1).Return(domain.User{}, fmt.Errorf("usecase GetUserByID failed"))

		userHandlers.Me(response, requestWithCtx)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})
}

func TestLogoutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)

	t.Run("logout ok", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		response := httptest.NewRecorder()

		userHandlers.Logout(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})
}

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)

	t.Run("update ok", func(t *testing.T) {
		req := domain.UpdateRequest{
			Email:     "newmail@mail.ru",
			FirstName: "NewName",
			LastName:  "NewLastName",
		}

		user := domain.User{
			ID:        1,
			Email:     "old@mail.ru",
			FirstName: "OldName",
			LastName:  "OldLastName",
		}

		updatedUser := domain.User{
			ID:        1,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/me", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		requestWithCtx := request.WithContext(ctx)
		response := httptest.NewRecorder()

		req.ID = 1
		mockUS.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)
		mockUS.EXPECT().IsEmailTaken(gomock.Any(), req.Email).Return(false)
		mockUS.EXPECT().UpdateUser(gomock.Any(), domain.UserFromUpdated(req)).Return(updatedUser, nil)

		userHandlers.Update(response, requestWithCtx)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("userID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		response := httptest.NewRecorder()

		userHandlers.Update(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/update", bytes.NewBufferString("bad json"))
		response := httptest.NewRecorder()

		userHandlers.Update(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	})

	t.Run("invalid request fields", func(t *testing.T) {
		req := domain.UpdateRequest{
			Email:     "BadEmail",
			FirstName: "NewName",
			LastName:  "NewLastName",
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/me", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		requestWithCtx := request.WithContext(ctx)
		response := httptest.NewRecorder()

		userHandlers.Update(response, requestWithCtx)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("email taken", func(t *testing.T) {
		req := domain.UpdateRequest{
			Email:     "new@taken.ru",
			FirstName: "NewName",
			LastName:  "NewLastName",
		}
		user := domain.User{
			ID:        1,
			Email:     "email@taken.ru",
			FirstName: "OldName",
			LastName:  "OldLastName",
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/me", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		requestWithCtx := request.WithContext(ctx)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)
		mockUS.EXPECT().IsEmailTaken(gomock.Any(), req.Email).Return(true)

		userHandlers.Update(response, requestWithCtx)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("usecase UpdateUser failed", func(t *testing.T) {
		req := domain.UpdateRequest{
			Email:     "email@mail.ru",
			FirstName: "NewName",
			LastName:  "NewLastName",
		}

		user := domain.User{
			ID:        1,
			Email:     "email@mail.ru",
			FirstName: "OldName",
			LastName:  "OldLastName",
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/me", bytes.NewBuffer(body)).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)
		req.ID = 1
		mockUS.EXPECT().UpdateUser(gomock.Any(), domain.UserFromUpdated(req)).Return(domain.User{}, fmt.Errorf("usecase UserUpdate failed"))

		userHandlers.Update(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}

func TestDeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockAuthUsecase(ctrl)
	cfg, _ := config.NewConfig()
	userHandlers := NewAuthHandler(mockUS, cfg)

	t.Run("DeleteImage ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		mockUS.EXPECT().DeleteImage(gomock.Any(), 1).Return(domain.User{
			Email: "email@mail.ru",
			Image: "",
		}, nil)

		userHandlers.DeleteImage(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("userId not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
		response := httptest.NewRecorder()

		userHandlers.DeleteImage(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})
	t.Run("usecase DeleteImage failed", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/auth/me", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		mockUS.EXPECT().DeleteImage(gomock.Any(), 1).Return(domain.User{}, fmt.Errorf("usecase DeleteImage failed"))

		userHandlers.DeleteImage(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}
