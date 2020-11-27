![roo-proxy](docs/roo.jpg)

# roo-proxy

Some experiments with using rust/envoy-plugins for the Knative Queue Proxy

# Experiments

## Identity

Literally just a reverse proxy. Does nothing else. Just to see what breaks.

This is an example of the most minimal QP that will start work with knative.

Requirements and Limitations:

 - Must set targetBurstCapacity=-1 on Knative Service (because scraping will not work with this QP).
 - No actual health check implemented here, all the healthchecks immediately return true.

Install:

```shell
ko apply -f config/identity
```
