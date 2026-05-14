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
		PerformOp(r, opcode, line)
	}
}

func ParseOpCode(line byte) OpCode {
	code := line >> 4
	if code == 0b1011 {
		return ImmediateToRegister
	}

	code = line >> 2
	if code == 0b100010 {
		return Normal
	}

	code = line >> 1
	if code == 0b1100011 {
		return ImmediateToRegisterMemory
	}
	if code == 0b1010000 {
		return MemoryToAccumulator
	}
	if code == 0b1010001 {
		return AccumulatorToMemory
	}

	code = line & 0b11000100
	if code == 0b00000000 {
		return ArithmeticNormal
	}
	if code == 0b10000000 {
		return ArithmeticToRegisterMemory
	}
	if code == 0b00000100 {
		return ArithmeticAccumulator
	}

	if line == JNZ {
		return Jnz
	}

	fmt.Printf("invalid op code, line: %08b\n", line)
	os.Exit(1)
	return NoOp
}

func PerformOp(r *Reader, code OpCode, line byte) {
	switch code {
	default:
		fmt.Println("invalid op code")
		os.Exit(1)
	case ArithmeticAccumulator:
		PerformArithmeticAccumulator(r, line)
	case ArithmeticNormal:
		PerformArithmeticNormal(r, line)
	case ArithmeticToRegisterMemory:
		PerformArithmeticToRegisterMemory(r, line)
	case AccumulatorToMemory:
		PerformAccumulatorToMemory(r, line)
	case ImmediateToRegister:
		PerformImmediateToRegister(r, line)
	case ImmediateToRegisterMemory:
		PerformImmediateToRegisterMemory(r, line)
	case MemoryToAccumulator:
		PerformMemoryToAccumulator(r, line)
	case Normal:
		PerformNormal(r, line)
	case Jnz:
		fmt.Print("jnz ")
		line = r.Read()
		fmt.Println(int8(line))
	}
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
	ArithmeticAccumulator OpCode = iota
	ArithmeticNormal
	ArithmeticToRegisterMemory
	AccumulatorToMemory
	ImmediateToRegister
	ImmediateToRegisterMemory
	MemoryToAccumulator
	Normal
	NoOp
	Jnz
)

const (
	JNZ byte = 0b01110101
)
