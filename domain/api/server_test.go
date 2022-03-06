package api_test

import (
	"fmt"
	"github.com/lripardo/lrw/domain/api"
	"testing"
)

func TestRoute_Append(t *testing.T) {
	root := "/api/v1"
	route := api.NewRootRoute(root)
	if route.Path != root {
		t.Fatalf("route must have a path with %s value", root)
	}
	other := "other"
	route2 := route.Append(api.Route{
		Path: other,
	})
	path := fmt.Sprintf("%s/%s", root, other)
	if route2.Path != path {
		t.Fatalf("route2 must have other appended to root path")
	}
	route3 := route2.Append(api.Route{
		Methods: []string{"POST", "GET", "DELETE", "PUT"},
	})
	if route3.Path != path || len(route3.Methods) != 4 {
		t.Fatal("route3 must have 4 methods and keep the path")
	}
	route4 := route3.Append(api.Route{
		Methods: []string{"PATCH"},
	})
	if route4.Path != path || len(route4.Methods) != 5 {
		t.Fatal("route4 must have 5 methods and keep the path")
	}
	route5 := route4.Append(api.Route{Handlers: []api.Handler{nil, nil, nil}})
	if route5.Path != path || len(route5.Methods) != 5 || len(route5.Handlers) != 3 {
		t.Fatal("route5 must have 5 methods, 3 handlers and keep the path")
	}
	route6 := route5.Append(api.Route{
		Handlers: []api.Handler{nil},
	})
	if route6.Path != path || len(route6.Methods) != 5 || len(route6.Handlers) != 4 {
		t.Fatal("route6 must have 5 methods, 4 handlers and keep the path")
	}
}
