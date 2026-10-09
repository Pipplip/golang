package httpapi

import "net/http"

// HomeHandler ..
func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Rest-API HOME"))
}
