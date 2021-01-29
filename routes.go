package main

/**
  Данный файл содержит изменяемую часть сервера Api
  Список методов и функции-обработчики
*/

import (
	"github.com/mikelpsv/data-mapping-service/app"
	"github.com/mikelpsv/data-mapping-service/routes"
)

func RegisterHandlers(routeItems app.Routes) app.Routes {
	routeItems = routes.RegisterControlHandlers(routeItems) // ping etc
	routeItems = routes.RegisterServiceHandlers(routeItems) // ping etc
	return routeItems
}
