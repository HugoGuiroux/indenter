// Copyright 2014 Guiroux Hugo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON response type.
type Response struct {
	Body  string `json:",omitempty"`
	Error string `json:",omitempty"`
}

// String function for JSON.
func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		log.Print("Bad marshalling:", err)
		return ""
	}

	return string(b)
}

func ipPrint(r *http.Request, v ...interface{}) {
	s := fmt.Sprintf("%s %s %s: ", r.RemoteAddr, r.Method, r.URL)
	b := make([]interface{}, 0, len(v)+1)
	b = append(b, s)
	b = append(b, v...)
	log.Print(b...)
}

// log function wrapper to log all access for access to http server.
// From https://groups.google.com/d/msg/golang-nuts/s7Xk1q0LSU0/vSvGnerlDZ4J
func logHttpAccess(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipPrint(r, "request")
		handler.ServeHTTP(w, r)
	})
}

// sendResponseToClient send a JSON response to the client with body or error.
func sendResponseToClient(w http.ResponseWriter, body string, err string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{body, err})
}
