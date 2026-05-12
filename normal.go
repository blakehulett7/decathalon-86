package main

import (
	"fmt"
	"os"
)

func ParseNormalArguments(line byte, w uint8) (mode uint8, first_register, second_register string) {
	mode = line >> 6

	first_register_row := (line >> 3) & 0b00000111

	if first_register_row > 7 {
		fmt.Println("invalid first register code, only 8 register codes are supported")
		os.Exit(1)
	}

	first_register = RegisterTable[first_register_row][w]

	second_register_row := line & 0b00000111

	second_register = RegisterTable[second_register_row][w]

	if second_register_row > 7 {
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
