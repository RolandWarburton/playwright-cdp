package routes

import (
	"errors"
)

// factory for creating routes
func GetRoute(sw string, middleware Middleware) (IRoute, error) {
	switch sw {
	case "create":
		route := NewBaiscRoute().WithPath("/create").WithMiddleware(middleware)
		return route, nil
	case "example":
		route := NewBaiscRoute().WithPath("/example").WithMiddleware(middleware)
		return route, nil
	}

	return nil, errors.New("wrong route specified")
}
