package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootPageRendersTaskSampleShell(t *testing.T) {
	app := buildApp("Simple Tasks", "Marionette end-to-end sample")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Simple Tasks",
		"Marionette end-to-end sample",
		"No tasks yet",
		`hx-post="/tasks/create"`,
		`id="task-list"`,
		`name="name"`,
		"Add Task",
		"htmx",
		"daisyui",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected root page to contain %q, got %q", want, body)
		}
	}
}

func TestTaskCreateActionRendersTaskListFragment(t *testing.T) {
	app := buildApp("Simple Tasks", "Marionette end-to-end sample")

	req := httptest.NewRequest(http.MethodPost, "/tasks/create", strings.NewReader("name=Write+tests"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{"Write tests", "ID", "Name"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected task list fragment to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected action to return a fragment, got full document %q", body)
	}
}
