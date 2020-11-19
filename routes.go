package main

/**
  Данный файл содержит изменяемую часть сервера Api
  Список методов и функции-обработчики
*/

import (
	"github.com/mikelpsv/data-mapping-service/app"
	//"github.com/mikelpsv/auth_service/routes"
)

func RegisterHandlers(routeItems app.Routes) app.Routes {
	//routeItems = routes.RegisterCustomerHandler(routeItems)
	return routeItems
}
