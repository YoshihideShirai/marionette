package main

import (
	"fmt"
	"log"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

type formula struct {
	Title string
	Tex   string
	Note  string
}

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
	app := mb.New()
	app.AddScript("https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js")
	app.AddJavaScript(customJavaScript())
	app.AddStyle(customStyles())
	app.Set("formulaIndex", 0)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx)
	})

	app.Action("formula/next", func(ctx *mb.Context) mf.Node {
		next := (ctx.GetInt("formulaIndex") + 1) % len(formulas)
		ctx.Set("formulaIndex", next)
		return formulaPanel(formulas[next])
	})

	if err := app.Run("127.0.0.1:8082"); err != nil {
		log.Fatal(err)
	}
}

func page(ctx *mb.Context) mf.Node {
	current := formulas[ctx.GetInt("formulaIndex")%len(formulas)]
	return mf.Container(mf.ContainerProps{MaxWidth: "3xl", Centered: true},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.PageHeader(mf.PageHeaderProps{
				Title:       "Custom JavaScript Sample",
				Description: "External MathJax plus app-level custom JavaScript.",
			}),
			mf.Alert(mf.AlertProps{
				Title:       "What this demonstrates",
				Description: "MathJax is registered with AddScript. The status line and MathJax re-rendering are handled by AddJavaScript.",
				Props:       mf.ComponentProps{Variant: "info"},
			}),
			mf.Region(mf.RegionProps{ID: "formula-panel"}, formulaPanel(current)),
			mf.Card(mf.CardProps{
				Title:       "Inline JavaScript hook",
				Description: "This status is written by custom JavaScript after the page loads.",
			}, mf.Raw(`<p id="custom-js-status" class="custom-js-status">Waiting for custom JavaScript...</p>`)),
		),
	)
}

func formulaPanel(f formula) mf.Node {
	return mf.Card(mf.CardProps{
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
