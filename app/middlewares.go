package app

import (
	"errors"
	"log"
	"net/http"
)

/*
	Добавляем Content-Type всем обработчикам
*/
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Vary", "Accept-Encoding")
		next(w, r)
	}
}

/*
	Проверяем токен
*/
func SetMiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Validating token...")
		err := TokenValid(r)
		if err != nil {
			ResponseERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
