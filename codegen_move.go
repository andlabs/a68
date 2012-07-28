// 15 july 2012
package main

import (
	"fmt"
)

// move <ea>,<ea>
func _move_ea_ea(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 1 },		// ...wait, what?!
		'w':	{ 1, 1 },
		'l':	{ 1, 0 },
	}

	if src.Type == 'a' && suffix == 'b' {		// no byte reads from address registers
		// TODO print the register?
		return fmt.Errorf("move.b cannot be used with an address register source")
	}
	WriteBits(0, 0)
	WriteBits(sizes[suffix]...)
	fDest := WriteEA(dest)
	WriteEANow(src)			// source, then destination
	if fDest != nil {
		fDest()
	}
	return nil
}

// move <ea>,ccr
func _move_ea_ccr(suffix rune, src Operand) error {
	if suffix != ' ' {
		return fmt.Errorf("move to ccr cannot have suffix")
	}
	if !addressingModeValid(src.Type, "d*+-$%^&wl#") {
		return fmt.Errorf("cannot move aN/ccr/sr/usp to ccr")
	}
	WriteBits(0, 1, 0, 0)
	WriteBits(0, 1, 0, 0)
	WriteBits(1, 1)
	WriteEANow(src)
	return nil
}

// move sr,<ea>
func _move_sr_ea(suffix rune, dest Operand) error {
	if suffix != ' ' {
		return fmt.Errorf("move from sr cannot have suffix")
	}
	if !addressingModeValid(src.Type, "d*+-$%wl") {
		return fmt.Errorf("cannot move sr to aN/(pc)/immediate/ccr/sr/usp")
	}
	WriteBits(0, 1, 0, 0)
	WriteBits(0, 0, 0, 0)
	WriteBits(1, 1)
	WriteEANow(dest)
	return nil
}

// move <ea>,<ea>
// move <ea>,ccr
// move sr,<ea>
func o_move(suffix rune, src Operand, dest Operand) error {
	if dest.Type == 'c' {
		return _move_ea_ccr(suffix, src)
	}
	if src.Type == 's' {
		return _move_sr_ea(suffix, dest)
	}
	return _move_ea_ea(suffix, src, dest)
}

// TODO movea (in another file)
