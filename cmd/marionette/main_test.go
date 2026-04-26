package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateUserStartDateValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr string
	}{
		{name: "missing date", date: "", wantErr: "Start date is required."},
		{name: "invalid date", date: "2026-02-30", wantErr: "Enter a valid date (YYYY-MM-DD)."},
		{name: "out of range", date: "2030-01-01", wantErr: "Start date must be between 2024-01-01 and 2026-12-31."},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := buildApp()
			form := url.Values{
				"name":       {"Demo User"},
				"email":      {"demo@example.com"},
				"role":       {"Viewer"},
				"start_date": {tc.date},
			}
			req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			app.Handler().ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rr.Code)
			}
			body := rr.Body.String()
			if !strings.Contains(body, `name="start_date"`) {
				t.Fatalf("expected start_date input in response, got %q", body)
			}
			if !strings.Contains(body, tc.wantErr) {
				t.Fatalf("expected error %q in response, got %q", tc.wantErr, body)
			}
		})
	}
}

func TestRenderUsersTableBodySwitchesBetweenLoadingEmptyAndData(t *testing.T) {
	loadingHTML, err := renderUsersTableBody(nil, true).Render()
	if err != nil {
		t.Fatalf("loading state render failed: %v", err)
	}
	if !strings.Contains(string(loadingHTML), `aria-busy="true"`) {
		t.Fatalf("expected loading skeleton markup, got %q", loadingHTML)
	}

	emptyHTML, err := renderUsersTableBody(nil, false).Render()
	if err != nil {
		t.Fatalf("empty state render failed: %v", err)
	}
	if !strings.Contains(string(emptyHTML), "No users yet") {
		t.Fatalf("expected empty state title, got %q", emptyHTML)
	}

	dataHTML, err := renderUsersTableBody([]user{{ID: 1, Name: "Aiko", Email: "aiko@example.com", Role: "Admin", StartDate: "2024-01-01"}}, false).Render()
	if err != nil {
		t.Fatalf("data state render failed: %v", err)
	}
	if !strings.Contains(string(dataHTML), "<table") {
		t.Fatalf("expected data table, got %q", dataHTML)
	}
}
