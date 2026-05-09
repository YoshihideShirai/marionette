package html

import (
	"html/template"
	"strings"
	"testing"
)

func TestElementEscapesTextAndAttrs(t *testing.T) {
	rendered, err := DivProps(
		ElementProps{
			ID:    "panel",
			Class: "card",
			Attrs: Attrs{
				"data-value": `<script>alert("x")</script>`,
			},
		},
		Text(`<b>Hello</b>`),
	).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := string(rendered)
	if !strings.Contains(got, `data-value="&lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;"`) {
		t.Fatalf("rendered attrs were not escaped: %s", got)
	}
	if !strings.Contains(got, `&lt;b&gt;Hello&lt;/b&gt;`) {
		t.Fatalf("rendered text was not escaped: %s", got)
	}
}

func TestElementRejectsInvalidTags(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{name: "empty", tag: ""},
		{name: "starts with digit", tag: "1div"},
		{name: "contains space", tag: "div script"},
		{name: "contains greater-than", tag: "div>"},
		{name: "contains slash", tag: "x/y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Element(tt.tag, ElementProps{}).Render()
			if err == nil {
				t.Fatal("Render() error = nil, want invalid tag error")
			}
		})
	}
}

func FuzzElementEscapesTextAndAttrs(f *testing.F) {
	for _, seed := range []struct {
		text string
		attr string
	}{
		{text: `<text> "quoted" & 'single'`, attr: `<attr> "quoted" & 'single'`},
		{text: `<`, attr: `>`},
		{text: `"`, attr: `&`},
		{text: `'`, attr: `<>&"'`},
	} {
		f.Add(seed.text, seed.attr)
	}

	f.Fuzz(func(t *testing.T, text, attr string) {
		rendered, err := ElementNode{
			Tag:   "div",
			Attrs: map[string]string{"data-value": attr},
			Text:  text,
		}.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		want := `<div data-value="` + template.HTMLEscapeString(attr) + `">` + template.HTMLEscapeString(text) + `</div>`
		if got := string(rendered); got != want {
			t.Fatalf("Render() = %q, want %q", got, want)
		}
	})
}

func TestElementAllowsValidCustomAndSvgTags(t *testing.T) {
	for _, tag := range []string{"my-element", "svg", "path"} {
		t.Run(tag, func(t *testing.T) {
			rendered, err := Element(tag, ElementProps{Attrs: Attrs{"data-tag": tag}}, Text("ok")).Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			want := `<` + tag + ` data-tag="` + tag + `"><span>ok</span></` + tag + `>`
			if got := string(rendered); got != want {
				t.Fatalf("Render() = %q, want %q", got, want)
			}
		})
	}
}

func TestElementAttrsRenderInDeterministicKeyOrder(t *testing.T) {
	rendered, err := ElementNode{
		Tag: "div",
		Attrs: map[string]string{
			"zeta":  "last",
			"alpha": "first",
			"data":  "middle",
		},
	}.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := `<div alpha="first" data="middle" zeta="last"></div>`
	if got := string(rendered); got != want {
		t.Fatalf("Render() = %q, want deterministic sorted attrs %q", got, want)
	}
}

func TestElementSkipsNilChildWithoutPanic(t *testing.T) {
	rendered, err := Div(Text("before"), nil, Text("after")).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := `<div><span>before</span><span>after</span></div>`
	if got := string(rendered); got != want {
		t.Fatalf("Render() = %q, want nil child skipped as %q", got, want)
	}
}

func TestRawDangerouslyRendersUnescapedTrustedHTML(t *testing.T) {
	const unsafe = `<script>alert("x")</script><b data-x="&">raw</b>`

	rendered, err := Raw(unsafe).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got := string(rendered); got != unsafe {
		t.Fatalf("Raw Render() = %q, want unescaped trusted HTML %q", got, unsafe)
	}
}

func TestSemanticElementHelpersRenderExpectedTags(t *testing.T) {
	rendered, err := LabelElementProps(
		ElementProps{Class: "input", Attrs: Attrs{"for": "search"}},
		Text("Search"),
		InputElement(ElementProps{ID: "search", Attrs: Attrs{"type": "search"}}),
	).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	got := string(rendered)
	for _, want := range []string{`<label`, `class="input"`, `for="search"`, `<input id="search" type="search"></input>`, `Search`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}

	list, err := UlProps(ElementProps{Class: "list"}, Li(Text("One")), Li(Text("Two"))).Render()
	if err != nil {
		t.Fatalf("Render() list error = %v", err)
	}
	if got := string(list); !strings.Contains(got, `<ul class="list"><li><span>One</span></li><li><span>Two</span></li></ul>`) {
		t.Fatalf("unexpected list markup: %q", got)
	}
}

func TestHTMXPrimitivesRenderExpectedMarkup(t *testing.T) {
	buttonHTML, err := HTMXButton("Save <now>").Post("/items/save").Target("#items").Render()
	if err != nil {
		t.Fatalf("button Render() error = %v", err)
	}
	button := string(buttonHTML)
	for _, want := range []string{`hx-post="/items/save"`, `hx-target="#items"`, `hx-swap="outerHTML"`, `Save &lt;now&gt;`} {
		if !strings.Contains(button, want) {
			t.Fatalf("expected %q in %q", want, button)
		}
	}

	formHTML, err := NewForm("items/create", Input("name", `<Aiko>`), Submit("Create")).Target("#items").Render()
	if err != nil {
		t.Fatalf("form Render() error = %v", err)
	}
	form := string(formHTML)
	for _, want := range []string{`hx-post="/items/create"`, `hx-target="#items"`, `name="name"`, `value="&lt;Aiko&gt;"`, `type="submit"`} {
		if !strings.Contains(form, want) {
			t.Fatalf("expected %q in %q", want, form)
		}
	}
}

func TestNavigationAndTablePrimitivesRenderExpectedMarkup(t *testing.T) {
	sidebarHTML, err := NewSidebar("Marionette", "Admin <Console>", SidebarLink("Home", "/").Active()).Note("Demo", `<unsafe>`).Render()
	if err != nil {
		t.Fatalf("sidebar Render() error = %v", err)
	}
	sidebar := string(sidebarHTML)
	for _, want := range []string{`<aside`, `href="/"`, `btn btn-primary justify-start`, `Admin &lt;Console&gt;`, `&lt;unsafe&gt;`} {
		if !strings.Contains(sidebar, want) {
			t.Fatalf("expected %q in %q", want, sidebar)
		}
	}

	tableHTML, err := HTMXTable([]string{"Name"}, TableRow(Text(`<Aiko>`))).Render()
	if err != nil {
		t.Fatalf("table Render() error = %v", err)
	}
	if got := string(tableHTML); !strings.Contains(got, `<td><span>&lt;Aiko&gt;</span></td>`) {
		t.Fatalf("expected escaped table cell in %q", got)
	}
}
