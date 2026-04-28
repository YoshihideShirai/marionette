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
	pg := pagination{Page: 1, PerPage: 5, TotalPages: 1}
	loadingHTML, err := renderUsersTableBody(nil, true, "", pg).Render()
	if err != nil {
		t.Fatalf("loading state render failed: %v", err)
	}
	if !strings.Contains(string(loadingHTML), `aria-busy="true"`) {
		t.Fatalf("expected loading skeleton markup, got %q", loadingHTML)
	}

	emptyHTML, err := renderUsersTableBody(nil, false, "", pg).Render()
	if err != nil {
		t.Fatalf("empty state render failed: %v", err)
	}
	if !strings.Contains(string(emptyHTML), "No users yet") {
		t.Fatalf("expected empty state title, got %q", emptyHTML)
	}

	dataHTML, err := renderUsersTableBody([]user{{ID: 1, Name: "Aiko", Email: "aiko@example.com", Role: "Admin", StartDate: "2024-01-01"}}, false, "", pg).Render()
	if err != nil {
		t.Fatalf("data state render failed: %v", err)
	}
	if !strings.Contains(string(dataHTML), "<table") {
		t.Fatalf("expected data table, got %q", dataHTML)
	}
}

func TestRenderUsersTableBodySortsByQueryColumn(t *testing.T) {
	pg := pagination{Page: 1, PerPage: 5, TotalPages: 1}
	dataHTML, err := renderUsersTableBody([]user{
		{ID: 1, Name: "Ren", Email: "ren@example.com", Role: "Editor", StartDate: "2024-07-01"},
		{ID: 2, Name: "Aiko", Email: "aiko@example.com", Role: "Admin", StartDate: "2024-03-18"},
	}, false, "name", pg).Render()
	if err != nil {
		t.Fatalf("data state render failed: %v", err)
	}
	got := string(dataHTML)
	if !strings.Contains(got, `href="/?page=1&amp;per_page=5&amp;sort=name"`) {
		t.Fatalf("expected sortable header link, got %q", got)
	}
	if strings.Index(got, "Aiko") > strings.Index(got, "Ren") {
		t.Fatalf("expected name-sorted rows, got %q", got)
	}
}

func TestUsersTableSortLinksKeepPerPageQuery(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodGet, "/?page=2&per_page=2", nil)
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `href="/?page=2&amp;per_page=2&amp;sort=name"`) {
		t.Fatalf("expected sort link to keep current page/per_page query, got %q", body)
	}
}

func TestUsersPaginationNavigation(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodGet, "/?page=2&per_page=2", nil)
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Page 2 / 2") {
		t.Fatalf("expected page indicator, got %q", body)
	}
	if !strings.Contains(body, `href="/?page=1&amp;per_page=2"`) {
		t.Fatalf("expected prev link, got %q", body)
	}
	if strings.Contains(body, `href="/?page=3&amp;per_page=2"`) {
		t.Fatalf("did not expect next link on last page, got %q", body)
	}
	if !strings.Contains(body, "Mina Suzuki") {
		t.Fatalf("expected second-page user entry, got %q", body)
	}
	if strings.Contains(body, "Aiko Tanaka") {
		t.Fatalf("did not expect first-page user entry on page 2, got %q", body)
	}
}

func TestUsersPageIncludesThemeToggleButton(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "🌓 Theme") {
		t.Fatalf("expected theme toggle button label, got %q", body)
	}
	if !strings.Contains(body, "window.mrnToggleTheme") {
		t.Fatalf("expected theme toggle button onclick handler, got %q", body)
	}
}
