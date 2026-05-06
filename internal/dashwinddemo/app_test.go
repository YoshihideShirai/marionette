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
	for _, want := range []string{"DashWind", "New Users", "Total Sales", "User Channels", "drawer lg:drawer-open dashwind-shell", "hx-post=\"/dashboard/period\""} {
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
	for _, want := range []string{`id="dashwind-main"`, "Period updated to This quarter", "btn btn-sm btn-primary"} {
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
		{path: "/integration", want: []string{"Slack", "Salesforce", "toggle toggle-success toggle-lg"}},
		{path: "/settings-billing", want: []string{"Billing History", "#4567", "Product usage invoices"}},
		{path: "/login", want: []string{"Login", "Email Id", "Password", "DashWind user page preview"}},
		{path: "/components", want: []string{"Components", "dashwind.Shell", "ActionFormWithOptions"}},
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
