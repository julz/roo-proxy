package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// this flag is passed when we're run as an exec probe.
var readinessProbeTimeout = flag.Duration("probe-period", -1, "run readiness probe with given timeout")

// identity/queue/main is a QP replacement that literally just reverse proxies to the user container.
// to make the rest of knative happy, it also
// - runs as an exec probe via /ko-app/queue (but just returns success always, sorry).
// - replies to the "K-Network-Probe" header so the activator will route to us
// (it does its own probing separately from the regular kubernetes readiness
// flow). Again this just always says success, again sorry.
func main() {
	flag.Parse()
	// just immediately exit if the flag is passed indicating we're being run
	// as an exec probe since this is supposed to be the most trivial possible QP.
	if *readinessProbeTimeout != -1 {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:" + os.Getenv("USER_PORT"),
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("QUEUE_SERVING_PORT"), logHandler(healthzHandler(proxy))))
}

func healthzHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We need to reply to the activator with "probe" because it does its own
		// probing separately from the kubernetes readiness probing.
		if r.Header.Get("K-Network-Probe") != "" {
			io.WriteString(w, "queue")
			return
		}

		h.ServeHTTP(w, r)
	}
}

func logHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle", r.Host, r.URL, r.Header)
		h.ServeHTTP(w, r)
	}
}
