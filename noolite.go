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

// Package noolite provide class for control Noolite Adapters PC11xx.
// Protocol described on url:
// http://www.noo.com.by/assets/files/software/PC11xx_HID_API.pdf
package noolite

import (
	"errors"
	"github.com/tonymagro/usb"
)

const (
	VID = 5824 // Vendor ID
	PID = 1503 // Product ID
)

// Available commands
type command byte

const (
	off command = iota
	decBrightnes
	on
	incBrightnes
	cSwitch
	invertBrightnes
	set
	callScenario
	saveScenario
	unbind
	stopColorSelection
	bind = iota + 0x04
	// Commands for SD111-180 only
	colorSelection
	colorSwitch
	modeSwitch
	effectSpeed
)

// Noolite Adapter class as USB HID
type NooliteAdapter struct {
	*usb.Device
	mode byte
}

// Return NooliteAdapter object.
//
//  work mode: range 0..7
//  bitrate: range 0..3
//  command repeats: range 0..7
func NewNooliteAdapter(mode, bitrate, repeats uint) (*NooliteAdapter, error) {
	usb.Init()

	d := usb.Open(VID, PID)

	if d == nil {
		return nil, errors.New("Device not found")
	}

	if mode > 7 {
		return nil, errors.New("Mode must be in 0..7 range")
	}

	if bitrate > 3 {
		return nil, errors.New("Bitrate must be in 0..3 range")
	}

	if repeats > 7 {
		return nil, errors.New("Repeats must be in 0..7 range")
	}

	m := (byte(repeats) << 5) | (byte(bitrate) << 3) | byte(mode)

	d.Interface(0)
	if d.LastError() != "No error" {
		defer d.Close()
		return nil, errors.New(d.LastError())
	}

	return &NooliteAdapter{d, m}, nil
}

// Return NooliteAdapter object with default work mode values:
//
//  work mode:  0
//  bitrate: 2 (for 1000 bit/sec)
//  command repeats: 2
func DefaultNooliteAdapter() (*NooliteAdapter, error) { // Default constructor
	return NewNooliteAdapter(0, 2, 2)
}

// Return NooliteAdapter method for string command
//
// Set command must be separately processed because have different signature
func (n *NooliteAdapter) FindCommand(command string) (func(int) error, bool) {
	m := map[string]func(int) error {
		"on": n.On,
		"off": n.Off,
		"switch": n.Switch,
		"decraseBrightnes": n.DecraseBrightnes,
		"incraseBrightnes": n.IncraseBrightnes,
		"invertBrightnes": n.InvertBrightnes,
		"callScenario": n.CallScenario,
		"saveScenario": n.SaveScenario,
		"unbind": n.UnbindChannel,
		"stopColorSelection": n.StopColorSelection,
		"bind": n.BindChannel,
		"colorSelection": n.ColorSelection,
		"colorSwitch": n.ColorSwitch,
		"modeSwitch": n.ModeSwitch,
		"effectSpeed": n.EffectSpeed,
	}

	cmd, ok := m[command]
	return cmd, ok
}

func (n *NooliteAdapter) composeCommand(cmd command, channel int, args ...int) []byte {
	c := make([]byte, 8)

	c[0] = n.mode
	c[1] = byte(cmd)
	c[4] = byte(channel)

	if cmd == set {
		l := len(args)
		switch l {
		case 1:
			{
				c[2] = 0x01
				c[5] = byte(args[0])
			}
		case 3:
			{
				c[2] = 0x03
				for i, v := range args {
					c[5+i] = byte(v)
				}
			}
		default:
			panic("Bad arguments for SET command")
		}
	}

	return c
}

func (n *NooliteAdapter) sendCommand(command []byte) error {
	n.ControlMsg(0x21, 0x09, 0x300, 0, command)
	if n.LastError() != "No error" {
		return errors.New(n.LastError())
	}
	return nil
}

// Turn power OFF for specified channel
func (n *NooliteAdapter) Off(channel int) error {
	cmd := n.composeCommand(off, channel)
	return n.sendCommand(cmd)
}

// Smooth brightnes decrase for specified channel
func (n *NooliteAdapter) DecraseBrightnes(channel int) error {
	cmd := n.composeCommand(decBrightnes, channel)
	return n.sendCommand(cmd)
}

// Turn power ON for specified channel
func (n *NooliteAdapter) On(channel int) error {
	cmd := n.composeCommand(on, channel)
	return n.sendCommand(cmd)
}

// Smooth brightnes incrase for specified channel
func (n *NooliteAdapter) IncraseBrightnes(channel int) error {
	cmd := n.composeCommand(incBrightnes, channel)
	return n.sendCommand(cmd)
}

// Switch power state between off and on for specified channel
func (n *NooliteAdapter) Switch(channel int) error {
	cmd := n.composeCommand(cSwitch, channel)
	return n.sendCommand(cmd)
}

// Smooth brightnes incrase or decrase for specified channel
func (n *NooliteAdapter) InvertBrightnes(channel int) error {
	cmd := n.composeCommand(invertBrightnes, channel)
	return n.sendCommand(cmd)
}

// Set brightnes value for specified channel
//
// Value must be in range 35..155.
// When value == 0 lights off.
// When value > 155 lights on for full brightness.
func (n *NooliteAdapter) SetBrightnesValue(channel, value int) error {
	cmd := n.composeCommand(set, channel, value)
	return n.sendCommand(cmd)
}

// Set brightnes values for independens channels
//
// Available for SD111-180 only
func (n *NooliteAdapter) SetBrightnesValues(channel, val1, val2, val3 int) error {
	cmd := n.composeCommand(set, channel, val1, val2, val3)
	return n.sendCommand(cmd)
}

// Call scenario for specified channel
func (n *NooliteAdapter) CallScenario(channel int) error {
	cmd := n.composeCommand(callScenario, channel)
	return n.sendCommand(cmd)
}

// Save scenario for specified channel
func (n *NooliteAdapter) SaveScenario(channel int) error {
	cmd := n.composeCommand(saveScenario, channel)
	return n.sendCommand(cmd)
}

// Unbind signal for specified channel
func (n *NooliteAdapter) UnbindChannel(channel int) error {
	cmd := n.composeCommand(unbind, channel)
	return n.sendCommand(cmd)
}

// Stop color selection for specified channel
//
// Available for SD111-180 only
func (n *NooliteAdapter) StopColorSelection(channel int) error {
	cmd := n.composeCommand(stopColorSelection, channel)
	return n.sendCommand(cmd)
}

// Set binding for specified channel
func (n *NooliteAdapter) BindChannel(channel int) error {
	cmd := n.composeCommand(bind, channel)
	return n.sendCommand(cmd)
}

// Smooth color changing for specified channel
//
// Stop with StopColorSelection method
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) ColorSelection(channel int) error {
	cmd := n.composeCommand(colorSelection, channel)
	return n.sendCommand(cmd)
}

// Switch color for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) ColorSwitch(channel int) error {
	cmd := n.composeCommand(colorSwitch, channel)
	return n.sendCommand(cmd)
}

// Switch work mode for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) ModeSwitch(channel int) error {
	cmd := n.composeCommand(modeSwitch, channel)
	return n.sendCommand(cmd)
}

// Set change color speed for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) EffectSpeed(channel int) error {
	cmd := n.composeCommand(effectSpeed, channel)
	return n.sendCommand(cmd)
}
