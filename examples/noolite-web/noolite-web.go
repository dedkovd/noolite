/*
Copyright 2016 Denis V. Dedkov (denis.v.dedkov@gmail.com)

This file is part of Noolite Go bindings.

Noolite Go bindings is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Noolite Go bindings is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Noolite Go bindings.  If not, see <http://www.gnu.org/licenses/>.
*/

// Example WEB-server for control noolite adapter
//
// Default run on :8080 with static dir /var/www/static
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dedkovd/noolite"
	"net/http"
	"strconv"
	"strings"
)

func sendCommand(n *noolite.NooliteAdapter, command string, channel, value, r, g, b int) error {
	if channel == -1 {
		return errors.New("Channel was not set")
	}

	if command == "" {
		return errors.New("Command was not set")
	}

	if command == "set" {
		if value != 0 {
			return n.SetBrightnesValue(channel, value)
		} else if r != 0 || g != 0 || b != 0 {
			return n.SetBrightnesValues(channel, r, g, b)
		} else {
			return errors.New("Need some value")
		}
	} else {
		cmd, ok := n.FindCommand(command)

		if !ok {
			return errors.New("Command not found")
		}

		return cmd(channel)
	}
}

func parseParams(path string) (string, int, int, int, int, int) {
	params := strings.Split(path, "/")[1:]

	command := ""
	channel := -1
	value := 0
	r := 0
	g := 0
	b := 0

	command = params[0]
	if len(params) > 1 {
		channel, _ = strconv.Atoi(params[1])
	}
	if len(params) > 2 {
		value, _ = strconv.Atoi(params[2])
	}
	if len(params) == 5 {
		value = 0
		r, _ = strconv.Atoi(params[2])
		g, _ = strconv.Atoi(params[3])
		b, _ = strconv.Atoi(params[4])
	}

	return command, channel, value, r, g, b
}

func main() {
	binding := *flag.String("bind", ":8080", "Address binding")
	static_dir := *flag.String("static", "/var/www/static", "Static directory")

	flag.Parse()

	n, err := noolite.DefaultNooliteAdapter()

	if err != nil {
		panic(err)
	}

	defer n.Close()

	http.HandleFunc("/noolite/", func(w http.ResponseWriter, r *http.Request) {
		command, channel, value, red, green, blue := parseParams(r.URL.Path[1:])

		err := sendCommand(n, command, channel, value, red, green, blue)

		if err != nil {
			fmt.Fprintf(w, "{\"error\": %q}", err)
		} else {
			fmt.Fprintf(w, "{\"command\": %q, \"channel\": \"%d\"}", command, channel)
		}
	})

	fs := http.FileServer(http.Dir(static_dir))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	panic(http.ListenAndServe(binding, nil))
}
