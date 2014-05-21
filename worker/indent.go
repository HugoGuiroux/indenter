// Copyright 2014 Guiroux Hugo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
