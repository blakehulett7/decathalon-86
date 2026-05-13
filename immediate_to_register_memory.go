package main

import (
	"fmt"
	"os"
)

func PerformImmediateToRegisterMemory(r *Reader, line byte) {
	w := line & 0b1
	line = r.Read()
	mode, reg_mem := ParseImmediateToRegisterMemoryArguments(line)

	_, register_memory := ParseRegisters(r, w, mode, 0, reg_mem)

	immediate := ParseImmediateData(r, w+1)
	size := "byte"
	if w > 0 {
		size = "word"
	}

	fmt.Printf("mov %s, %s %d\n", register_memory, size, immediate)
}

func ParseImmediateToRegisterMemoryArguments(line byte) (mode, reg_mem uint8) {
	mode = line >> 6

	reg_mem = line & 0b00000111
	if reg_mem > 7 {
		fmt.Println("invalid second register code, only 8 register codes are supported")
		os.Exit(1)
	}

	return
}
