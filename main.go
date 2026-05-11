package main

import (
	"flag"
	"fmt"
	"os"
)

const MovOperand = byte(34)
const RegisterToRegisterMod = byte(3)

var RegisterTable = [8][2]string{
	{"al", "ax"},
	{"cl", "cx"},
	{"dl", "dx"},
	{"bl", "bx"},
	{"ah", "sp"},
	{"ch", "bp"},
	{"dh", "si"},
	{"bh", "di"},
}

func main() {
	flag.Parse()
	path := flag.Arg(0)

	assembled, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(assembled)%2 != 0 {
		fmt.Println("invalid instruction length, must be a multiple of 2 bytes")
		os.Exit(1)
	}

	for i := 1; i < len(assembled); i += 2 {
		line := assembled[i-1 : i+1]
		ParseLine(line)
	}

}

func ParseLine(line []byte) {
	if len(line) != 2 {
		fmt.Println("unexpected instruction line length")
		os.Exit(1)
	}

	first_byte := line[0]
	second_byte := line[1]

	operand := first_byte >> 2
	if operand != MovOperand {
		fmt.Println("unexpected operand, only MOV is supported")
		os.Exit(1)
	}

	w := first_byte & 0b00000001

	d := (first_byte >> 1) & 0b00000001
	invert_registers := d == 1

	mod := second_byte >> 6
	if mod != RegisterToRegisterMod {
		fmt.Println("unexpected mode, only register to register is supported")
		os.Exit(1)
	}

	first_register_row := (second_byte >> 3) & 0b00000111

	if first_register_row > 7 {
		fmt.Println("invalid first register code, only 8 register codes are supported")
		os.Exit(1)
	}

	first_register := RegisterTable[first_register_row][w]

	second_register_row := second_byte & 0b00000111

	second_register := RegisterTable[second_register_row][w]

	if second_register_row > 7 {
		fmt.Println("invalid second register code, only 8 register codes are supported")
		os.Exit(1)
	}

	dest := first_register
	src := second_register

	if !invert_registers {
		src = first_register
		dest = second_register
	}

	fmt.Println("bits 16")
	fmt.Print("mov ")
	fmt.Print(dest)
	fmt.Print(", ")
	fmt.Print(src)
	fmt.Println()
}

func PrintByte(b byte) {
	fmt.Printf("%08b\n", b)
}
