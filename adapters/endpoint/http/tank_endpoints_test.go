package http

import (
	"archimedes-server/core/tank/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeReadTank struct {
	tanks      []domain.Tank
	tanksErr   error
	tankStatus *domain.TankStatus
	tankErr    error
}

func (f *fakeReadTank) GetTanks(ctx context.Context) ([]domain.Tank, error) {
	return f.tanks, f.tanksErr
}

func (f *fakeReadTank) GetTankByID(ctx context.Context, tankID string) (*domain.TankStatus, error) {
	return f.tankStatus, f.tankErr
}

func TestTankAPI_GetTanksHandler(t *testing.T) {
	tests := []struct {
		name       string
		repo       *fakeReadTank
		wantStatus int
		wantBody   string
	}{
		{
			name:       "returns tanks as json",
			repo:       &fakeReadTank{tanks: []domain.Tank{{ID: "1", Name: "Tank A"}}},
			wantStatus: http.StatusOK,
			wantBody:   `{"tanks":[{"id":"1","name":"Tank A"}]}`,
		},
		{
			name:       "empty tank list",
			repo:       &fakeReadTank{tanks: []domain.Tank{}},
			wantStatus: http.StatusOK,
			wantBody:   `{"tanks":[]}`,
		},
		{
			name:       "repository error returns 500",
			repo:       &fakeReadTank{tanksErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			api := &tankAPI{readTankRepository: tt.repo}
			req := httptest.NewRequest(http.MethodGet, "/read/tank", nil)
			rec := httptest.NewRecorder()

			api.GetTanksHandler(rec, req)

			is.Equal(tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				is.JSONEq(tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestTankAPI_GetTankByIDHandler(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		repo       *fakeReadTank
		wantStatus int
		wantBody   string
	}{
		{
			name: "returns tank status as json",
			repo: &fakeReadTank{tankStatus: &domain.TankStatus{
				Capacity:  100,
				Volume:    42.5,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}},
			wantStatus: http.StatusOK,
			wantBody: mustMarshal(t, &domain.TankStatus{
				Capacity:  100,
				Volume:    42.5,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}),
		},
		{
			name:       "not found returns 404",
			repo:       &fakeReadTank{tankStatus: nil},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "repository error returns 500",
			repo:       &fakeReadTank{tankErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			api := &tankAPI{readTankRepository: tt.repo}
			req := httptest.NewRequest(http.MethodGet, "/read/tank/tank-1", nil)
			req.SetPathValue("id", "tank-1")
			rec := httptest.NewRecorder()

			api.GetTankByIDHandler(rec, req)

			is.Equal(tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				is.JSONEq(tt.wantBody, rec.Body.String())
			}
		})
	}
}

func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal test fixture: %v", err)
	}
	return string(b)
}
