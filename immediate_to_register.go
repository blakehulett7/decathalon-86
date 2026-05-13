package main

import (
	"encoding/binary"
	"fmt"
)

func PerformImmediateToRegister(r *Reader, line byte) {
	displacement, register := ParseImmediateToRegisterParams(line)
	num := ParseImmediateData(r, displacement)
	fmt.Printf("mov %s, %d\n", register, num)
}

func ParseImmediateData(r *Reader, displacement uint8) int16 {
	if displacement == 1 {
		line := r.Read()
		if isNegative(line) {
			return int16(int8(line))
		}
		return int16(line)
	}

	line := r.Read()
	wide_line := r.Read()
	data := []byte{line, wide_line}

	return int16(binary.LittleEndian.Uint16(data))
}

func ParseImmediateToRegisterParams(line byte) (displacement uint8, register string) {
	data := line & 0b00001111
	w := data >> 3
	register_code := data & 0b00000111

	displacement = w + 1
	register = RegisterTable[register_code][w]
	return
}

func isNegative(line byte) bool {
	shifted := line >> 7
	leading_bit := shifted & 0b1
	return leading_bit == 1
}
