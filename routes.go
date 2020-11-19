package main

/**
  Данный файл содержит изменяемую часть сервера Api
  Список методов и функции-обработчики
*/

import (
	"github.com/mikelpsv/data-mapping-service/app"
)

func RegisterHandlers(routeItems app.Routes) app.Routes {
	//routeItems = routes.RegisterCustomerHandler(routeItems)
	return routeItems
}
