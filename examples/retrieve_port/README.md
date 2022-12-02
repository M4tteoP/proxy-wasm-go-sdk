## PoC retriving port Connection attributes
Tested against: envoy version: `15baf56003f33a07e0ab44f82f75a660040db438/1.24.0/Distribution/RELEASE/BoringSSL`
[Envoy Connection attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#connection-attributes)

Attributes:
- `source.port`
- `destination.port`

## How to run the PoC:
```sh
# Go to repo root folder
cd ..
# Build a specific example.
make build.example name=retrieve_port
# Run a specific example.
make run name=retrieve_port
# Let's make a request to trigger callbacks
curl 0.0.0.0:18000
```

docker run -p 8080:80 kennethreitz/httpbin


It works :) 
