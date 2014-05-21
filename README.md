# Indenter

Indenter is a small project that can be seen as a proof of concept of a
scalable, fault-tolerant, distributed web service copying the go playground
architecture.

![playground schema](http://blog.golang.org/playground/overview.png)

## Problem
I wanted to know how to build scalable and fault tolerant web application. This
start as an experiment and Indenter is a simple example showing an example of
such application.

The application is divided into three parts: the front end, serving static
content (as well as light dynamic content), the back end (a.k.a. worker, which
handle heavy tasks, here launching gofmt command) and the service discovery.

## Front End
Front End both scale vertically and horizontally. Thanks to goroutines, front
end scales vertically (adding resources to a single node allows to handle more
requests). But it also scales horizontally (adding new nodes) using some router
in front of the front end. Such router can be a simple DNS round-robin
strategy. This also enable usage of elastic cloud technique to adapt front end
to the workload.

The front end is in charge of serving simple pages (here the static index.html
page). It can also handle more dynamic pages which does not require heavy
computing.
When it receives a request for indenting a file (POST request), it contacts the
service discovery daemon (here etcd) to get a list of currently available
workers. Then, it choose given a scheduling policy (i.e. pick random one) a
worker and asks the indenting of the file using RPC request. The result of this
request is sent as result for the client.

## Worker
A worker is in charge of handling RPC requests from front end entities. Again, at
each request a goroutine serves the RPC client, allowing vertical scaling
(internally it uses net/http). When a worker starts, it registers itself to the
service discovery daemon with a random unique identifier and his (host, port)
tuple where a front end can contact it. The entry inside the service discovery
daemon has a time to live (ttl) of 20 seconds. Then every 20 seconds a goroutine
is responsible to register again to the daemon, allowing to catch workers
failures.

## Service discovery
To allow front ends entities to find workers to communicate with, a service
discovery daemon is used. The chosen one is
[etcd](http://coreos.com/docs/distributed-configuration/getting-started-with-etcd). etcd
is a "highly-available key value store for shared configuration and service
discovery". I chose this one as it is distributed, handles failures, is
lightweight and written in Go.

All workers register themselves inside a /workers/ directory. This directory is
fetched whenever a front end needs the list of available workers.

## To go further...
I wanted to let Indenter as simple as possible in order to be easily reusable by
any other projects (the code is under BSD license).

Some possible improvements/usages:

* Adding (distributed) caching to front end to avoid indenting/processing common tasks (as on the playground)
* Front End can handle small dynamic tasks like account management as a real web-site
* Worker can handle other things like launching/compiling an application, sending email, etc.

# How to use it
Simply install etcd (via package manager or
[source](https://github.com/coreos/etcd)) and launch it (./etcd).

Then
```sh
go get github.com/GHugo/indenter/frontend/
go get github.com/GHugo/indenter/worker/
```

The frontend binary waits the index.html file inside the current directory (a
provided one is frontend/index.html).
And run the front end (./frontend) and any number of workers (./worker) you want.
