package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestOrdersFilterFragmentKeepsMainContentTarget(t *testing.T) {
	app := buildApp()
	app.SetGlobal("loggedIn", true)
	app.SetGlobal("currentPage", "overview")
	handler := app.Handler()

	first := postFilter(t, handler, "Active")
	if strings.Contains(first, "<!doctype html>") {
		t.Fatalf("expected fragment response, got full document")
	}
	for _, want := range []string{`id="main-content"`, `Filter applied: Active`, `selected="selected" value="Active"`} {
		if !strings.Contains(first, want) {
			t.Fatalf("expected first filter response to contain %q, got %q", want, first)
		}
	}

	second := postFilter(t, handler, "Review")
	for _, want := range []string{`id="main-content"`, `Filter applied: Review`, `selected="selected" value="Review"`} {
		if !strings.Contains(second, want) {
			t.Fatalf("expected second filter response to contain %q, got %q", want, second)
		}
	}
}

func TestDashboardUsesOverlayDrawerNavigation(t *testing.T) {
	app := buildApp()
	app.SetGlobal("loggedIn", true)
	handler := app.Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{`class="drawer min-h-screen bg-base-200"`, `id="admin-nav-drawer"`, `drawer-side z-40`, `Open navigation drawer`, `Revenue Ops`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected dashboard to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "lg:drawer-open") {
		t.Fatalf("expected overlay drawer without permanently open sidebar, got %q", body)
	}
}

func TestDashboardLinksDealsToDetailPages(t *testing.T) {
	app := buildApp()
	app.SetGlobal("loggedIn", true)
	handler := app.Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{`href="/orders/detail?id=ORD-1042"`, `ORD-1042`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected dashboard to contain %q, got %q", want, body)
		}
	}
}

func TestOrderDetailPageShowsSelectedDeal(t *testing.T) {
	app := buildApp()
	app.SetGlobal("loggedIn", true)
	handler := app.Handler()

	req := httptest.NewRequest(http.MethodGet, "/orders/detail?id=ORD-1043", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{`Deal detail`, `ORD-1043 - Northwind Health`, `Workflow action`, `Projected ARR`, `$64,000`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected detail page to contain %q, got %q", want, body)
		}
	}
}

func TestLoginReplacesNarrowLoginContainerWithDashboard(t *testing.T) {
	app := buildApp()
	handler := app.Handler()

	form := url.Values{"provider": {"demo-sso"}}
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	for _, want := range []string{`id="app-body"`, `drawer min-h-screen bg-base-200`, `max-w-none`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected login response to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, `max-w-5xl`) {
		t.Fatalf("expected login response to replace the narrow login container, got %q", body)
	}
}

func postFilter(t *testing.T, handler http.Handler, status string) string {
	t.Helper()
	form := url.Values{"status": {status}}
	req := httptest.NewRequest(http.MethodPost, "/orders/filter", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	return rr.Body.String()
}
