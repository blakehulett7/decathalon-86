package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func PerformArithmeticAccumulator(r *Reader, line byte) {
	operation := ParseOperation(line)

	w := line & 0b1
	memory_address := ParseImmediateData(r, w+1)
	fmt.Printf("%s ax, [%d]\n", operation, memory_address)
}

func PerformArithmeticNormal(r *Reader, line byte) {
	operation := ParseOperation(line)

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

	fmt.Printf("%s %s, %s\n", operation, dest, src)
}

func PerformArithmeticToRegisterMemory(r *Reader, line byte) {
	operation := ParseOperation(line)

	s := line >> 1 & 0b1
	w := line & 0b1

	line = r.Read()
	mode, reg_mem := ParseImmediateToRegisterMemoryArguments(line)

	_, register_memory := ParseRegisters(r, w, mode, 0, reg_mem)
	immediate := ParseArithmeticData(r, w, s)

	if s == w {
		size := "byte"
		if w == 1 {
			size = "word"
		}
		fmt.Printf("%s %s %s, %d\n", operation, size, register_memory, immediate)
		return
	}

	fmt.Printf("%s %s, %d\n", operation, register_memory, immediate)
}

func ParseOperation(line byte) string {
	code := line & 0b00111000
	switch code {
	default:
		fmt.Printf("invalid arithmetic operation: %b\n", line)
		os.Exit(1)
		return ""
	case 0b000000:
		return "add"
	case 0b101000:
		return "sub"
	case 0b111000:
		return "cmp"
	}
}

func ParseArithmeticData(r *Reader, w, s uint8) int16 {
	if s == 1 || w == 0 {
		line := r.Read()
		if s == 1 && isNegative(line) {
			return int16(int8(line))
		}
		return int16(line)
	}

	line := r.Read()
	wide_line := r.Read()
	data := []byte{line, wide_line}

	return int16(binary.LittleEndian.Uint16(data))
}
