// Copyright Hugo Guiroux
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
