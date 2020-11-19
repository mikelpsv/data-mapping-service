package app

import (
	"net/http"
)

type Route struct {
	Name          string
	Method        string
	Pattern       string
	SetHeaderJSON bool
	ValidateToken bool
	HandlerFunc   http.HandlerFunc
}

type Routes []Route
