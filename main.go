// main.go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/icholy/digest"
)

func main() {
	// 1. Parse command-line arguments
	username := flag.String("username", "", "Digest Auth Username")
	password := flag.String("password", "", "Digest Auth Password")
	backendURLStr := flag.String("backend", "", "Backend URL (e.g., http://backend-server:8080)")
	listenHost := flag.String("listen-host", "", "Host to listen on (default: all interfaces)")
	listenPort := flag.String("listen-port", "8080", "Port to listen on (default: 8080)")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	var logLevel slog.Level
	if *debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	if *username == "" || *password == "" || *backendURLStr == "" {
		fmt.Fprintln(os.Stderr, "usage: digest-auth-removal-proxy -username <user> -password <pass> -backend <url>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 2. Parse the backend URL
	backendURL, err := url.Parse(*backendURLStr)
	if err != nil {
		logger.Debug("Invalid backend URL", "error", err)
		os.Exit(1)
	}

	// 3. Create a digest transport with the provided credentials
	digestTransport := &digest.Transport{
		Username:  *username,
		Password:  *password,
		Transport: http.DefaultTransport,
	}

	// 4. Create the reverse proxy
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			originalURL := req.URL.String()

			// Set the target URL for the backend
			req.URL.Scheme = backendURL.Scheme
			req.URL.Host = backendURL.Host
			req.URL.Path = singleJoiningSlash(backendURL.Path, req.URL.Path)
			if backendURL.RawQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = backendURL.RawQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = backendURL.RawQuery + "&" + req.URL.RawQuery
			}
			// Set the Host header to the backend's host
			req.Host = backendURL.Host

			logger.Debug("Director: Proxying to backend", "method", req.Method, "url", req.URL.String(), "originalUrl", originalURL)
		},
		// Use the custom client with digest auth for making requests to the backend
		Transport: digestTransport,
	}

	// 5. Define the handler function for all incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handler: Received request", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)
		// Let the reverse proxy handle the request
		reverseProxy.ServeHTTP(w, r)
	})

	// 6. Start the HTTP server on configurable host and port
	listenAddr := fmt.Sprintf("%s:%s", *listenHost, *listenPort)
	logger.Info("Starting reverse proxy", "addr", listenAddr, "backend", backendURL.String(), "username", *username)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// singleJoiningSlash is a helper function copied from the standard library's
// httputil implementation to correctly join paths.
func singleJoiningSlash(a, b string) string {
	aslash := len(a) > 0 && a[len(a)-1] == '/'
	bslash := len(b) > 0 && b[0] == '/'
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
