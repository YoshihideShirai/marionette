package twailwindcss

const (
	TemplateName = "tailwindcss"
	BrowserURL   = "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"
)

func FrameworkStylesheets() []string {
	return nil
}

func FrameworkScripts() []string {
	return []string{BrowserURL}
}
