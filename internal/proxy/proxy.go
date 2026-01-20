// Package proxy handles proxying requests from the gateway to backend services.
package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Handler manages HTTP proxy operations with an HTTP client.
type Handler struct {
	client *http.Client
}

// NewHandler creates a new proxy handler with configured timeouts and redirect behavior.
func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// ProxyRequest creates an HTTP handler that forwards requests to the target URL.
// It strips the specified prefix from the request path if provided.
func (h *Handler) ProxyRequest(targetURL, stripPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse(targetURL)
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		proxyPath := r.URL.Path
		if stripPrefix != "" {
			proxyPath = strings.TrimPrefix(proxyPath, stripPrefix)
		}

		targetURL := target.Scheme + "://" + target.Host + proxyPath
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			log.Printf("Error creating proxy request: %v", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		proxyReq.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
		proxyReq.Header.Set("X-Forwarded-Host", r.Host)

		resp, err := h.client.Do(proxyReq)
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}
}
