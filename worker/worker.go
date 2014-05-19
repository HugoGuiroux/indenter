// Copyright Hugo Guiroux.
// This file is part of Indenter.
//
// Indenter is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
