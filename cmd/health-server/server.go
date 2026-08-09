package healthserver

import (
	"archimedes-worker/core/log"
	"fmt"
	"net/http"
)

func Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprint(w, "ok")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Log("Server is running on http://localhost:8080/health")
	server.ListenAndServe()
}
