package dashwinddemo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestDashboardRendersDashWindTemplateSections(t *testing.T) {
	app := BuildApp()
	handler := app.Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{"DashWind", "fa-solid fa-house", "font-awesome/6.5.2/css/all.min.css", "New Users", "Total Sales", "User Signup Source", "Amount to be Collected", "Cash in hand", "Refresh Data", "Share", "Email Digests", "Download", "drawer lg:drawer-open dashwind-shell", "hx-post=\"/dashboard/period\""} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected dashboard to contain %q, got %q", want, body)
		}
	}
}

func TestPeriodActionReturnsDashboardFragment(t *testing.T) {
	app := BuildApp()
	handler := app.Handler()

	form := url.Values{"period": {"This quarter"}}
	req := httptest.NewRequest(http.MethodPost, "/dashboard/period", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{`id="dashwind-main"`, "Period updated to This quarter", "2026-04-01 ~ 2026-06-30"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected period response to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected fragment response, got full document")
	}
}

func TestLeadsAddAndDeleteActionsUpdateFragment(t *testing.T) {
	app := BuildApp()
	handler := app.Handler()

	addReq := httptest.NewRequest(http.MethodPost, "/leads/add", nil)
	addRR := httptest.NewRecorder()
	handler.ServeHTTP(addRR, addReq)
	if addRR.Code != http.StatusOK {
		t.Fatalf("expected add 200, got %d", addRR.Code)
	}
	addBody := addRR.Body.String()
	for _, want := range []string{"Demo Lead 6", "Added a demo lead", `id="dashwind-main"`} {
		if !strings.Contains(addBody, want) {
			t.Fatalf("expected add response to contain %q, got %q", want, addBody)
		}
	}

	form := url.Values{"email": {"demo6@example.com"}}
	deleteReq := httptest.NewRequest(http.MethodPost, "/leads/delete", strings.NewReader(form.Encode()))
	deleteReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	deleteRR := httptest.NewRecorder()
	handler.ServeHTTP(deleteRR, deleteReq)
	if deleteRR.Code != http.StatusOK {
		t.Fatalf("expected delete 200, got %d", deleteRR.Code)
	}
	deleteBody := deleteRR.Body.String()
	if strings.Contains(deleteBody, "Demo Lead 6") {
		t.Fatalf("expected deleted lead to be absent, got %q", deleteBody)
	}
	if !strings.Contains(deleteBody, "Lead removed") {
		t.Fatalf("expected delete notice, got %q", deleteBody)
	}
}

func TestTemplateRouteGroupsRenderDashWindPages(t *testing.T) {
	app := BuildApp()
	handler := app.Handler()

	cases := []struct {
		path string
		want []string
	}{
		{path: "/integration", want: []string{"Slack", "Salesforce", "toggle toggle-success toggle-lg", "Slack logo", `hx-post="/integration/toggle"`}},
		{path: "/settings-billing", want: []string{"Billing History", "#4567", "Product usage invoices"}},
		{path: "/login", want: []string{"Login", "Email Id", "Password", "DashWind user page preview"}},
		{path: "/components", want: []string{"Components", "dashwind.Shell", "ActionFormWithOptions", "Documentation", "Getting Started"}},
	}
	for _, tt := range cases {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rr.Code)
			}
			body := rr.Body.String()
			for _, want := range tt.want {
				if !strings.Contains(body, want) {
					t.Fatalf("expected %s to contain %q, got %q", tt.path, want, body)
				}
			}
		})
	}
}

func TestIntegrationToggleAndCalendarDayActionsUpdateFragments(t *testing.T) {
	app := BuildApp()
	handler := app.Handler()

	integrationForm := url.Values{"integration": {"Slack"}}
	integrationReq := httptest.NewRequest(http.MethodPost, "/integration/toggle", strings.NewReader(integrationForm.Encode()))
	integrationReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	integrationRR := httptest.NewRecorder()
	handler.ServeHTTP(integrationRR, integrationReq)
	if integrationRR.Code != http.StatusOK {
		t.Fatalf("expected integration toggle 200, got %d", integrationRR.Code)
	}
	integrationBody := integrationRR.Body.String()
	for _, want := range []string{`id="dashwind-main"`, "Slack disabled", "Slack logo", "Disabled"} {
		if !strings.Contains(integrationBody, want) {
			t.Fatalf("expected integration response to contain %q, got %q", want, integrationBody)
		}
	}

	calendarForm := url.Values{"day": {"20"}}
	calendarReq := httptest.NewRequest(http.MethodPost, "/calendar/day", strings.NewReader(calendarForm.Encode()))
	calendarReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	calendarRR := httptest.NewRecorder()
	handler.ServeHTTP(calendarRR, calendarReq)
	if calendarRR.Code != http.StatusOK {
		t.Fatalf("expected calendar day 200, got %d", calendarRR.Code)
	}
	calendarBody := calendarRR.Body.String()
	for _, want := range []string{`id="dashwind-main"`, "May 20 selected", "May 20 details", "Product webinar", "hx-post=\"/calendar/day\""} {
		if !strings.Contains(calendarBody, want) {
			t.Fatalf("expected calendar response to contain %q, got %q", want, calendarBody)
		}
	}
	if strings.Contains(calendarBody, "<!doctype html>") {
		t.Fatalf("expected fragment response, got full document")
	}
}
