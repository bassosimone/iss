// SPDX-License-Identifier: GPL-3.0-or-later

package iss_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bassosimone/iss"
	"github.com/bassosimone/runtimex"
	"github.com/bassosimone/uis"
)

func Example() {
	// Create a simulation with the default IPv4 scenario
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := iss.NewDefaultRouter()
	sim := iss.MustNewSimulation(ctx, "testdata", iss.ScenarioV4(), router)
	defer func() {
		cancel()
		sim.Wait()
	}()

	// Start a PCAP trace to capture all packets for debugging.
	// The resulting file can be inspected with Wireshark or tcpdump.
	filep := runtimex.PanicOnError1(os.Create("testdata/example.pcap"))
	pcaptrace := uis.NewPCAPTrace(filep, uis.MTUJumbo)
	defer pcaptrace.Close()

	// Set a packet filter that captures all packets to the trace
	router.SetPacketFilter(iss.PacketFilterFunc(func(pkt uis.VNICFrame) bool {
		pcaptrace.Dump(pkt.Packet)
		return false
	}))

	// Resolve www.example.com
	addrs := runtimex.PanicOnError1(sim.LookupHost(ctx, "www.example.com"))
	fmt.Printf("%+v\n", addrs)

	// Fetch www.example.com
	txp := &http.Transport{
		DialContext:       sim.DialContext,
		ForceAttemptHTTP2: true,
		TLSClientConfig:   &tls.Config{RootCAs: sim.CertPool()},
	}
	clnt := &http.Client{Transport: txp}
	hr := runtimex.PanicOnError1(clnt.Get("https://www.example.com/"))
	defer hr.Body.Close()

	fmt.Printf("%d\n", hr.StatusCode)

	body := runtimex.PanicOnError1(io.ReadAll(hr.Body))
	fmt.Printf("%d\n", len(body))

	// Output:
	// [104.18.26.120]
	// 200
	// 605
}
