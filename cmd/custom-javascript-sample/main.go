package main

import (
	"fmt"
	"log"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type formula struct {
	Title string
	Tex   string
	Note  string
}

const mathJaxCHTMLURL = "https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js"

var formulas = []formula{
	{
		Title: "Euler identity",
		Tex:   `e^{i\pi} + 1 = 0`,
		Note:  "MathJax is loaded through app.AddScript.",
	},
	{
		Title: "Quadratic formula",
		Tex:   `x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`,
		Note:  "The equation panel is swapped by htmx.",
	},
	{
		Title: "Gaussian integral",
		Tex:   `\int_{-\infty}^{\infty} e^{-x^2}\,dx = \sqrt{\pi}`,
		Note:  "Custom JavaScript re-typesets MathJax after swaps.",
	},
}

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8082"); err != nil {
		log.Fatal(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.AddScript(mathJaxCHTMLURL)
	app.AddJavaScript(customJavaScript())
	app.AddStyle(customStyles())
	app.SetGlobal("formulaIndex", 0)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx)
	})

	app.Action("formula/next", func(ctx *mb.Context) mf.Node {
		next := ctx.UpdateGlobal("formulaIndex", func(old any) any {
			oldInt, _ := old.(int)
			return (oldInt + 1) % len(formulas)
		}).(int)
		return formulaPanel(formulas[next])
	})

	return app
}

func page(ctx *mb.Context) mf.Node {
	current := formulas[ctx.GetGlobalInt("formulaIndex")%len(formulas)]
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Custom JavaScript Sample",
			Description: "External MathJax plus app-level custom JavaScript.",
			Actions:     mf.Button("View source", mf.ComponentProps{Class: "btn-outline btn-sm"}),
		}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Script", Value: "MathJax", Description: "Loaded with AddScript", Trend: "External dependency", TrendTone: dw.ToneNeutral, Icon: metricIcon("M")},
			{Title: "Hook", Value: "1", Description: "AddJavaScript block", Trend: "afterSwap aware", TrendTone: dw.ToneSuccess, Icon: metricIcon("J")},
			{Title: "Panels", Value: fmt.Sprintf("%d", len(formulas)), Description: "Formula examples", Trend: "Cycles by action", TrendTone: dw.ToneSuccess, Icon: metricIcon("F")},
			{Title: "Swap", Value: "htmx", Description: "#formula-panel", Trend: "Partial render", TrendTone: dw.ToneNeutral, Icon: metricIcon("H")},
		}}),
		mf.Split(mf.SplitProps{
			Main: mf.Region(mf.RegionProps{ID: "formula-panel"}, formulaPanel(current)),
			Aside: mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
				dw.CardPanel(dw.CardPanelProps{
					Title:       "What this demonstrates",
					Description: "MathJax is registered with AddScript. The status line and MathJax re-rendering are handled by AddJavaScript.",
				},
					mf.Badge(mf.BadgeProps{Label: "Browser hook", Props: mf.ComponentProps{Class: "badge-info"}}),
				),
				dw.CardPanel(dw.CardPanelProps{
					Title:       "Inline JavaScript hook",
					Description: "This status is written by custom JavaScript after the page loads.",
				}, mf.Raw(`<p id="custom-js-status" class="custom-js-status">Waiting for custom JavaScript...</p>`)),
			),
			AsideWidth: "sm",
			Gap:        "6",
		}),
	)
	return appShell(body)
}

func formulaPanel(f formula) mf.Node {
	return dw.CardPanel(dw.CardPanelProps{
		Title:       f.Title,
		Description: f.Note,
		Actions: mf.ActionForm(mf.ActionFormProps{
			Action: "/formula/next",
			Target: "#formula-panel",
			Swap:   "innerHTML",
		}, mf.SubmitButton("Next formula", mf.ComponentProps{Variant: "primary", Size: "sm"})),
	},
		mf.Raw(fmt.Sprintf(`<div class="math-sample-equation" data-mathjax>\[%s\]</div>`, f.Tex)),
	)
}

func appShell(body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "Formula Lab", Subtitle: "Custom JS sample", Mark: "F"},
		CurrentPath:       "/",
		SearchPlaceholder: "Search formulas",
		Navigation: []dw.NavGroup{{
			Label: "Sample",
			Items: []dw.NavItem{{Path: "/", Label: "MathJax", Icon: "M"}, {Path: "#", Label: "Hooks", Icon: "J", Badge: "1"}},
		}},
		User: dw.UserMenu{Name: "Demo User", Email: "math@example.com", Initials: "DU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}

func customJavaScript() string {
	return `
		(function() {
			function setStatus(message) {
				var status = document.getElementById("custom-js-status");
				if (status) status.textContent = message;
			}

			function typeset(root) {
				if (!window.MathJax || !window.MathJax.typesetPromise) return;
				var scope = root || document;
				var nodes = scope.querySelectorAll ? scope.querySelectorAll("[data-mathjax]") : [];
				if (nodes.length === 0) return;
				window.MathJax.typesetPromise(Array.prototype.slice.call(nodes));
			}

			document.addEventListener("DOMContentLoaded", function() {
				setStatus("Custom JavaScript loaded. MathJax is ready for equations.");
				typeset(document);
			});

			document.addEventListener("htmx:afterSwap", function(event) {
				typeset(event.detail && event.detail.elt ? event.detail.elt : document);
			});
		})();
	`
}

func customStyles() string {
	return `
		.math-sample-equation {
			border: 1px solid color-mix(in oklab, var(--color-primary) 28%, transparent);
			border-radius: 0.5rem;
			background: color-mix(in oklab, var(--color-primary) 8%, var(--color-base-100));
			padding: 1.5rem;
			font-size: 1.35rem;
			overflow-x: auto;
		}

		.custom-js-status {
			color: var(--color-success);
			font-weight: 600;
		}
	`
}
