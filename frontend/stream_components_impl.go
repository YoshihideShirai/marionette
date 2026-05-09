package frontend

import (
	"fmt"
	"strings"
)

// StreamTrigger renders a hidden htmx trigger for polling a server-side stream action.
func StreamTrigger(props StreamTriggerProps) Node {
	action := strings.TrimSpace(props.Action)
	if action == "" {
		return renderErrorNode{err: fmt.Errorf("stream trigger action is required")}
	}
	target := strings.TrimSpace(props.Target)
	if target == "" {
		return renderErrorNode{err: fmt.Errorf("stream trigger target is required")}
	}
	swap := strings.TrimSpace(props.Swap)
	if swap == "" {
		swap = "outerHTML"
	}
	trigger := strings.TrimSpace(props.Trigger)
	if trigger == "" {
		delay := strings.TrimSpace(props.Delay)
		if delay == "" {
			delay = "350ms"
		}
		trigger = "load delay:" + delay
	}

	attrs := map[string]string{
		"aria-hidden": "true",
		"class":       strings.TrimSpace("hidden " + props.Props.Class),
		"hx-post":     action,
		"hx-trigger":  trigger,
		"hx-target":   target,
		"hx-swap":     swap,
	}
	if strings.TrimSpace(props.ID) != "" {
		attrs["id"] = strings.TrimSpace(props.ID)
	}
	return element{Tag: "div", Attrs: attrs}
}
