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

func (n *NooliteAdapter) sendCommand(command []byte) {
	n.Configuration(1)
	n.Interface(0)
	n.ControlMsg(0x21, 0x09, 0x300, 0, command)
}

// Turn power OFF for specified channel
func (n *NooliteAdapter) Off(channel int) {
	cmd := n.composeCommand(off, channel)
	n.sendCommand(cmd)
}

// Smooth brightnes decrase for specified channel
func (n *NooliteAdapter) DecraseBrightnes(channel int) {
	cmd := n.composeCommand(decBrightnes, channel)
	n.sendCommand(cmd)
}

// Turn power ON for specified channel
func (n *NooliteAdapter) On(channel int) {
	cmd := n.composeCommand(on, channel)
	n.sendCommand(cmd)
}

// Smooth brightnes incrase for specified channel
func (n *NooliteAdapter) IncraseBrightnes(channel int) {
	cmd := n.composeCommand(incBrightnes, channel)
	n.sendCommand(cmd)
}

// Switch power state between off and on for specified channel
func (n *NooliteAdapter) Switch(channel int) {
	cmd := n.composeCommand(cSwitch, channel)
	n.sendCommand(cmd)
}

// Smooth brightnes incrase or decrase for specified channel
func (n *NooliteAdapter) InvertBrightnes(channel int) {
	cmd := n.composeCommand(invertBrightnes, channel)
	n.sendCommand(cmd)
}

// Set brightnes value for specified channel
//
// Value must be in range 35..155.
// When value == 0 lights off.
// When value > 155 lights on for full brightness.
func (n *NooliteAdapter) SetBrightnesValue(channel, value int) {
	cmd := n.composeCommand(set, channel, value)
	n.sendCommand(cmd)
}

// Set brightnes values for independens channels
//
// Available for SD111-180 only
func (n *NooliteAdapter) SetBrightnesValues(channel, val1, val2, val3 int) {
	cmd := n.composeCommand(set, channel, val1, val2, val3)
	n.sendCommand(cmd)
}

// Call scenario for specified channel
func (n *NooliteAdapter) CallScenario(channel int) {
	cmd := n.composeCommand(callScenario, channel)
	n.sendCommand(cmd)
}

// Save scenario for specified channel
func (n *NooliteAdapter) SaveScenario(channel int) {
	cmd := n.composeCommand(saveScenario, channel)
	n.sendCommand(cmd)
}

// Unbind signal for specified channel
func (n *NooliteAdapter) UnbindChannel(channel int) {
	cmd := n.composeCommand(unbind, channel)
	n.sendCommand(cmd)
}

// Stop color selection for specified channel
//
// Available for SD111-180 only
func (n *NooliteAdapter) StopColorSelection(channel int) {
	cmd := n.composeCommand(stopColorSelection, channel)
	n.sendCommand(cmd)
}

// Set binding for specified channel
func (n *NooliteAdapter) BindChannel(channel int) {
	cmd := n.composeCommand(bind, channel)
	n.sendCommand(cmd)
}

// Smooth color changing for specified channel
//
// Stop with StopColorSelection method
// 
// Avialable for SD111-180 only
func (n *NooliteAdapter) ColorSelection(channel int) {
	cmd := n.composeCommand(colorSelection, channel)
	n.sendCommand(cmd)
}

// Switch color for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) ColorSwitch(channel int) {
	cmd := n.composeCommand(colorSwitch, channel)
	n.sendCommand(cmd)
}

// Switch work mode for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) ModeSwitch(channel int) {
	cmd := n.composeCommand(modeSwitch, channel)
	n.sendCommand(cmd)
}

// Set change color speed for specified channel
//
// Avialable for SD111-180 only
func (n *NooliteAdapter) EffectSpeed(channel int) {
	cmd := n.composeCommand(effectSpeed, channel)
	n.sendCommand(cmd)
}
