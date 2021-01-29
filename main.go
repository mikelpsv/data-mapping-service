package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mikelpsv/data-mapping-service/app"
	"github.com/mikelpsv/data-mapping-service/models"
	"log"
	"net/http"
	"os"
)

type AppCfg struct {
	AppAddr    string
	AppPort    string
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
}

var Config AppCfg

func main() {
	ReadEnv()
	routeItems := app.Routes{}
	routeItems = RegisterHandlers(routeItems)
	router := NewRouter(routeItems)

	app.InitDb(Config.DbHost, Config.DbName, Config.DbUser, Config.DbPassword)

	/*
		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "eddrfsd",
			ValInt:      "453",
			Payload:     "",
		}).Store()

		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "456",
			ValInt:      "453",
			Payload:     "",
		}).Store()

		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "756",
			ValInt:      "6",
			Payload:     "",
		}).Store()

		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "345",
			ValInt:      "63",
			Payload:     "",
		}).Store()

		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "345",
			ValInt:      "6",
			Payload:     "",
		}).Store()


		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "756",
			ValInt:      "453",
			Payload:     "",
		}).Store()

		(&models.Mapping{
			NamespaceId: 1,
			KeyId:       2,
			ValExt:      "345",
			ValInt:      "453",
			Payload:     "",
		}).Store()

	*/

	m := models.Mappings{}
	m.GetByExtValue(1, 2, "345")

	m = models.Mappings{}
	m.GetByIntValue(1, 2, "453")

	log.Fatal(http.ListenAndServe(Config.AppAddr+":"+Config.AppPort, router))
}

func ReadEnv() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	Config.AppAddr = os.Getenv("APP_ADDR")
	Config.AppPort = os.Getenv("APP_PORT")

	Config.DbHost = os.Getenv("POSTGRES_HOST")
	Config.DbPort = os.Getenv("POSTGRES_PORT")
	Config.DbName = os.Getenv("POSTGRES_DB")
	Config.DbUser = os.Getenv("POSTGRES_USER")
	Config.DbPassword = os.Getenv("POSTGRES_PASSWORD")
}

func NewRouter(routeItems app.Routes) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routeItems {
		handlerFunc := route.HandlerFunc
		if route.ValidateToken {
			handlerFunc = app.SetMiddlewareAuth(handlerFunc)
		}

		if route.SetHeaderJSON {
			handlerFunc = app.SetMiddlewareJSON(handlerFunc)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(handlerFunc)
	}

	return router
}
