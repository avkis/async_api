package apiserver

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiResponse[T any] struct {
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r SignupRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *ApiServer) signupHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[SignupRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, fmt.Errorf("invalid request %w", err))
		}

		existingUser, err := s.store.Users.ByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if existingUser != nil {
			return NewErrWithStatus(http.StatusConflict, fmt.Errorf("email already registerd"))
		}

		_, err = s.store.Users.CreateUser(r.Context(), req.Email, req.Password)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := encode(ApiResponse[struct{}]{
			Message: "successfully signed up user",
		}, http.StatusCreated, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		return nil
	})
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r SigninRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *ApiServer) signinHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[SigninRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, fmt.Errorf("invalid request %w", err))
		}

		user, err := s.store.Users.ByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := user.ComparePassword(req.Password); err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokenStore.DeleteUserTokens(r.Context(), user.ID)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokenStore.Create(r.Context(), user.ID, tokenPair.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := encode(ApiResponse[SigninResponse]{
			Message: "successfully signed in user",
			Data: &SigninResponse{
				AccessToken:  tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.RefreshToken.Raw,
			},
		}, int(http.StatusOK), w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		return nil
	})
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r TokenRefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh token is required")
	}
	return nil
}

func (s *ApiServer) tokenRefreshHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[TokenRefreshRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, fmt.Errorf("invalid request %w", err))
		}

		currentRefreshToken, err := s.jwtManager.Parse(req.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		userIDStr, err := currentRefreshToken.Claims.GetSubject()
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		currentRefreshTokenRecord, err := s.store.RefreshTokenStore.ByPrimaryKey(r.Context(), userID, currentRefreshToken)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, sql.ErrNoRows) {
				status = http.StatusUnauthorized
			}
			return NewErrWithStatus(status, err)
		}

		if currentRefreshTokenRecord.ExpiresAt.Before(time.Now()) {
			return NewErrWithStatus(http.StatusUnauthorized, fmt.Errorf("refresh token expired"))
		}

		tokenPair, err := s.jwtManager.GenerateTokenPair(userID)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if _, err := s.store.RefreshTokenStore.DeleteUserTokens(r.Context(), userID); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if _, err := s.store.RefreshTokenStore.Create(r.Context(), userID, tokenPair.RefreshToken); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := encode(ApiResponse[TokenRefreshResponse]{
			Data: &TokenRefreshResponse{
				AccessToken:  tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.RefreshToken.Raw,
			},
		}, int(http.StatusOK), w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		return nil
	})
}
