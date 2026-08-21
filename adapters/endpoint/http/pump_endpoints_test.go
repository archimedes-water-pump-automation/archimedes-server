package http

import (
	"archimedes-server/core/pump/domain"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeReadPumpStatus struct {
	pumps      []domain.Pump
	pumpsErr   error
	pumpStatus *domain.PumpStatus
	pumpErr    error
	history    []domain.PumpStatus
	historyErr error
}

func (f *fakeReadPumpStatus) GetPumps(ctx context.Context) ([]domain.Pump, error) {
	return f.pumps, f.pumpsErr
}

func (f *fakeReadPumpStatus) GetPumpStatus(ctx context.Context, pumpID string) (*domain.PumpStatus, error) {
	return f.pumpStatus, f.pumpErr
}

func (f *fakeReadPumpStatus) GetPumpStatusHistory(ctx context.Context, pumpID string) ([]domain.PumpStatus, error) {
	return f.history, f.historyErr
}

func TestPumpAPI_GetPumpsHandler(t *testing.T) {
	tests := []struct {
		name       string
		repo       *fakeReadPumpStatus
		wantStatus int
		wantBody   string
	}{
		{
			name:       "returns pumps as json",
			repo:       &fakeReadPumpStatus{pumps: []domain.Pump{{ID: "1", Name: "Pump A"}}},
			wantStatus: http.StatusOK,
			wantBody:   `{"pumps":[{"id":"1","name":"Pump A"}]}`,
		},
		{
			name:       "repository error returns 500",
			repo:       &fakeReadPumpStatus{pumpsErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			api := &pumpAPI{readPumpStatusRepository: tt.repo}
			req := httptest.NewRequest(http.MethodGet, "/read/pump", nil)
			rec := httptest.NewRecorder()

			api.GetPumpsHandler(rec, req)

			is.Equal(tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				is.JSONEq(tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestPumpAPI_GetPumpByIDHandler(t *testing.T) {
	startedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		repo       *fakeReadPumpStatus
		wantStatus int
		wantBody   string
	}{
		{
			name: "returns pump status as json",
			repo: &fakeReadPumpStatus{pumpStatus: &domain.PumpStatus{
				ID:        "s1",
				PumpID:    "pump-1",
				StartedAt: startedAt,
			}},
			wantStatus: http.StatusOK,
			wantBody:   mustMarshal(t, &domain.PumpStatus{ID: "s1", PumpID: "pump-1", StartedAt: startedAt}),
		},
		{
			name:       "not found returns 404",
			repo:       &fakeReadPumpStatus{pumpStatus: nil},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "repository error returns 500",
			repo:       &fakeReadPumpStatus{pumpErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			api := &pumpAPI{readPumpStatusRepository: tt.repo}
			req := httptest.NewRequest(http.MethodGet, "/read/pump/pump-1", nil)
			req.SetPathValue("id", "pump-1")
			rec := httptest.NewRecorder()

			api.GetPumpByIDHandler(rec, req)

			is.Equal(tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				is.JSONEq(tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestPumpAPI_GetPumpHistoricHandler(t *testing.T) {
	startedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		repo       *fakeReadPumpStatus
		wantStatus int
		wantBody   string
	}{
		{
			name:       "returns pump history as json",
			repo:       &fakeReadPumpStatus{history: []domain.PumpStatus{{ID: "s1", PumpID: "pump-1", StartedAt: startedAt}}},
			wantStatus: http.StatusOK,
			wantBody:   mustMarshal(t, []domain.PumpStatus{{ID: "s1", PumpID: "pump-1", StartedAt: startedAt}}),
		},
		{
			name:       "nil history returns 404",
			repo:       &fakeReadPumpStatus{history: nil},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "repository error returns 500",
			repo:       &fakeReadPumpStatus{historyErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			api := &pumpAPI{readPumpStatusRepository: tt.repo}
			req := httptest.NewRequest(http.MethodGet, "/read/pump/pump-1/historic", nil)
			req.SetPathValue("id", "pump-1")
			rec := httptest.NewRecorder()

			api.GetPumpHistoricHandler(rec, req)

			is.Equal(tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				is.JSONEq(tt.wantBody, rec.Body.String())
			}
		})
	}
}
