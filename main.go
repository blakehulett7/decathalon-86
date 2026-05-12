package main

import (
	"flag"
	"fmt"
	"os"
)

var RegisterMemoryTable = [8]string{
	"bx + si",
	"bx + di",
	"bp + si",
	"bp + di",
	"si",
	"di",
	"bp",
	"bx",
}

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

	fmt.Println("bits 16")

	r := &Reader{
		Data:   assembled,
		Cursor: 0,
	}

	for r.Cursor < uint8(len(r.Data)) {
		line := r.Read()
		opcode := ParseOpCode(line)

		switch opcode {
		default:
			fmt.Println("invalid op code")
			os.Exit(1)
		case ImmediateToRegister:
			displacement, register := ParseImmediateToRegisterParams(line)
			num := ParseImmediateData(r, displacement)
			fmt.Printf("mov %s, %d\n", register, num)
		case Normal:
			direction, w := ParseNormalModifiers(line)
			line = r.Read()
			mode, reg, reg_mem := ParseNormalArguments(line, w)

			register, register_memory := ParseRegisters(r, w, mode, reg, reg_mem)

			invert_registers := direction == 0
			dest := register
			src := register_memory

			if invert_registers {
				src = register
				dest = register_memory
			}

			fmt.Printf("mov %s, %s\n", dest, src)
		}
	}
}

func ParseOpCode(line byte) OpCode {
	code := line >> 4
	if code == 0b00001011 {
		return ImmediateToRegister
	}

	code = line >> 2
	if code == 0b00100010 {
		return Normal
	}

	fmt.Printf("invalid op code, line: %08b\n", line)
	os.Exit(1)
	return NoOp
}

func PrintByte(b byte) {
	fmt.Printf("%08b\n", b)
}

type Reader struct {
	Data   []byte
	Cursor uint8
}

func (r *Reader) Read() byte {
	if r.Cursor > uint8(len(r.Data)) {
		fmt.Println("cursor can't exceed data lenght")
		os.Exit(1)
	}

	b := r.Data[r.Cursor]
	r.Cursor++
	return b
}

type OpCode uint8

const (
	ImmediateToRegister OpCode = iota
	Normal
	NoOp
)
