package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

const memSize = 1<<16 - 1

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("No input file provided")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(file)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewReader(os.Stdin)

	var mem [memSize]int8
	ptr := 0

	idx := 0

	for idx != len(content) {
		switch content[idx] {
		case '+':
			mem[ptr]++

		case '-':
			mem[ptr]--

		case '<':
			ptr--

		case '>':
			ptr++

		case '[':
			if mem[ptr] != 0 {
				break
			}

			open := 1
			for open > 0 {
				idx++

				switch content[idx] {
				case '[':
					open++
				case ']':
					open--
				}
			}

		case ']':
			open := 1
			for open > 0 {
				idx--

				switch content[idx] {
				case '[':
					open--
				case ']':
					open++
				}
			}

			idx--

		case '.':
			fmt.Print(string(mem[ptr]))

		case ',':
			b, err := sc.ReadByte()
			if err != nil {
				log.Fatal(b)
			}

			mem[ptr] = int8(b)
		}

		idx++
	}
}
