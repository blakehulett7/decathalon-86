package main

import "fmt"

func PerformAccumulatorToMemory(r *Reader, line byte) {
	w := line & 0b1
	memory_address := ParseImmediateData(r, w+1)
	fmt.Printf("mov [%d], ax\n", memory_address)
}

func PerformMemoryToAccumulator(r *Reader, line byte) {
	w := line & 0b1
	memory_address := ParseImmediateData(r, w+1)
	fmt.Printf("mov ax, [%d]\n", memory_address)
}
