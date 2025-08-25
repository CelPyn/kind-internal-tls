package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var customCertFile = "/usr/local/private-ca/bundle.pem"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	destination := os.Getenv("DESTINATION")
	if destination == "" {
		log.Fatal().Msg("DESTINATION must be set")
	}

	log.Info().Str("destination", destination).Msg("Starting client")
	client := setupClient()

	go serveProbes()

	for {
		<-time.After(5 * time.Second)
		go callHTTP(destination, client)
		go callHTTPS(destination, client)
	}
}

func setupClient() *http.Client {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load system cert pool")
	}

	customCert, err := os.ReadFile(customCertFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read custom cert file")
	}
	if ok := certPool.AppendCertsFromPEM(customCert); !ok {
		log.Fatal().Msg("Failed to append custom cert to cert pool")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	return &http.Client{
		Transport: transport,
	}
}

func serveProbes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start probe server")
	}
}

func callHTTP(destination string, client *http.Client) {
	log.Info().Msg("Calling HTTP endpoint")
	url := fmt.Sprintf("http://%s:8080/http", destination)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return
	}

	do, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform request")
		return
	}
	defer func() {
		_ = do.Body.Close()
	}()

	log.Info().Str("status", do.Status).Msg("Received response from HTTP endpoint")
}

func callHTTPS(destination string, client *http.Client) {
	log.Info().Msg("Calling HTTPS endpoint")
	url := fmt.Sprintf("https://%s:8443/https", destination)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return
	}

	do, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform request")
		return
	}
	defer func() {
		_ = do.Body.Close()
	}()

	log.Info().Str("status", do.Status).Msg("Received response from HTTPS endpoint")
}
