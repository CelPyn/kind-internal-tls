package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Body        map[string]interface{} `json:"body"`
	Headers     map[string][]string    `json:"headers"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Proto       string                 `json:"proto"`
	QueryParams map[string][]string    `json:"query_params"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &body)

		response := Response{
			Body:        body,
			Headers:     r.Header,
			Method:      r.Method,
			Path:        r.URL.Path,
			Proto:       r.Proto,
			QueryParams: r.URL.Query(),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("proto", r.Proto).
			Interface("headers", r.Header).
			Interface("query_params", r.URL.Query()).
			Interface("body", body).
			Msg("Incoming request")
	})

	go func() {
		server := http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		log.Info().Msg("Listening on :8080")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	cert, err := tls.LoadX509KeyPair("/etc/tls/tls.crt", "/etc/tls/tls.key")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load TLS certificate")
	}

	server := http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
	}

	log.Info().Msg("Listening on :8443")
	err = server.ListenAndServeTLS("/etc/tls/tls.crt", "/etc/tls/tls.key")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
