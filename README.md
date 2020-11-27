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

 - Must set `targetBurstCapacity: -1` on Knative Service (because scraping will this QP implementation doesn't support scraping).
 - No actual health check implemented here, all the healthchecks immediately return true.


