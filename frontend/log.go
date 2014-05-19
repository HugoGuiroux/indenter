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
