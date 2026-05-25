package handler

import "net/http"

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}