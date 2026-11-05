package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler returns a liveness handler: it always responds 200 as long as
// the process is running and able to serve HTTP. It performs no dependency
// checks, which is the correct semantics for a Kubernetes liveness probe.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// ReadinessFunc is a dependency check used by the readiness handler. It should
// return nil when the dependency is reachable. Keeping it as a function keeps
// the handler decoupled from the concrete MongoDB driver types.
type ReadinessFunc func(ctx context.Context) error

// ReadinessHandler returns a readiness handler. When check is nil it behaves
// like the liveness handler. Otherwise it runs the check with a short timeout
// and responds 503 if the dependency is not reachable.
func ReadinessHandler(check ReadinessFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check == nil {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := check(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"error":  err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
