package tools

import (
	"net/http"
)

func HttpAllowCrossDomainAccess(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
