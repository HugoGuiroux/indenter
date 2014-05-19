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
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

type IndentRequest struct{}

// indentCode function take a source code to indent and launch the gofmt on it.
func indentCode(body string) (string, error) {
	commandStarted := new(bool)
	*commandStarted = false

	deferSafe := func(c io.Closer) {
		if *commandStarted {
			c.Close()
		}
	}

	// Launch without arg, file content will be passed using stdin.
	cmd := exec.Command("gofmt")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer deferSafe(stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer deferSafe(stderr)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	bodyByte := []byte(body)

	// Write to standard pipe.
	if _, err := stdin.Write(bodyByte); err != nil {
		return "", err
	}

	// Close stdin in order to work.
	stdin.Close()

	// Get output.
	var output []byte
	output, err = ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	// Get error.
	var errput []byte
	errput, err = ioutil.ReadAll(stderr)
	errstring := string(errput)

	// Now ack the end.
	if err := cmd.Wait(); err != nil {
		if len(errstring) != 0 {
			log.Print(len(errput))
			return "", errors.New(errstring)
		}

		return "", err
	}

	// Wait close stdin/err/out so trick the defer.
	*commandStarted = true
	if len(errstring) != 0 {
		log.Print(len(errput))
		return "", errors.New(errstring)
	}

	return string(output), nil
}

// Indent function is the RPC function to be called by a frontend.
func (s *IndentRequest) Indent(body string, reply *string) error {
	log.Print("Getting RPC call for indenting source code")

	var err error
	*reply, err = indentCode(body)

	if err != nil {
		log.Print("Error indenting source code: ", err)
	}

	return err
}
