// Package dashwind provides reusable DashWind-inspired admin shell components.
package dashwind

import (
	"strings"

	mf "github.com/YoshihideShirai/marionette/frontend"
	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
)

// DefaultMainTargetID is the default region id used by ShellContent for htmx fragment swaps.
const DefaultMainTargetID = "dashwind-main"

// DefaultDrawerID is the safe default drawer toggle id used by Shell.
const DefaultDrawerID = "dashwind-drawer"

// DefaultCSS contains small layout tweaks used by the DashWind shell.
const DefaultCSS = `
#marionette-root { width: 100%; max-width: none; padding: 0; }
#marionette-root > * { animation: none; }
.dashwind-shell .drawer-side .menu a.active { background: var(--color-primary); color: var(--color-primary-content); font-weight: 700; }
.dashwind-shell .card, .dashwind-shell .stats { border: 1px solid color-mix(in oklab, var(--color-base-content) 10%, transparent); }
.dashwind-shell .stat-value { font-size: clamp(1.75rem, 2vw, 2.25rem); }
`

// Brand describes the shell brand block rendered at the top of the sidebar.
type Brand struct {
	Title    string
	Subtitle string
	Mark     string
	Href     string
	Class    string
}

// Navigation describes a set of sidebar menu groups.
type Navigation []NavGroup

// NavItem describes a single sidebar link.
type NavItem struct {
	Path     string
	Label    string
	Icon     string
	Badge    string
	Disabled bool
	External bool
	Children []NavItem

	// Deprecated: use Path.
	Href string
	// Deprecated: active state is normally derived by RenderNavigation from the current path.
	Active bool
	// Class appends classes to the rendered anchor.
	Class string
}

// NavGroup describes a labeled sidebar menu section.
type NavGroup struct {
	Label string
	Items []NavItem
	Class string
}

// UserMenu configures the user affordance rendered in the topbar action area.
type UserMenu struct {
	Name        string
	Email       string
	Initials    string
	AvatarURL   string
	Items       []NavItem
	Actions     []mf.Node
	Class       string
	AvatarClass string
	MenuClass   string
}

// ShellProps configures the DashWind application shell.
type ShellProps struct {
	ID                string
	DrawerID          string
	MainTargetID      string
	Brand             Brand
	CurrentPath       string
	Navigation        Navigation
	User              UserMenu
	Actions           []mf.Node
	SearchPlaceholder string
	Class             string

	// Optional class overrides. Defaults match the DashWind demo shell and are safe for responsive layouts.
	DrawerClass  string
	ContentClass string
	SideClass    string
	OverlayClass string
	NavbarClass  string
	MainClass    string
	SidebarClass string

	// Deprecated: use Brand.Title, Brand.Subtitle, Brand.Mark, Navigation, User.Initials, Actions, and ShellContent/body instead.
	CurrentTitle  string
	BrandTitle    string
	BrandSubtitle string
	BrandMark     string
	NavGroups     Navigation
	Content       mf.Node
	UserInitials  string
	SidebarFooter mf.Node
	TopbarActions []mf.Node
}

// Shell renders a responsive dashboard shell with a drawer sidebar and sticky topbar.
func Shell(props ShellProps, body mf.Node) mf.Node {
	props = normalizeShellProps(props)
	if body == nil {
		body = props.Content
	}
	id := defaultString(props.ID, "dashwind-app")
	drawerID := defaultString(props.DrawerID, DefaultDrawerID)
	content := mf.DivProps(mf.ElementProps{Class: shellContentClass(props)}, topbar(props, drawerID), ShellContentWithClass(props.MainTargetID, props.MainClass, body))
	return mf.Region(mf.RegionProps{ID: id}, daisy.DrawerWithProps(daisy.DrawerProps{
		ID:           drawerID,
		Class:        strings.TrimSpace("lg:drawer-open dashwind-shell " + props.DrawerClass + " " + props.Class),
		ContentClass: defaultString(props.ContentClass, "flex min-h-screen flex-col bg-base-200"),
		SideClass:    defaultString(props.SideClass, "z-40"),
		OverlayClass: props.OverlayClass,
		Content:      content,
		Side:         sidebar(props),
	}))
}

// ShellContent wraps page content in the htmx-friendly main region used by Shell.
func ShellContent(id string, body mf.Node) mf.Node {
	return ShellContentWithClass(id, "", body)
}

// ShellContentWithClass wraps page content with an optional class override for Shell's main region.
func ShellContentWithClass(id, className string, body mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: defaultString(id, DefaultMainTargetID), Props: mf.ComponentProps{Class: defaultString(className, "flex-1 p-4 md:p-6 lg:p-8 space-y-6")}}, body)
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

// Tone describes the semantic color applied to dashboard metric trends.
type Tone string

const (
	// ToneSuccess marks positive KPI movement.
	ToneSuccess Tone = "success"
	// ToneWarning marks KPI movement that needs attention.
	ToneWarning Tone = "warning"
	// ToneError marks negative KPI movement.
	ToneError Tone = "error"
	// ToneNeutral marks informational or unchanged KPI movement.
	ToneNeutral Tone = "neutral"
)

// Metric describes a KPI tile rendered by MetricGrid or MetricCard.
type Metric struct {
	Title       string
	Value       string
	Description string
	Trend       string
	TrendTone   Tone
	Icon        mf.Node
	Href        string
}

// MetricGridProps configures a responsive grid of metric cards.
type MetricGridProps struct {
	Items []Metric
	Class string
}

// MetricGrid renders KPI metric cards in a responsive grid.
func MetricGrid(props MetricGridProps) mf.Node {
	cards := make([]mf.Node, 0, len(props.Items))
	for _, item := range props.Items {
		cards = append(cards, MetricCard(item))
	}
	return div(defaultString(props.Class, "grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4"), cards...)
}

// MetricCard renders a single DashWind KPI card.
func MetricCard(metric Metric) mf.Node {
	children := []mf.Node{
		div("flex items-start justify-between gap-4",
			div("space-y-1", paragraph("text-sm font-medium text-base-content/60", metric.Title), div("text-3xl font-bold tracking-tight text-primary", mf.Text(metric.Value))),
			metricIcon(metric.Icon),
		),
	}
	if strings.TrimSpace(metric.Description) != "" {
		children = append(children, paragraph("text-sm text-base-content/60", metric.Description))
	}
	if strings.TrimSpace(metric.Trend) != "" {
		children = append(children, paragraph(strings.TrimSpace("text-sm font-medium "+toneTextClass(metric.TrendTone)), metric.Trend))
	}

	className := "card bg-base-100 shadow transition hover:-translate-y-0.5 hover:shadow-md"
	body := mf.DivProps(mf.ElementProps{Class: "card-body gap-3"}, children...)
	if strings.TrimSpace(metric.Href) == "" {
		return mf.DivProps(mf.ElementProps{Class: className}, body)
	}
	return mf.AnchorProps(mf.ElementProps{Class: strings.TrimSpace(className + " block no-underline text-base-content"), Attrs: mf.Attrs{"href": metric.Href}}, body)
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
// Deprecated: use MetricGrid with Metric items.
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

// Column describes a high-level DashWind data table column for a row of type T.
type Column[T any] struct {
	Header      string
	Cell        func(T) mf.Node
	Class       string
	HeaderClass string
	Sortable    bool
}

// DataTableProps configures a responsive DashWind data table.
type DataTableProps[T any] struct {
	Columns      []Column[T]
	Rows         []T
	Empty        mf.Node
	Zebra        bool
	Compact      bool
	Actions      mf.Node
	WrapperClass string
}

// DataTable renders a responsive high-level DashWind table backed by DaisyUI's table primitive.
func DataTable[T any](props DataTableProps[T]) mf.Node {
	headers := make([]string, 0, len(props.Columns))
	headerClasses := make([]string, 0, len(props.Columns))
	headerSortables := make([]bool, 0, len(props.Columns))
	for _, column := range props.Columns {
		headers = append(headers, column.Header)
		headerClasses = append(headerClasses, column.HeaderClass)
		headerSortables = append(headerSortables, column.Sortable)
	}

	rows := make([][]mf.Node, 0, len(props.Rows))
	cellClasses := make([][]string, 0, len(props.Rows))
	for _, row := range props.Rows {
		cells := make([]mf.Node, 0, len(props.Columns))
		classes := make([]string, 0, len(props.Columns))
		for _, column := range props.Columns {
			cell := mf.Node(mf.Text(""))
			if column.Cell != nil {
				cell = column.Cell(row)
			}
			cells = append(cells, cell)
			classes = append(classes, column.Class)
		}
		rows = append(rows, cells)
		cellClasses = append(cellClasses, classes)
	}

	className := ""
	if props.Zebra {
		className = strings.TrimSpace(className + " table-zebra")
	}
	if props.Compact {
		className = strings.TrimSpace(className + " table-sm")
	}
	table := daisy.TableWithProps(daisy.TableProps{
		Headers:         headers,
		Rows:            rows,
		Class:           className,
		WrapperClass:    props.WrapperClass,
		HeaderClasses:   headerClasses,
		HeaderSortables: headerSortables,
		CellClasses:     cellClasses,
		Empty:           props.Empty,
		EmptyColSpan:    len(props.Columns),
	})
	if props.Actions == nil {
		return table
	}
	return div("space-y-4", div("flex justify-end", props.Actions), table)
}

func normalizeShellProps(props ShellProps) ShellProps {
	if props.Brand.Title == "" {
		props.Brand.Title = props.BrandTitle
	}
	if props.Brand.Subtitle == "" {
		props.Brand.Subtitle = props.BrandSubtitle
	}
	if props.Brand.Mark == "" {
		props.Brand.Mark = props.BrandMark
	}
	if len(props.Navigation) == 0 {
		props.Navigation = props.NavGroups
	}
	if props.User.Initials == "" {
		props.User.Initials = props.UserInitials
	}
	if len(props.Actions) == 0 {
		props.Actions = props.TopbarActions
	}
	if strings.TrimSpace(props.CurrentPath) == "" && strings.TrimSpace(props.CurrentTitle) != "" {
		if item, ok := findNavItemByLabel(props.Navigation, props.CurrentTitle); ok {
			props.CurrentPath = navItemPath(item)
		}
	}
	return props
}

func topbar(props ShellProps, drawerID string) mf.Node {
	actions := append([]mf.Node{}, props.Actions...)
	if len(actions) == 0 {
		actions = []mf.Node{
			daisy.ButtonWithAttrs("◐", mf.ComponentProps{Class: "btn-ghost btn-circle"}, map[string]string{"type": "button", "onclick": "mrnToggleTheme()", "aria-label": "toggle theme"}),
			daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-ghost btn-circle indicator"}, map[string]string{"type": "button", "aria-label": "notifications"}, span("indicator-item badge badge-primary badge-xs", ""), mf.Text("🔔")),
			renderUserMenu(props.User),
		}
	}
	return daisy.NavbarWithProps(daisy.NavbarProps{Class: defaultString(props.NavbarClass, "sticky top-0 z-30 border-b border-base-300 bg-base-100/90 backdrop-blur")},
		div("flex-none lg:hidden", mf.LabelElementProps(mf.ElementProps{Class: "btn btn-square btn-ghost", Attrs: mf.Attrs{"for": drawerID, "aria-label": "open sidebar"}}, mf.Text("☰"))),
		div("flex-1", mf.H1Props(mf.ElementProps{Class: "text-xl font-semibold"}, mf.Text(currentTitle(props)))),
		div("hidden max-w-md flex-1 md:block", searchInput(defaultString(props.SearchPlaceholder, "Search DashWind"))),
		div("flex-none gap-2", actions...),
	)
}

func sidebar(props ShellProps) mf.Node {
	footer := props.SidebarFooter
	if footer == nil {
		footer = div("mt-6 rounded-box bg-primary/10 p-4 text-sm", paragraph("font-semibold", "Marionette port"), paragraph("mt-1 opacity-70", "DashWind template patterns rebuilt as Go handlers and htmx fragments."))
	}
	brand := props.Brand
	brand.Title = defaultString(brand.Title, "DashWind")
	brand.Mark = defaultString(brand.Mark, "D")
	panelChildren := []mf.Node{brandBlock(brand)}
	panelChildren = append(panelChildren, RenderNavigation(props.Navigation, props.CurrentPath))
	panelChildren = append(panelChildren, footer)
	return div(defaultString(props.SidebarClass, "min-h-full w-80 bg-base-100 text-base-content shadow-xl"), div("p-5", panelChildren...))
}

func brandBlock(brand Brand) mf.Node {
	mark := div("grid h-11 w-11 place-items-center rounded-2xl bg-primary text-xl font-black text-primary-content", mf.Text(brand.Mark))
	text := div("", mf.H2Props(mf.ElementProps{Class: "text-lg font-bold"}, mf.Text(brand.Title)), paragraph("text-xs text-base-content/60", brand.Subtitle))
	children := []mf.Node{mark, text}
	className := strings.TrimSpace("mb-6 flex items-center gap-3 " + brand.Class)
	if strings.TrimSpace(brand.Href) != "" {
		return mf.AnchorProps(mf.ElementProps{Class: className, Attrs: mf.Attrs{"href": brand.Href}}, children...)
	}
	return div(className, children...)
}

// RenderNavigation renders DashWind sidebar navigation and marks the item whose Path matches currentPath active.
func RenderNavigation(nav Navigation, currentPath string) mf.Node {
	menus := make([]mf.Node, 0, len(nav))
	for _, group := range nav {
		menus = append(menus, NavMenu(group, currentPath))
	}
	return mf.DivProps(mf.ElementProps{}, menus...)
}

// NavMenu renders a menu group and marks an item active by CurrentPath when Active is false.
func NavMenu(group NavGroup, currentPath string) mf.Node {
	items := []mf.Node{}
	if strings.TrimSpace(group.Label) != "" {
		items = append(items, daisy.MenuTitle(group.Label))
	}
	for _, item := range group.Items {
		items = append(items, renderNavItem(item, currentPath))
	}
	return daisy.MenuWithProps(daisy.MenuProps{Class: strings.TrimSpace("rounded-box gap-1 p-0 " + group.Class)}, items...)
}

func renderNavItem(item NavItem, currentPath string) mf.Node {
	path := navItemPath(item)
	active := item.Active || (path != "" && path == currentPath)
	className := strings.TrimSpace(item.Class)
	if active {
		className = strings.TrimSpace(className + " active")
	}
	if item.Disabled {
		className = strings.TrimSpace(className + " disabled")
	}
	attrs := mf.Attrs{"class": className}
	if item.Disabled {
		attrs["aria-disabled"] = "true"
		attrs["tabindex"] = "-1"
	} else {
		if path == "" {
			path = "#"
		}
		attrs["href"] = path
		if item.External {
			attrs["target"] = "_blank"
			attrs["rel"] = "noopener noreferrer"
		}
	}
	children := make([]mf.Node, 0, 3)
	if item.Icon != "" {
		children = append(children, span("w-6 text-center", item.Icon))
	}
	children = append(children, span("flex-1", item.Label))
	if item.Badge != "" {
		children = append(children, span("badge badge-sm", item.Badge))
	}
	liChildren := []mf.Node{mf.Element("a", mf.ElementProps{Attrs: attrs}, children...)}
	if len(item.Children) > 0 {
		childItems := make([]mf.Node, 0, len(item.Children))
		for _, child := range item.Children {
			childItems = append(childItems, renderNavItem(child, currentPath))
		}
		liChildren = append(liChildren, mf.UlProps(mf.ElementProps{}, childItems...))
	}
	return mf.Li(liChildren...)
}

func navItemPath(item NavItem) string {
	if strings.TrimSpace(item.Path) != "" {
		return item.Path
	}
	return item.Href
}

func renderUserMenu(user UserMenu) mf.Node {
	initials := defaultString(user.Initials, "DW")
	avatarClass := defaultString(user.AvatarClass, "bg-primary text-primary-content w-10 rounded-full")
	avatar := daisy.AvatarPlaceholder(initials, "", avatarClass)
	if strings.TrimSpace(user.AvatarURL) != "" {
		avatar = mf.DivProps(mf.ElementProps{Class: "avatar"}, div(defaultString(user.AvatarClass, "w-10 rounded-full"), mf.Element("img", mf.ElementProps{Attrs: mf.Attrs{"src": user.AvatarURL, "alt": defaultString(user.Name, "User")}})))
	}
	menuItems := make([]mf.Node, 0, len(user.Items)+len(user.Actions)+1)
	if user.Name != "" || user.Email != "" {
		menuItems = append(menuItems, mf.Li(mf.Element("div", mf.ElementProps{Class: "flex flex-col gap-0"}, span("font-semibold", user.Name), span("text-xs opacity-60", user.Email))))
	}
	for _, item := range user.Items {
		menuItems = append(menuItems, daisy.MenuLink(daisy.MenuLinkProps{Label: item.Label, Href: navItemPath(item), Icon: item.Icon, Active: item.Active, Class: item.Class}))
	}
	menuItems = append(menuItems, user.Actions...)
	if len(menuItems) == 0 {
		return avatar
	}
	return div(strings.TrimSpace("dropdown dropdown-end "+user.Class),
		mf.Element("button", mf.ElementProps{Class: "btn btn-ghost btn-circle", Attrs: mf.Attrs{"type": "button", "aria-label": defaultString(user.Name, "user menu")}}, avatar),
		daisy.MenuWithProps(daisy.MenuProps{Class: strings.TrimSpace("dropdown-content z-50 mt-3 w-56 rounded-box bg-base-100 p-2 shadow " + user.MenuClass)}, menuItems...),
	)
}

func currentTitle(props ShellProps) string {
	if props.CurrentTitle != "" {
		return props.CurrentTitle
	}
	if item, ok := findActiveNavItem(props.Navigation, props.CurrentPath); ok {
		return item.Label
	}
	return "Dashboard"
}

func findNavItemByLabel(nav Navigation, label string) (NavItem, bool) {
	for _, group := range nav {
		if item, ok := findNavItemInItems(group.Items, func(item NavItem) bool { return item.Label == label }); ok {
			return item, true
		}
	}
	return NavItem{}, false
}

func findActiveNavItem(nav Navigation, currentPath string) (NavItem, bool) {
	for _, group := range nav {
		if item, ok := findNavItemInItems(group.Items, func(item NavItem) bool {
			path := navItemPath(item)
			return item.Active || (path != "" && path == currentPath)
		}); ok {
			return item, true
		}
	}
	return NavItem{}, false
}

func findNavItemInItems(items []NavItem, match func(NavItem) bool) (NavItem, bool) {
	for _, item := range items {
		if match(item) {
			return item, true
		}
		if child, ok := findNavItemInItems(item.Children, match); ok {
			return child, true
		}
	}
	return NavItem{}, false
}

func shellContentClass(props ShellProps) string {
	return defaultString(props.ContentClass, "flex min-h-screen flex-col bg-base-200")
}

func searchInput(placeholder string) mf.Node {
	return mf.LabelElementProps(mf.ElementProps{Class: "input input-bordered flex items-center gap-2"},
		span("opacity-60", "⌕"),
		mf.InputElement(mf.ElementProps{Class: "grow", Attrs: mf.Attrs{"type": "search", "placeholder": placeholder}}),
	)
}

func metricIcon(icon mf.Node) mf.Node {
	if icon == nil {
		return mf.Raw("")
	}
	return div("grid h-12 w-12 shrink-0 place-items-center rounded-box bg-primary/10 text-2xl font-bold text-primary", icon)
}

func toneTextClass(tone Tone) string {
	switch tone {
	case ToneSuccess:
		return "text-success"
	case ToneWarning:
		return "text-warning"
	case ToneError:
		return "text-error"
	case ToneNeutral, "":
		return "text-base-content/60"
	default:
		return "text-" + string(tone)
	}
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
