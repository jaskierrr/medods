package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/internal/controller"
	"main/internal/handlers"
	"main/internal/lib/logger"
	"main/internal/models"
	service "main/internal/service/token"
	"main/test/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		tokenRepo *mock.MockRepositoryToken
		emailRepo *mock.MockRepositoryEmail
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger()

	tokenRepo := mock.NewMockRepositoryToken(ctrl)
	emailRepo := mock.NewMockRepositoryEmail(ctrl)

	testFields := &fields{
		tokenRepo: tokenRepo,
		emailRepo: emailRepo,
	}

	secret := "SECRET"
	accessTokenTTL := 1
	refreshTokenTTL := 30

	s := service.New(tokenRepo, emailRepo, logger, secret, accessTokenTTL, refreshTokenTTL)
	c := controller.New(s, logger)
	h := handlers.New(c, logger)

	user := models.User{
		ID: "1",
		IP: "123",
	}

	validRequest := httptest.NewRequest(http.MethodPost, "/login?id="+user.ID, nil)
	validRequest.Header.Set("X-Real-Ip", user.IP)

	invalidRequest := httptest.NewRequest(http.MethodPost, "/login?id=", nil)

	w := httptest.NewRecorder()
	ctx := context.Background()

	tests := []struct {
		name    string
		args    *http.Request
		prepare func(f *fields)
		wantID  string
		wantIP  string
		wantJSONErr error
		wantRefreshTokenErr error
		tokenExpFlag bool
	}{
		{
			name: "valid",
			args: validRequest,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.tokenRepo.EXPECT().Login(ctx, user, gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			wantID:  user.ID,
			wantIP:  user.IP,
			wantJSONErr: nil,
			wantRefreshTokenErr: nil,
			tokenExpFlag: false,
		},
		{
			name: "bad request",
			args: invalidRequest,
			prepare: func(f *fields) {
				gomock.InOrder(
					// f.tokenRepo.EXPECT().Login(ctx, user, gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			wantID:  "",
			wantIP:  "",
			wantJSONErr: io.EOF,
			wantRefreshTokenErr: bcrypt.ErrHashTooShort,
			tokenExpFlag: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			h.Login(w, tt.args)
			resp := w.Result()
			defer resp.Body.Close()
			response := models.Response{}
			jsonDecodeErr := json.NewDecoder(resp.Body).Decode(&response)

			// access token validation
			idAccessToken, ipAccessToken, expAccessToken, _ := validateToken(response.Access, secret)
			expAccessTokenTime := time.Unix(expAccessToken, 0)
			t.Log("Login")

			switch {
			case !errors.Is(jsonDecodeErr, tt.wantJSONErr):
				t.Errorf("\nJSON decode error = %v\nwant = %v", jsonDecodeErr, tt.wantJSONErr)
			case !reflect.DeepEqual(idAccessToken, tt.wantID):
				t.Errorf("\nAccess token: \nid = %v\nwant = %v", idAccessToken, tt.wantID)
			case !reflect.DeepEqual(ipAccessToken, tt.wantIP):
				t.Errorf("\nAccess token: \nip = %v\nwant = %v", ipAccessToken, tt.wantIP)
			case expAccessTokenTime.Before(time.Now()) != tt.tokenExpFlag:
				t.Errorf("\ntoken is expired")
			}

			// refresh token validation
			refreshTokenHashFromResponse, _ := base64.StdEncoding.DecodeString(response.Refresh)
			notHashedRefreshToken := tt.wantID+service.RefreshTokenSeparator+tt.wantIP
			validateRefreshTokenErr := bcrypt.CompareHashAndPassword(refreshTokenHashFromResponse, []byte(notHashedRefreshToken))
			if !errors.Is(validateRefreshTokenErr, tt.wantRefreshTokenErr) {
				t.Errorf("\nRefresh token error = %v\nwant = %v", validateRefreshTokenErr, tt.wantRefreshTokenErr)
				t.Errorf("\nRefresh token from response = %v", string(refreshTokenHashFromResponse))
				t.Errorf("\nNew Refresh token = %v", notHashedRefreshToken)
			}
		})
	}

}

func validateToken(tokenStr string, secret string) (string, string, int64, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return "", "", 0, err
	}
	if !token.Valid {
		return "", "", 0, fmt.Errorf("token invalid")
	}

	id, ok := claims["guid"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("id not found")
	}

	ip, ok := claims["ip"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("ip not found")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", "", 0, fmt.Errorf("exp not found")
	}

	return id, ip, int64(exp), nil
}
