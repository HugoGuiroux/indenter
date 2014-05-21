// Copyright 2014 Guiroux Hugo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main for worker is in charge to take
// RPC call on Indent method by indenting the file.
// It also register itself to the etcd server discovery to be found by the frontend component.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// RPC bind argument.
var (
	rpcBind = flag.String("rpc_bind", "0.0.0.0", "Network interface address on which listen for RPC requests")
	rpcPort = flag.Int("rpc_port", 54321, "Network port on which listen for RPC requests")
)

// main function of the Worker package.
func main() {
	// Listen for RPC request.
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *rpcBind, *rpcPort))
	if err != nil {
		log.Fatal("Unable to listen: ", err)
	}

	// Announce to service discovery.
	go serveAnouncement()

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Unable to serve: ", err)
	}
}

func init() {
	// Register RPC argument.
	if err := rpc.Register(new(IndentRequest)); err != nil {
		log.Fatal("Error while registering rpc request: ", err)
	}

	// Do it the HTTP way.
	rpc.HandleHTTP()

	// Parse arguments.
	flag.Parse()
}
