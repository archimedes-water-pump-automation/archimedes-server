package http

import (
	"fmt"
	"net/http"
)

// healthCheckHandler reports process liveness only; it does not check the
// database pool or MQTT connections. The method check is redundant with
// the "GET /health" mux pattern, kept as a guard in case the route is ever
// registered under a broader pattern.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
