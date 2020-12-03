package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

// this flag is passed when we're run as an exec probe.
var readinessProbeTimeout = flag.Duration("probe-period", -1, "run readiness probe with given timeout")

// mtls/queue/main is a QP that serves istio TLS certs _without_ using the istio sidecar.
// the certs are expected to be mounted in /etc/certs.
// You can cause istio to add a sidecar _just_ to download and rotate the certs
// (without actually terminating TLS etc) with the following annotations:
//   sidecar.istio.io/inject: "true"
//   sidecar.istio.io/interceptionMode: "NONE"
//   status.sidecar.istio.io/port: "0"
//   sidecar.istio.io/rewriteAppHTTPProbers: "false"
//   proxy.istio.io/config: |
//     proxyMetadata:
//       OUTPUT_CERTS: /etc/istio-certs
//   sidecar.istio.io/userVolume: '[{"name": "istio-certs", "emptyDir": {"medium": "Memory"}}]'
//   sidecar.istio.io/userVolumeMount: '[{"name": "istio-certs", "mountPath": "/etc/istio-certs"}]'
// You also need to add a volumeMount to the queue proxy pod to mount these certs.
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

	time.Sleep(4 * time.Second) // hack to wait for the cert to show up

	log.Fatal(http.ListenAndServeTLS(
		":"+os.Getenv("QUEUE_SERVING_PORT"),
		"/etc/certs/cert-chain.pem",
		"/etc/certs/key.pem",
		logHandler(healthzHandler(proxy))),
	)
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
