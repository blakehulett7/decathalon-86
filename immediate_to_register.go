package main

import "encoding/binary"

func ParseImmediateData(r *Reader, displacement uint8) uint16 {
	if displacement == 1 {
		line := r.Read()
		return uint16(line)
	}

	line := r.Read()
	wide_line := r.Read()
	data := []byte{line, wide_line}

	return binary.LittleEndian.Uint16(data)
}

func ParseImmediateToRegisterParams(line byte) (displacement uint8, register string) {
	data := line & 0b00001111
	w := data >> 3
	register_code := data & 0b00000111

	displacement = w + 1
	register = RegisterTable[register_code][w]
	return
}
