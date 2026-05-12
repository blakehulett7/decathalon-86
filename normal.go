package main

import (
	"fmt"
	"os"
)

func ParseNormalArguments(line byte, w uint8) (mode, reg, reg_mem uint8) {
	mode = line >> 6

	reg = (line >> 3) & 0b00000111
	if reg > 7 {
		fmt.Println("invalid first register code, only 8 register codes are supported")
		os.Exit(1)
	}

	reg_mem = line & 0b00000111
	if reg_mem > 7 {
		fmt.Println("invalid second register code, only 8 register codes are supported")
		os.Exit(1)
	}

	return
}

func ParseNormalModifiers(line byte) (direction, w uint8) {
	direction = (line >> 1) & 0b00000001
	w = line & 0b00000001
	return
}

func ParseRegisters(r *Reader, w, mode, reg, reg_mem uint8) (register, register_memory string) {
	switch mode {
	default:
		fmt.Println("mode cannot exceed 3")
		os.Exit(1)
	case 3:
		register = RegisterTable[reg][w]
		register_memory = RegisterTable[reg_mem][w]
		return
	}

	return "", ""
}
