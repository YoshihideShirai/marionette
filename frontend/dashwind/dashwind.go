// Package dashwind provides reusable DashWind-inspired admin shell components.
package dashwind

import (
	"strings"

	mf "github.com/YoshihideShirai/marionette/frontend"
	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
)

// DefaultMainTargetID is the default region id used by ShellContent for htmx fragment swaps.
const DefaultMainTargetID = "dashwind-main"

// DefaultCSS contains small layout tweaks used by the DashWind shell.
const DefaultCSS = `
#marionette-root { width: 100%; max-width: none; padding: 0; }
#marionette-root > * { animation: none; }
.dashwind-shell .drawer-side .menu a.active { background: var(--color-primary); color: var(--color-primary-content); font-weight: 700; }
.dashwind-shell .card, .dashwind-shell .stats { border: 1px solid color-mix(in oklab, var(--color-base-content) 10%, transparent); }
.dashwind-shell .stat-value { font-size: clamp(1.75rem, 2vw, 2.25rem); }
`

// NavItem describes a single sidebar link.
type NavItem struct {
	Label  string
	Href   string
	Icon   string
	Active bool
	Class  string
}

// NavGroup describes a labeled sidebar menu section.
type NavGroup struct {
	Label string
	Items []NavItem
	Class string
}

// ShellProps configures the DashWind application shell.
type ShellProps struct {
	ID                string
	DrawerID          string
	CurrentTitle      string
	BrandTitle        string
	BrandSubtitle     string
	BrandMark         string
	NavGroups         []NavGroup
	Content           mf.Node
	MainTargetID      string
	SearchPlaceholder string
	UserInitials      string
	SidebarFooter     mf.Node
	TopbarActions     []mf.Node
	Class             string
}

// Shell renders a responsive dashboard shell with a drawer sidebar and sticky topbar.
func Shell(props ShellProps) mf.Node {
	id := defaultString(props.ID, "dashwind-app")
	drawerID := defaultString(props.DrawerID, "dashwind-drawer")
	current := defaultString(props.CurrentTitle, "Dashboard")
	content := mf.DivProps(mf.ElementProps{Class: "flex min-h-screen flex-col bg-base-200"}, topbar(props, drawerID, current), ShellContent(props.MainTargetID, props.Content))
	return mf.Region(mf.RegionProps{ID: id}, daisy.DrawerWithProps(daisy.DrawerProps{
		ID:           drawerID,
		Class:        strings.TrimSpace("lg:drawer-open dashwind-shell " + props.Class),
		ContentClass: "flex min-h-screen flex-col bg-base-200",
		SideClass:    "z-40",
		Content:      content,
		Side:         sidebar(props, current),
	}))
}

// ShellContent wraps page content in the htmx-friendly main region used by Shell.
func ShellContent(id string, body mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: defaultString(id, DefaultMainTargetID), Props: mf.ComponentProps{Class: "flex-1 p-4 md:p-6 lg:p-8 space-y-6"}}, body)
}

// PageHeaderProps configures a page title, description, and optional action node.
type PageHeaderProps struct {
	Title       string
	Description string
	Actions     mf.Node
	Class       string
}

// PageHeader renders a responsive title row for dashboard pages.
func PageHeader(props PageHeaderProps) mf.Node {
	children := []mf.Node{div("", mf.H1Props(mf.ElementProps{Class: "text-3xl font-bold"}, mf.Text(props.Title)), paragraph("mt-1 text-base-content/60", props.Description))}
	if props.Actions != nil {
		children = append(children, props.Actions)
	}
	return div(strings.TrimSpace("flex flex-col gap-4 md:flex-row md:items-center md:justify-between "+props.Class), children...)
}

// CardPanelProps configures a DashWind card surface.
type CardPanelProps struct {
	Title       string
	Description string
	Class       string
	BodyClass   string
	Actions     mf.Node
}

// CardPanel renders a shadowed DaisyUI card panel.
func CardPanel(props CardPanelProps, children ...mf.Node) mf.Node {
	className := strings.TrimSpace("bg-base-100 shadow " + props.Class)
	return daisy.CardPanel(daisy.CardPanelProps{Title: props.Title, Description: props.Description, Class: className, BodyClass: props.BodyClass, Actions: props.Actions}, children...)
}

// Stat describes a single metric item inside StatsGrid.
type Stat struct {
	Title            string
	Value            string
	Description      string
	Figure           mf.Node
	Class            string
	TitleClass       string
	ValueClass       string
	DescriptionClass string
	FigureClass      string
}

// StatsGridProps configures a grid of stat cards.
type StatsGridProps struct {
	Items []Stat
	Class string
}

// StatsGrid renders each stat as an individual card, matching DashWind dashboard tiles.
func StatsGrid(props StatsGridProps) mf.Node {
	cards := make([]mf.Node, 0, len(props.Items))
	for _, item := range props.Items {
		cards = append(cards, daisy.StatsWithProps(daisy.StatsProps{Class: "shadow bg-base-100"}, daisy.StatItem(daisy.StatProps{
			Title: item.Title, Value: item.Value, Description: item.Description, Figure: item.Figure, Class: item.Class,
			TitleClass: item.TitleClass, ValueClass: defaultString(item.ValueClass, "text-primary"), DescriptionClass: item.DescriptionClass, FigureClass: defaultString(item.FigureClass, "text-primary text-3xl"),
		})))
	}
	return div(defaultString(props.Class, "grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4"), cards...)
}

// DataTableProps configures a responsive DaisyUI table.
type DataTableProps struct {
	Headers      []string
	Rows         [][]mf.Node
	Class        string
	WrapperClass string
}

// DataTable renders a responsive table for DashWind pages.
func DataTable(props DataTableProps) mf.Node {
	return daisy.TableWithProps(daisy.TableProps{Headers: props.Headers, Rows: props.Rows, Class: props.Class, WrapperClass: props.WrapperClass})
}

func topbar(props ShellProps, drawerID, current string) mf.Node {
	actions := append([]mf.Node{}, props.TopbarActions...)
	if len(actions) == 0 {
		actions = []mf.Node{
			daisy.ButtonWithAttrs("◐", mf.ComponentProps{Class: "btn-ghost btn-circle"}, map[string]string{"type": "button", "onclick": "mrnToggleTheme()", "aria-label": "toggle theme"}),
			daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-ghost btn-circle indicator"}, map[string]string{"type": "button", "aria-label": "notifications"}, span("indicator-item badge badge-primary badge-xs", ""), mf.Text("🔔")),
			daisy.AvatarPlaceholder(defaultString(props.UserInitials, "DW"), "", "bg-primary text-primary-content w-10 rounded-full"),
		}
	}
	return daisy.NavbarWithProps(daisy.NavbarProps{Class: "sticky top-0 z-30 border-b border-base-300 bg-base-100/90 backdrop-blur"},
		div("flex-none lg:hidden", mf.LabelElementProps(mf.ElementProps{Class: "btn btn-square btn-ghost", Attrs: mf.Attrs{"for": drawerID, "aria-label": "open sidebar"}}, mf.Text("☰"))),
		div("flex-1", mf.H1Props(mf.ElementProps{Class: "text-xl font-semibold"}, mf.Text(current))),
		div("hidden max-w-md flex-1 md:block", searchInput(defaultString(props.SearchPlaceholder, "Search DashWind"))),
		div("flex-none gap-2", actions...),
	)
}

func sidebar(props ShellProps, current string) mf.Node {
	footer := props.SidebarFooter
	if footer == nil {
		footer = div("mt-6 rounded-box bg-primary/10 p-4 text-sm", paragraph("font-semibold", "Marionette port"), paragraph("mt-1 opacity-70", "DashWind template patterns rebuilt as Go handlers and htmx fragments."))
	}
	brandTitle := defaultString(props.BrandTitle, "DashWind")
	brandMark := defaultString(props.BrandMark, "D")
	panelChildren := []mf.Node{
		div("mb-6 flex items-center gap-3",
			div("grid h-11 w-11 place-items-center rounded-2xl bg-primary text-xl font-black text-primary-content", mf.Text(brandMark)),
			div("", mf.H2Props(mf.ElementProps{Class: "text-lg font-bold"}, mf.Text(brandTitle)), paragraph("text-xs text-base-content/60", props.BrandSubtitle)),
		),
	}
	for _, group := range props.NavGroups {
		panelChildren = append(panelChildren, NavMenu(group, current))
	}
	panelChildren = append(panelChildren, footer)
	return div("min-h-full w-80 bg-base-100 text-base-content shadow-xl", div("p-5", panelChildren...))
}

// NavMenu renders a menu group and marks an item active by CurrentTitle when Active is false.
func NavMenu(group NavGroup, currentTitle string) mf.Node {
	items := []mf.Node{daisy.MenuTitle(group.Label)}
	for _, item := range group.Items {
		active := item.Active || (item.Label != "" && item.Label == currentTitle)
		items = append(items, daisy.MenuLink(daisy.MenuLinkProps{Label: item.Label, Href: item.Href, Icon: item.Icon, Active: active, Class: item.Class}))
	}
	return daisy.MenuWithProps(daisy.MenuProps{Class: strings.TrimSpace("rounded-box gap-1 p-0 " + group.Class)}, items...)
}

func searchInput(placeholder string) mf.Node {
	return mf.LabelElementProps(mf.ElementProps{Class: "input input-bordered flex items-center gap-2"},
		span("opacity-60", "⌕"),
		mf.InputElement(mf.ElementProps{Class: "grow", Attrs: mf.Attrs{"type": "search", "placeholder": placeholder}}),
	)
}

func div(className string, children ...mf.Node) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: className}, children...)
}
func paragraph(className, text string) mf.Node {
	return mf.PProps(mf.ElementProps{Class: className}, mf.Text(text))
}
func span(className, text string) mf.Node {
	return mf.SpanProps(mf.ElementProps{Class: className}, mf.Text(text))
}
func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
