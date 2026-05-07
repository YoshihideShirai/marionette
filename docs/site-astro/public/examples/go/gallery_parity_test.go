package goexamples

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	mb "github.com/YoshihideShirai/marionette/backend"
)

func TestGallerySamplesMatchGoExamples(t *testing.T) {
	examples := []struct {
		id       string
		sample   string
		register func(*mb.App)
	}{
		{id: "actions", sample: "actions.html", register: RegisterActionsExample},
		{id: "alert", sample: "alert.html", register: RegisterAlertExample},
		{id: "app-shell", sample: "app-shell.html", register: RegisterAppShellExample},
		{id: "avatar", sample: "avatar.html", register: RegisterAvatarExample},
		{id: "badge", sample: "badge.html", register: RegisterBadgeExample},
		{id: "box", sample: "box.html", register: RegisterBoxExample},
		{id: "breadcrumb", sample: "breadcrumb.html", register: RegisterBreadcrumbExample},
		{id: "button", sample: "button.html", register: RegisterButtonExample},
		{id: "card", sample: "card.html", register: RegisterCardExample},
		{id: "carousel-item", sample: "carousel-item.html", register: RegisterCarouselItemExample},
		{id: "carousel", sample: "carousel.html", register: RegisterCarouselExample},
		{id: "chart-bar", sample: "chart-bar.html", register: RegisterChartBarExample},
		{id: "chart-doughnut", sample: "chart-doughnut.html", register: RegisterChartDoughnutExample},
		{id: "chart-line", sample: "chart-line.html", register: RegisterChartLineExample},
		{id: "chart-pie", sample: "chart-pie.html", register: RegisterChartPieExample},
		{id: "chart-scatter", sample: "chart-scatter.html", register: RegisterChartScatterExample},
		{id: "chart", sample: "chart.html", register: RegisterChartExample},
		{id: "chat-bubble", sample: "chat-bubble.html", register: RegisterChatBubbleExample},
		{id: "checkbox", sample: "checkbox.html", register: RegisterCheckboxExample},
		{id: "code", sample: "code.html", register: RegisterCodeExample},
		{id: "collapse", sample: "collapse.html", register: RegisterCollapseExample},
		{id: "container", sample: "container.html", register: RegisterContainerExample},
		{id: "dataframe-chart", sample: "dataframe-chart.html", register: RegisterDataFrameChartExample},
		{id: "dataframe", sample: "dataframe.html", register: RegisterDataFrameExample},
		{id: "dashwind-auth-card", sample: "dashwind-auth-card.html", register: RegisterDashwindAuthCardExample},
		{id: "dashwind-card-panel", sample: "dashwind-card-panel.html", register: RegisterDashwindCardPanelExample},
		{id: "dashwind-data-table", sample: "dashwind-data-table.html", register: RegisterDashwindDataTableExample},
		{id: "dashwind-metric-grid", sample: "dashwind-metric-grid.html", register: RegisterDashwindMetricGridExample},
		{id: "dashwind-page-header", sample: "dashwind-page-header.html", register: RegisterDashwindPageHeaderExample},
		{id: "dashwind-resource-page", sample: "dashwind-resource-page.html", register: RegisterDashwindResourcePageExample},
		{id: "dashwind-settings-section", sample: "dashwind-settings-section.html", register: RegisterDashwindSettingsSectionExample},
		{id: "dashwind-shell", sample: "dashwind-shell.html", register: RegisterDashwindShellExample},
		{id: "divider", sample: "divider.html", register: RegisterDividerExample},
		{id: "drawer", sample: "drawer.html", register: RegisterDrawerExample},
		{id: "dropdown", sample: "dropdown.html", register: RegisterDropdownExample},
		{id: "empty-state", sample: "empty_state.html", register: RegisterEmptyStateExample},
		{id: "feedback", sample: "feedback.html", register: RegisterFeedbackExample},
		{id: "file-upload", sample: "file-upload.html", register: RegisterFileUploadExample},
		{id: "font-icon", sample: "font-icon.html", register: RegisterFontIconExample},
		{id: "footer", sample: "footer.html", register: RegisterFooterExample},
		{id: "form-field", sample: "form_field.html", register: RegisterFormFieldExample},
		{id: "form", sample: "form.html", register: RegisterFormExample},
		{id: "grid", sample: "grid.html", register: RegisterGridExample},
		{id: "heading", sample: "heading.html", register: RegisterHeadingExample},
		{id: "hero", sample: "hero.html", register: RegisterHeroExample},
		{id: "image", sample: "image.html", register: RegisterImageExample},
		{id: "indicator", sample: "indicator.html", register: RegisterIndicatorExample},
		{id: "input", sample: "input.html", register: RegisterInputExample},
		{id: "join", sample: "join.html", register: RegisterJoinExample},
		{id: "kbd", sample: "kbd.html", register: RegisterKbdExample},
		{id: "layout", sample: "layout.html", register: RegisterLayoutExample},
		{id: "link", sample: "link.html", register: RegisterLinkExample},
		{id: "loading-with-variants", sample: "loading-with-variants.html", register: RegisterLoadingWithVariantsExample},
		{id: "loading", sample: "loading.html", register: RegisterLoadingExample},
		{id: "markdown", sample: "markdown.html", register: RegisterMarkdownExample},
		{id: "mask", sample: "mask.html", register: RegisterMaskExample},
		{id: "menu", sample: "menu.html", register: RegisterMenuExample},
		{id: "mockup-window", sample: "mockup-window.html", register: RegisterMockupWindowExample},
		{id: "modal", sample: "modal.html", register: RegisterModalExample},
		{id: "navbar", sample: "navbar.html", register: RegisterNavbarExample},
		{id: "navigation", sample: "navigation.html", register: RegisterNavigationExample},
		{id: "overlay-system", sample: "overlay-system.html", register: RegisterOverlaySystemExample},
		{id: "page-header", sample: "page_header.html", register: RegisterPageHeaderExample},
		{id: "pagination", sample: "pagination.html", register: RegisterPaginationExample},
		{id: "paragraph", sample: "paragraph.html", register: RegisterParagraphExample},
		{id: "progress", sample: "progress.html", register: RegisterProgressExample},
		{id: "radial-progress", sample: "radial-progress.html", register: RegisterRadialProgressExample},
		{id: "radio-group", sample: "radio_group.html", register: RegisterRadioGroupExample},
		{id: "range", sample: "range.html", register: RegisterRangeExample},
		{id: "rating", sample: "rating.html", register: RegisterRatingExample},
		{id: "region", sample: "region.html", register: RegisterRegionExample},
		{id: "section", sample: "section.html", register: RegisterSectionExample},
		{id: "select", sample: "select.html", register: RegisterSelectExample},
		{id: "sidebar", sample: "sidebar.html", register: RegisterSidebarExample},
		{id: "skeleton", sample: "skeleton.html", register: RegisterSkeletonExample},
		{id: "span", sample: "span.html", register: RegisterSpanExample},
		{id: "split", sample: "split.html", register: RegisterSplitExample},
		{id: "stack", sample: "stack.html", register: RegisterStackExample},
		{id: "stat", sample: "stat.html", register: RegisterStatExample},
		{id: "step", sample: "step.html", register: RegisterStepExample},
		{id: "steps", sample: "steps.html", register: RegisterStepsExample},
		{id: "switch", sample: "switch.html", register: RegisterSwitchExample},
		{id: "table-with-variants", sample: "table-with-variants.html", register: RegisterTableWithVariantsExample},
		{id: "table", sample: "table.html", register: RegisterTableExample},
		{id: "tabs", sample: "tabs.html", register: RegisterTabsExample},
		{id: "text-component", sample: "text-component.html", register: RegisterTextComponentExample},
		{id: "textarea", sample: "textarea.html", register: RegisterTextareaExample},
		{id: "theme-controller", sample: "theme-controller.html", register: RegisterThemeControllerExample},
		{id: "theme-toggle-button", sample: "theme-toggle-button.html", register: RegisterThemeToggleButtonExample},
		{id: "timeline", sample: "timeline.html", register: RegisterTimelineExample},
		{id: "toast", sample: "toast.html", register: RegisterToastExample},
		{id: "toggle", sample: "toggle.html", register: RegisterToggleExample},
		{id: "tooltip-with-variants", sample: "tooltip-with-variants.html", register: RegisterTooltipWithVariantsExample},
		{id: "tooltip", sample: "tooltip.html", register: RegisterTooltipExample},
	}

	for _, example := range examples {
		t.Run(example.id, func(t *testing.T) {
			app := mb.New()
			app.UseDaisyUITemplate()
			example.register(app)

			req := httptest.NewRequest(http.MethodGet, "/"+example.id, nil)
			rec := httptest.NewRecorder()
			app.Handler().ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("render %s: status %d", example.id, rec.Code)
			}

			rendered := extractMarionetteRoot(t, rec.Body.String())
			samplePath := filepath.Join("..", example.sample)
			sampleBytes, err := os.ReadFile(samplePath)
			if err != nil {
				t.Fatalf("read sample: %v", err)
			}
			sample := extractSampleBody(t, string(sampleBytes), example.id)
			if normalizeGalleryHTML(sample) != normalizeGalleryHTML(rendered) {
				t.Fatalf("sample %s does not match rendered Go example\n--- sample ---\n%s\n--- rendered ---\n%s", samplePath, sample, rendered)
			}
		})
	}
}

var (
	marionetteRootRe = regexp.MustCompile(`(?s)<main id="marionette-root"[^>]*>(.*)</main>\s*</body>`)
	bodyRe           = regexp.MustCompile(`(?s)<body\b[^>]*>(.*)</body>`)
	chartInitRe      = regexp.MustCompile(`(?s)<script>\s*\(function\(\) \{.*?initCharts\(document\);.*?</script>`)
	spaceRe          = regexp.MustCompile(`\s+`)
)

func extractMarionetteRoot(t *testing.T, page string) string {
	t.Helper()
	matches := marionetteRootRe.FindStringSubmatch(page)
	if matches == nil {
		t.Fatal("marionette root not found")
	}
	return strings.TrimSpace(matches[1])
}

func extractSampleBody(t *testing.T, page string, _ string) string {
	t.Helper()
	matches := bodyRe.FindStringSubmatch(page)
	if matches == nil {
		t.Fatal("sample body not found")
	}
	body := strings.TrimSpace(matches[1])
	body = chartInitRe.ReplaceAllString(body, "")
	return strings.TrimSpace(body)
}

func normalizeGalleryHTML(value string) string {
	return strings.TrimSpace(spaceRe.ReplaceAllString(value, " "))
}
