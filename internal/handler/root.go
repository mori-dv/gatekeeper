package handler

import "net/http"

func Root(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]any{
		"service": "gatekeeper",
	})
}