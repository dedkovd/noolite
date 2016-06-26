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

package noolite

import "testing"

func TestDefaultNooliteAdapter(t *testing.T) {
	n, err := DefaultNooliteAdapter()

	if err != nil {
		t.Error(err)
	}
	defer n.Close()

	cmd := n.composeCommand(0x00, 0x00)

	if cmd[0] != 0x50 {
		t.Error("Adapter mode (first byte): expected 0x50, got", cmd[0])
	}
}

func TestNewNooliteAdapter(t *testing.T) {
	n, err := NewNooliteAdapter(0, 2, 2)
	if err != nil {
		t.Error(err)
	}
	n.Close()

	n, err = NewNooliteAdapter(8, 2, 2)
	if err == nil {
		t.Error("Bad value for mode. Must be in range 0..7")
	}

	n, err = NewNooliteAdapter(0, 4, 2)
	if err == nil {
		t.Error("Bad value for bitrate. Must be in range 0..3")
	}

	n, err = NewNooliteAdapter(0, 2, 8)
	if err == nil {
		t.Error("Bad value for mode. Must be in range 0..7")
	}
}

type teststr struct {
	values []byte
	cmd    command
	args   []int
}

func TestComposeCommand(t *testing.T) {
	n, err := DefaultNooliteAdapter()

	if err != nil {
		t.Error(err)
	}

	defer n.Close()

	var tests = []teststr{
		{[]byte{0x50, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, off, []int{}},
		{[]byte{0x50, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, decBrightnes, []int{}},
		{[]byte{0x50, 0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, on, []int{}},
		{[]byte{0x50, 0x03, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, incBrightnes, []int{}},
		{[]byte{0x50, 0x04, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, cSwitch, []int{}},
		{[]byte{0x50, 0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, invertBrightnes, []int{}},
		{[]byte{0x50, 0x06, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00}, set, []int{1}},
		{[]byte{0x50, 0x06, 0x03, 0x00, 0x01, 0x01, 0x01, 0x01}, set, []int{1, 1, 1}},
		{[]byte{0x50, 0x07, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, callScenario, []int{}},
		{[]byte{0x50, 0x08, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, saveScenario, []int{}},
		{[]byte{0x50, 0x09, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, unbind, []int{}},
		{[]byte{0x50, 0x0a, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, stopColorSelection, []int{}},
		{[]byte{0x50, 0x0f, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, bind, []int{}},
		{[]byte{0x50, 0x10, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, colorSelection, []int{}},
		{[]byte{0x50, 0x11, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, colorSwitch, []int{}},
		{[]byte{0x50, 0x12, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, modeSwitch, []int{}},
		{[]byte{0x50, 0x13, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, effectSpeed, []int{}},
	}

	for _, tstr := range tests {
		cmd := n.composeCommand(tstr.cmd, 1, tstr.args...)
		for i, v := range cmd {
			if v != tstr.values[i] {
				t.Error("For command",
					tstr.cmd,
					"expected",
					tstr.values[i],
					"got", v,
					"in position", i)
			}
		}
	}
}
