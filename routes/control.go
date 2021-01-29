package routes

import (
	"github.com/mikelpsv/data-mapping-service/app"
	"net/http"
)

type PingResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func RegisterControlHandlers(routeItems app.Routes) app.Routes {
	routeItems = append(routeItems, app.Route{
		Name:          "Ping",
		Method:        "GET",
		Pattern:       "/ping",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   Ping,
	})
	return routeItems
}

func Ping(w http.ResponseWriter, r *http.Request) {
	app.ResponseJSON(w, http.StatusOK, PingResponse{
		Code:        http.StatusOK,
		Description: "",
	})
}
