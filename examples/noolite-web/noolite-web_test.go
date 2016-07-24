package main

import (
	"testing"
	"github.com/dedkovd/noolite"
)

func TestSendCommand(t *testing.T) {
	n, err := noolite.DefaultNooliteAdapter()

	if err != nil {
		t.Error(err)
	}

	defer n.Close()

	err = sendCommand(n, "on", 7, 0, 0, 0, 0)

	if err != nil {
		t.Error(err)
	}

	err = sendCommand(n, "", 0, 0, 0, 0, 0)

	if err == nil {
		t.Error("Command was not set expected")
	}

	err = sendCommand(n, "on", -1, 0, 0, 0, 0)

	if err == nil {
		t.Error("Channel was not set expected")
	}

	err = sendCommand(n, "set", 7, 0, 0, 0, 0)

	if err == nil {
		t.Error("Need some value expected")
	}

	err = sendCommand(n, "qwerty", 2, 0, 0, 0, 0)

	if err == nil {
		t.Error("Command not found expected")
	}
}

func TestParseParams(t *testing.T) {
	cmd, ch, v, r, g, b := parseParams("/set/7/45")

	if cmd != "set" {
		t.Error("Command SET exptected")
	}

	if ch != 7 {
		t.Error("Channel 7 expected")
	}

	if v != 45 {
		t.Error("Value 45 expected")
	}

	if r != 0 || g != 0 || b != 0 {
		t.Error("RGB 000 expected")
	}

	cmd, ch, v, r, g, b = parseParams("")

	if cmd != "" || ch != -1 || v != 0 || r != 0 || g != 0 || b != 0 {
		t.Error("Default values expected")
	}
}
