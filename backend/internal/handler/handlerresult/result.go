package handlerresult

import "github.com/romsar/gonertia/v3"

// Page
type PageResult struct {
	Component string
	Props     gonertia.Props
}

func Page(component string, props gonertia.Props) PageResult {
	return PageResult{
		Component: component,
		Props:     props,
	}
}

// Action
type ActionResult struct {
	RedirectTo string
}

func Redirect(redirectTo string) ActionResult {
	return ActionResult{
		RedirectTo: redirectTo,
	}
}
