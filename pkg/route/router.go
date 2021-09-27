package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router 路由对象
var Router *mux.Router

func RouteName2URL(routeName string, pairs ...string) string {
	url, err := Router.GetRoute(routeName).URL(pairs...)
	if err != nil {
		return ""
	}
	return url.String()
}

func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
