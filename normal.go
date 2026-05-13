package main

import (
	"fmt"
	"os"
)

func PerformNormal(r *Reader, line byte) {
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
	register = RegisterTable[reg][w]

	switch mode {
	default:
		fmt.Println("mode cannot exceed 3")
		os.Exit(1)
		return
	case 0:
		if reg_mem == 0b110 {
			direct_address := ParseImmediateData(r, 2)
			register_memory = fmt.Sprintf("[%d]", direct_address)
			return
		}
		register_memory = fmt.Sprintf("[%s]", RegisterMemoryTable[reg_mem])
		return
	case 1, 2:
		num := ParseImmediateData(r, mode)
		if num == 0 {
			register_memory = fmt.Sprintf("[%s]", RegisterMemoryTable[reg_mem])
			return
		}
		if num < 0 {
			register_memory = fmt.Sprintf("[%s - %d]", RegisterMemoryTable[reg_mem], num*-1)
			return
		}
		register_memory = fmt.Sprintf("[%s + %d]", RegisterMemoryTable[reg_mem], num)
		return
	case 3:
		register_memory = RegisterTable[reg_mem][w]
		return
	}
}
