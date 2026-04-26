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
	RedirectTo     string
	SuccessMessage string
}

func Redirect(to string, successMessage string) ActionResult {
	return ActionResult{
		RedirectTo:     to,
		SuccessMessage: successMessage,
	}
}
