package main

import (
	"flag"
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
	APP_ADDR    string
	APP_PORT    string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
}

var Config AppCfg

func main() {
	var pFlagInstall, pFlagTestData bool

	flag.BoolVar(&pFlagInstall, "install", false, "Install database. Data will be lost!")
	flag.BoolVar(&pFlagTestData, "testdata", false, "Load test data. Use with a key -install")
	flag.Parse()

	ReadEnv()
	routeItems := app.Routes{}
	routeItems = RegisterHandlers(routeItems)
	router := NewRouter(routeItems)

	if pFlagInstall {
		/*
			_, err := os.Stat(".data/authdata.db")
			if err == nil {
				log.Println("Data file is exist!")
				os.Remove(".data/authdata.db")
			}

		*/
	}

	app.InitDb(Config.DB_HOST, Config.DB_NAME, Config.DB_USER, Config.DB_PASSWORD)

	if pFlagInstall {
		app.Install(pFlagTestData)
	}

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

	log.Fatal(http.ListenAndServe(Config.APP_ADDR+":"+Config.APP_PORT, router))
}

func ReadEnv() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	Config.APP_ADDR = os.Getenv("APP_ADDR")
	Config.APP_PORT = os.Getenv("APP_PORT")

	Config.DB_HOST = os.Getenv("POSTGRES_HOST")
	Config.DB_PORT = os.Getenv("POSTGRES_PORT")
	Config.DB_NAME = os.Getenv("POSTGRES_DB")
	Config.DB_USER = os.Getenv("POSTGRES_USER")
	Config.DB_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
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
