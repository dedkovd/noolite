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
		err := sendCommand(n, parseParams(r.URL.Path[1:]))

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
