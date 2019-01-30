package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Pattern).
			HandlerFunc(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"MovieList",
		"GET",
		"/movies",
		MovieList,
	},
	Route{
		"MovieDetail",
		"GET",
		"/movies/{movieId}",
		MovieDetail,
	},
	Route{
		"MovieAdd",
		"POST",
		"/movies",
		MovieAdd,
	},
	Route{
		"MovieUpdate",
		"PUT",
		"/movies/{movieId}",
		MovieUpdate,
	},
	Route{
		"MovieRemove",
		"DELETE",
		"/movies/{movieId}",
		MovieRemove,
	},
}
