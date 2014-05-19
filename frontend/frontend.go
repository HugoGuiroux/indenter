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

// main package for frontend application.
// The frontend application is in charge to respond for the static index page
// (containing form to ask for indentation). It is also responsible for sending
// the requests for indenting file to one worker using service discovery etcd.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Retrieves from command line argument.
var (
	httpBind = flag.String("http_bind", "0.0.0.0", "Network interface address on which listen for HTTP requests")
	httpPort = flag.Int("http_port", 1234, "Network port on which listen for HTTP requests")
)

// Template caching.
var templates *template.Template

// indexPageHandle only load a template to serve the static index page (KISS)
// This page contains a form to send the source file to indent and receive back
// the JSON formatted response.
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ipPrint(r, "Error sending static response", err)
	}
}

// requestPageHandler get from body parameter the content of a file needed to be
// indented. Then it follows this simple algorithm:
// - find one worker using service discovery (etcd request on special array).
// - contact this worker.
// - perform RPC request where function prototype is known.
// - send back the result (or error) to client using JSON formated output.
func requestPageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve body.
	body := r.FormValue("body")
	if body == "" {
		ipPrint(r, "Request body empty")
		sendResponseToClient(w, "", "Request body empty")
		return
	}

	// Find worker.
	addr := getWorkerAddrFromServiceDiscovery()
	if addr == "" {
		ipPrint(r, "No worker found")

		// Send error to client.
		sendResponseToClient(w, "", "No worker found")
		return
	}

	// Contact it.
	worker, err := contactWorker(addr)
	if err != nil {
		ipPrint(r, "Error while contacting worker:", err)
		sendResponseToClient(w, "", "Error while contacting worker")
		return
	}

	// Perform RPC request.
	var result string
	if result, err = performIndent(worker, body); err != nil {
		ipPrint(r, "RPC error: ", err)
		sendResponseToClient(w, "", err.Error())
		return
	}

	// Send back result.
	sendResponseToClient(w, result, "")
}

func main() {
	conn := fmt.Sprintf("%s:%d", *httpBind, *httpPort)

	// Listen for web server.
	log.Println("Listening for request on ", conn)
	if err := http.ListenAndServe(conn, logHttpAccess(http.DefaultServeMux)); err != nil {
		log.Fatal("Error listening: ", err)
	}
}

func init() {
	// Setup routing.
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/request", requestPageHandler)

	// Get command line argument.
	flag.Parse()

	// Cache template.
	var err error
	if templates, err = template.ParseFiles("index.html"); err != nil {
		log.Fatal("Error caching template: ", err)
	}
}
