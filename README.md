# Internet Simulation Scenarios

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/iss)](https://pkg.go.dev/github.com/bassosimone/iss) [![Build Status](https://github.com/bassosimone/iss/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/iss/actions) [![codecov](https://codecov.io/gh/bassosimone/iss/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/iss)

Internet Simulation Scenarios (`iss`) is a Go package that provides
a declarative network integration testing environment built on top of
[uis](https://github.com/bassosimone/uis) (userspace internet simulation).

Given a `Scenario` describing DNS servers, HTTP servers, and a client
stack, `MustNewSimulation` creates a complete simulated internet with
working DNS resolution, TLS certificates, and HTTP/HTTPS servers.

A `Router` controls packet delivery between hosts. The `DefaultRouter`
supports a swappable `PacketFilter` for simulating censorship conditions.

Use `ScenarioV4` for a ready-made IPv4 topology, or build a custom
`Scenario` for specific test needs.

Basic usage is like:

```go
import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/bassosimone/iss"
)

// Create a new simulation using the default IPv4 scenario
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

router := iss.NewDefaultRouter()
sim := iss.MustNewSimulation(ctx, "testdata", iss.ScenarioV4(), router)
defer func() {
	cancel()
	sim.Wait()
}()

// Resolve the www.example.com domain name
addrs, _ := sim.LookupHost(ctx, "www.example.com")
fmt.Printf("%+v\n", addrs)

// Get the https://www.example.com/ URL
txp := &http.Transport{
	DialContext:     sim.DialContext,
	TLSClientConfig: &tls.Config{RootCAs: sim.CertPool()},
}
clnt := &http.Client{Transport: txp}
resp, _ := clnt.Get("https://www.example.com/")
```

The [example_test.go](example_test.go) file shows a complete example.

## Installation

```sh
go get github.com/bassosimone/iss
```

## Development

To run the tests:

```sh
go test -v .
```

To measure test coverage:

```sh
go test -v -cover .
```

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```
