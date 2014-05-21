// Copyright 2014 Guiroux Hugo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var (
	// etcd options.
	etcdHost               = flag.String("etcd_host", "localhost", "etcd host name")
	etcdPort               = flag.Int("etcd_port", 4001, "etcd host port")
	etcdAnnouncementPeriod = flag.Int("etcd_period", 20, "etcd announcement period for service discovery")

	// host option.
	hostAddr = flag.String("host_addr", "localhost", "Host name to use for system discovery (which IP is used to connect to this worker)")
)

const etcdDirectory = "/workers/"

// register the current worker to etcd directory.
// Must check if a service with the same name does not exists already.
func register(client *etcd.Client, ttl int) error {
	// Generate random name while taken.
	var name string
	var err error

	taken := true
	for taken {
		name = etcdDirectory + strconv.Itoa(rand.Int())
		taken = nameTaken(client, name)

		if taken {
			log.Print("Name ", name, " was already taken (unlikely)")
		}
	}

	log.Print("Registering ", name, " to etcd with ttl ", ttl)

	// Insert key with good ttl.
	_, err = client.Set(name, fmt.Sprintf("%s:%d", *hostAddr, *rpcPort), uint64(ttl))

	return err
}

// nameTaken check if a worker with the same name is not already registered.
func nameTaken(client *etcd.Client, name string) bool {
	resp, _ := client.Get(name, false, false)
	return resp != nil && resp.Node != nil
}

// serveAnouncement regurarly announce itself to the service discovery daemon.
func serveAnouncement() {
	a := time.Duration(*etcdAnnouncementPeriod) * time.Second

	// Init http client each time (lightweight).
	client := etcd.NewClient([]string{
		fmt.Sprintf("http://%s:%d", *etcdHost, *etcdPort),
	})

	// Announce the first time (don't wait period).
	if err := register(client, *etcdAnnouncementPeriod); err != nil {
		log.Print("Error while registering to system discovery: ", err)
	}

	// Use ticker & channel.
	c := time.Tick(a)
	for _ = range c {
		if err := register(client, *etcdAnnouncementPeriod); err != nil {
			log.Print("Error while registering to system discovery: ", err)
		}
	}
}

func init() {
	// Init only at launch.
	rand.Seed(time.Now().UTC().UnixNano())
}
