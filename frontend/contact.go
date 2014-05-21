// Copyright 2014 Guiroux Hugo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

// etcd options.
var (
	etcdHost   = flag.String("etcd_host", "localhost", "etcd host name")
	etcdPort   = flag.Int("etcd_port", 4001, "etcd host port")
	etcdClient *etcd.Client
)

const etcdDirectory = "/workers/"

type IndentRequest struct{}

// getWorkerAddrFromServiceDiscovery contact etcd server to get the list of
// available workers.
// Thus we choose randomly one available. Using metrics we could elaborate
// better load balanced strategy.
// return "" if no worker found, otherwise return a string of format ip:port for
// RPC request.
func getWorkerAddrFromServiceDiscovery() string {
	resp, err := etcdClient.Get(etcdDirectory, false, false)
	if err != nil {
		log.Print("Etcd workers retrieval fails: ", err)
		return ""
	}

	l := resp.Node.Nodes.Len()
	if l == 0 {
		return ""
	}

	b := make([]string, 0, l)

	for _, n := range resp.Node.Nodes {
		b = append(b, n.Value)
	}

	res := rand.Intn(len(b))

	return b[res]
}

// contactWorker try to establish RPC connection using the addr paramaeter as
// connection string.
func contactWorker(addr string) (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", addr)
}

// performIndent emit a RPC call for the remote Indent function which take as
// argument the code and return the code indented or an error.
// Warning, this function do the RPC call in a blocking way.
func performIndent(worker *rpc.Client, body string) (s string, e error) {
	e = worker.Call("IndentRequest.Indent", body, &s)

	return s, e
}

func init() {
	// Init only at launch.
	rand.Seed(time.Now().UTC().UnixNano())

	// To be sure args are parsed.
	flag.Parse()

	// Do this only once as it will spawn TCP connection each time and keep them
	// alive.
	etcdClient = etcd.NewClient([]string{
		fmt.Sprintf("http://%s:%d", *etcdHost, *etcdPort),
	})
}
