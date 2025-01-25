package gateway

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func ProxyHandler(envVar string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := os.Getenv(envVar)
		if baseURL == "" {
			http.Error(w, "Service URL not set: "+envVar, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		rest := vars["rest"]
		if !strings.HasPrefix(rest, "/") {
			rest = "/" + rest
		}
		fullURL := baseURL + rest

		outReq, err := http.NewRequest(r.Method, fullURL, r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		outReq.Header = r.Header.Clone()

		client := &http.Client{}
		resp, err := client.Do(outReq)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		w.Write(body)
	}
}
