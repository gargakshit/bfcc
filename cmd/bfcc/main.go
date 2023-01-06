package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
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

	module := ir.NewModule()

	getChar := module.NewFunc("getchar", types.I8)

	putChar := module.NewFunc("putchar", types.I32,
		ir.NewParam("p1", types.I8))

	memset := module.NewFunc("memset", types.Void,
		ir.NewParam("ptr", types.I8Ptr),
		ir.NewParam("val", types.I8),
		ir.NewParam("len", types.I64))

	entrypoint := module.NewFunc("main", types.I32)
	builder := entrypoint.NewBlock("")

	arrayType := types.NewArray(memSize, types.I8)

	mem := builder.NewAlloca(arrayType)
	ptr := builder.NewAlloca(types.I32)

	builder.NewStore(constant.NewInt(types.I32, 0), ptr)
	builder.NewCall(memset,
		builder.NewGetElementPtr(
			arrayType,
			mem,
			constant.NewInt(types.I64, 0),
			constant.NewInt(types.I64, 0),
		),
		constant.NewInt(types.I8, 0),
		constant.NewInt(types.I64, memSize),
	)

	s := newLoopStack()

	for _, r := range content {
		switch r {
		case '+':
			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			t1 := builder.NewAdd(builder.NewLoad(types.I8, d), constant.NewInt(types.I8, 1))
			builder.NewStore(t1, d)

		case '-':
			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			t1 := builder.NewAdd(builder.NewLoad(types.I8, d), constant.NewInt(types.I8, -1))
			builder.NewStore(t1, d)

		case '<':
			t1 := builder.NewAdd(builder.NewLoad(types.I32, ptr), constant.NewInt(types.I8, -1))
			builder.NewStore(t1, ptr)

		case '>':
			t1 := builder.NewAdd(builder.NewLoad(types.I32, ptr), constant.NewInt(types.I8, 1))
			builder.NewStore(t1, ptr)

		case '.':
			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			builder.NewCall(putChar, builder.NewLoad(types.I8, d))

		case ',':
			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			ch := builder.NewCall(getChar)
			builder.NewStore(ch, d)

		case '[':
			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			cmp := builder.NewICmp(
				enum.IPredNE,
				builder.NewLoad(types.I8, d),
				constant.NewInt(types.I8, 0),
			)

			l := &loop{
				start: entrypoint.NewBlock(""),
				end:   entrypoint.NewBlock(""),
			}
			s.push(l)

			builder.NewCondBr(cmp, l.start, l.end)
			builder = l.start

		case ']':
			l, ok := s.pop()
			if !ok {
				log.Fatal("unbalanced brackets")
			}

			d := builder.NewGetElementPtr(
				arrayType,
				mem,
				constant.NewInt(types.I64, 0),
				builder.NewLoad(types.I32, ptr),
			)

			cmp := builder.NewICmp(
				enum.IPredNE,
				builder.NewLoad(types.I8, d),
				constant.NewInt(types.I8, 0),
			)

			builder.NewCondBr(cmp, l.start, l.end)
			builder = l.end
		}
	}

	builder.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(module.String())
}
