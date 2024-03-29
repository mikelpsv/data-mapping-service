package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mikelpsv/data-mapping-service/app"
	"github.com/mikelpsv/data-mapping-service/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

func RegisterServiceHandlers(routeItems app.Routes) app.Routes {
	routeItems = append(routeItems, app.Route{
		Name:          "GetNamespaces",
		Method:        "GET",
		Pattern:       "/namespaces",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   GetNamespaces,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "CreateNamespace",
		Method:        "POST",
		Pattern:       "/namespaces/{nsName}",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   CreateNamespace,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "GetKeys",
		Method:        "GET",
		Pattern:       "/namespaces/{nsName}/keys",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   GetKeys,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "CreateKey",
		Method:        "POST",
		Pattern:       "/namespaces/{nsName}/keys/{keyName}",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   CreateKey,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "List",
		Method:        "GET",
		Pattern:       "/list",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   ListMappings,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "ListNewSintax",
		Method:        "GET",
		Pattern:       "/map/{nsName}/{keyName}",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   ListMappingsNewSintax,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "StoreMap",
		Method:        "POST",
		Pattern:       "/map/{nsName}/{keyName}",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   StoreMapping,
	})
	return routeItems
}

func GetNamespaces(w http.ResponseWriter, r *http.Request) {
	ns := new(models.Namespace)
	nsList, err := ns.Get()
	if err != nil {
		app.ResponseERROR(w, http.StatusInternalServerError, err)
		return
	}
	app.ResponseJSON(w, http.StatusOK, nsList)
}

func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	var valNs string
	var exist bool

	vars := mux.Vars(r)
	if valNs, exist = vars["nsName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter required"))
		return
	}
	ns := new(models.Namespace)
	ns.Name = valNs
	ns, err := ns.Store()
	if err != nil {
		app.ResponseERROR(w, http.StatusInternalServerError, err)
		return
	}
	app.ResponseJSON(w, http.StatusOK, ns)
}

func GetKeys(w http.ResponseWriter, r *http.Request) {
	var valNs string
	var exist bool
	var err error

	vars := mux.Vars(r)
	if valNs, exist = vars["nsName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter required"))
		return
	}

	nsItem := new(models.Namespace)
	if _, ok := r.Header["Request-By-Name"]; ok {
		nsItem, err = nsItem.FindByName(valNs)
		if err != nil {
			app.ResponseERROR(w, http.StatusNotFound, err)
			return
		}
	} else {
		i, err := strconv.Atoi(valNs)
		if err != nil {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter failed"))
			return
		}
		nsItem, err = nsItem.FindById(int64(i))
		if err != nil {
			app.ResponseERROR(w, http.StatusNotFound, err)
			return
		}
	}

	keys, err := nsItem.GetKeys()
	if err != nil {
		app.ResponseERROR(w, http.StatusInternalServerError, err)
	}
	app.ResponseJSON(w, http.StatusOK, keys)
}

func CreateKey(w http.ResponseWriter, r *http.Request) {
	var valNs string
	var exist bool
	var err error

	vars := mux.Vars(r)
	if valNs, exist = vars["nsName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter required"))
		return
	}
	if vars["keyName"] == "" {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("keyname parameter required"))
		return
	}

	nsItem := new(models.Namespace)
	if _, ok := r.Header["Request-By-Name"]; ok {
		nsItem, err = nsItem.FindByName(valNs)
		if err != nil {
			app.ResponseERROR(w, http.StatusNotFound, err)
			return
		}
	} else {
		i, err := strconv.Atoi(valNs)
		if err != nil {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter failed"))
			return
		}
		nsItem, err = nsItem.FindById(int64(i))
		if err != nil {
			app.ResponseERROR(w, http.StatusNotFound, err)
			return
		}
	}

	isExist := nsItem.KeyExists(vars["keyName"])
	if isExist {
		app.ResponseERROR(w, http.StatusBadRequest, fmt.Errorf("key exists"))
		return
	}

	key, err := nsItem.CreateKey(vars["keyName"])
	if err != nil {
		app.ResponseERROR(w, http.StatusInternalServerError, err)
	}
	app.ResponseJSON(w, http.StatusOK, key)
}

func ListMappings(w http.ResponseWriter, r *http.Request) {

	var valExt, valInt string
	var rel int

	keys := r.URL.Query()

	valNs := keys.Get("namespace")
	valKey := keys.Get("key")

	// только связанные
	valRelated := keys.Get("rel")
	if valRelated == "1" {
		rel = 1
	} else if valRelated == "0" {
		rel = 0
	} else {
		rel = -1
	}

	if valExt = keys.Get("val_ext"); valExt == "" {
		valExt = keys.Get("id_external") // old syntax
	}
	if valInt = keys.Get("val_int"); valInt == "" {
		valInt = keys.Get("id_internal") // old syntax
	}

	ns := new(models.Namespace)
	_, err := ns.FindByName(valNs)
	if valNs != "" && err != nil {
		app.ResponseERROR(w,
			http.StatusNotFound,
			errors.New(fmt.Sprintf("Namespace %s not found", valNs)))
	}

	key := new(models.Key)
	_, err = key.FindByName(valKey)
	if valKey != "" && err != nil {
		app.ResponseERROR(w,
			http.StatusNotFound,
			errors.New(fmt.Sprintf("key %s not found", valKey)))
	}

	mappings := models.Mappings{}
	mappings.ListMappings(ns, key, valInt, valExt, rel)

	response := new(models.MappingsResponse)
	response.Mappings = mappings

	app.ResponseJSON(w, http.StatusOK, response)
	return

}
func ListMappingsNewSintax(w http.ResponseWriter, r *http.Request) {
	var valNs, valKey string
	var exist bool
	var valExt, valInt string
	var rel int

	vars := mux.Vars(r)

	if valNs, exist = vars["nsName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter required"))
		return
	}
	if valKey, exist = vars["keyName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("key parameter required"))
		return
	}

	keys := r.URL.Query()
	// только связанные
	valRelated := keys.Get("rel")
	if valRelated == "1" {
		rel = 1
	} else if valRelated == "0" {
		rel = 0
	} else {
		rel = -1
	}

	if valExt = keys.Get("val_ext"); valExt == "" {
		valExt = keys.Get("id_external") // old syntax
	}
	if valInt = keys.Get("val_int"); valInt == "" {
		valInt = keys.Get("id_internal") // old syntax
	}

	ns := new(models.Namespace)
	_, err := ns.FindByName(valNs)
	if valNs != "" && err != nil {
		app.ResponseERROR(w,
			http.StatusNotFound,
			errors.New(fmt.Sprintf("Namespace %s not found", valNs)))
	}

	key := new(models.Key)
	_, err = key.FindByName(valKey)
	if valKey != "" && err != nil {
		app.ResponseERROR(w,
			http.StatusNotFound,
			errors.New(fmt.Sprintf("key %s not found", valKey)))
	}

	mappings := models.Mappings{}
	mappings.ListMappings(ns, key, valInt, valExt, rel)

	response := new(models.MappingsResponse)
	response.Mappings = mappings

	app.ResponseJSON(w, http.StatusOK, response)
	return

}

func StoreMapping(w http.ResponseWriter, r *http.Request) {
	/*
		{
		  "mappings": [
		    {
		      "val_ext": "000123",
		      "val_int": "321",
		      "payload": "{'inn': '12345'}"
		    }
		  ]
		}
	*/
	var valNs, valKey string
	var exist bool

	vars := mux.Vars(r)

	if valNs, exist = vars["nsName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("namespace parameter required"))
		return
	}
	if valKey, exist = vars["keyName"]; !exist {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("key parameter required"))
		return
	}

	ns := models.Namespace{}
	_, err := ns.FindByName(valNs)
	if err != nil {
		app.ResponseERROR(w, http.StatusNotFound, errors.New(fmt.Sprintf("namespace %s not found", valNs)))
		return
	}

	key := models.Key{}
	_, err = key.FindByName(valKey)
	if err != nil {
		app.ResponseERROR(w, http.StatusNotFound, errors.New(fmt.Sprintf("key %s not found", valKey)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.ResponseERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	mapReq := models.MappingRequest{}
	err = json.Unmarshal(body, &mapReq)
	if err != nil {
		app.ResponseERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	for _, item := range mapReq.Mappings {
		item.NamespaceId = ns.Id
		item.KeyId = key.Id
		item.Store() // тут может быть ошибка неуникальной вставки
	}

	app.ResponseJSON(w, http.StatusOK, "")

}
