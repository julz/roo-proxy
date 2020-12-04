![roo-proxy](docs/roo.jpg)

# roo-proxy

Some experiments with using rust/envoy-plugins for the Knative Queue Proxy

# Experiments

## Identity

Literally just a reverse proxy. Does nothing else. Just to see what breaks.

This is an example of the most minimal QP that will work with knative.

### Install:

```shell
ko apply -f config/identity
```

### Requirements and Main Limitations:

 - Must set `targetBurstCapacity: -1` on Knative Service (because this QP implementation doesn't support scraping).
 - No actual health check implemented here, all the healthchecks immediately return true.

## QP TLS Termination

Uses an envoy sidecar to put the correct keys on the filesystem for mTLS, but
does the actual TLS Termination in the Queue Proxy. This means (a) don't need
any istio iptables/init-container/CNI stuff, but you still get mTLS from
activator/ingress to QP, (b) no envoy in the routing path, so one less hop, and
no waiting for endpoints to propagate before containers can route.

### Install:

```shell
ko apply -f config/mtls
```

### Requirements and Main Limitations:

 - Must set `targetBurstCapacity: -1` on Knative Service (because this QP implementation doesn't support scraping).
 - No actual health check implemented here, all the healthchecks immediately return true.
 - No outgoing TLS from user container (since there's no interception of
   outgoing connections). Could mount the certificate in to the user container, though.
 - Need to turn on istio sidecar injection, and add various annotations to the ksvc. Also hack
   knative resource/queue.go to add volume mount to the pod.
   See [knative WIP branch](https://github.com/knative/serving/compare/master...julz:qpmtls?expand=1).

