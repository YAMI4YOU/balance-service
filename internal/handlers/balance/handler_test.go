package balance_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/YAMI4YOU/balance-service/internal/handlers/balance"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

type response struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		prepare      func(*Mockrepo)
		expectations func(t assert.TestingT, body string, status int)
	}{
		{
			name:   "success case",
			method: http.MethodGet,
			url:    "/balance?user_id=123",
			prepare: func(mockRepo *Mockrepo) {
				mockRepo.EXPECT().
					FetchBalance(gomock.Any(), 123).
					Return(models.NewMoneyFromKopecks(25075), nil)
			},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusOK, status)
				var resp response
				err := json.Unmarshal([]byte(body), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 123, resp.UserID)
				assert.Equal(t, 250.75, resp.Balance)
			},
		},
		{
			name:    "method not allowed",
			method:  http.MethodPost,
			url:     "/balance?user_id=123",
			prepare: func(mockRepo *Mockrepo) {},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusMethodNotAllowed, status)
				assert.Equal(t, "Method not allowed\n", body)
			},
		},
		{
			name:    "missing user_id parameter",
			method:  http.MethodGet,
			url:     "/balance",
			prepare: func(mockRepo *Mockrepo) {},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
				assert.Equal(t, "Missing user_id parameter\n", body)
			},
		},
		{
			name:    "invalid user_id parameter",
			method:  http.MethodGet,
			url:     "/balance?user_id=invalid",
			prepare: func(mockRepo *Mockrepo) {},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
				assert.Equal(t, "Invalid user_id parameter\n", body)
			},
		},
		{
			name:   "user not found",
			method: http.MethodGet,
			url:    "/balance?user_id=999",
			prepare: func(mockRepo *Mockrepo) {
				mockRepo.EXPECT().
					FetchBalance(gomock.Any(), 999).
					Return(models.Money{}, models.ErrUserNotFound)
			},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusNotFound, status)
				assert.Equal(t, "User not found\n", body)
			},
		},
		{
			name:   "internal server error",
			method: http.MethodGet,
			url:    "/balance?user_id=555",
			prepare: func(mockRepo *Mockrepo) {
				mockRepo.EXPECT().
					FetchBalance(gomock.Any(), 555).
					Return(models.Money{}, fmt.Errorf("database connection failed"))
			},
			expectations: func(t assert.TestingT, body string, status int) {
				assert.Equal(t, http.StatusInternalServerError, status)
				assert.Equal(t, "database connection failed\n", body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockrepo(ctrl)

			if tt.prepare != nil {
				tt.prepare(mockRepo)
			}

			handler := balance.NewHandler(mockRepo)
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			handler.Handle(rr, req)

			tt.expectations(t, rr.Body.String(), rr.Code)
		})
	}
}
