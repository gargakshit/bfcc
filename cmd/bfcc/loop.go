package main

import (
	"github.com/llir/llvm/ir"
)

type loop struct {
	start *ir.Block
	end   *ir.Block
}

type loopStack struct {
	underlying []*loop
}

func newLoopStack() *loopStack {
	return &loopStack{}
}

func (l *loopStack) push(loop *loop) {
	l.underlying = append(l.underlying, loop)
}

func (l *loopStack) pop() (*loop, bool) {
	underlyingLen := len(l.underlying)
	if underlyingLen == 0 {
		return nil, false
	}

	ret := l.underlying[underlyingLen-1]
	l.underlying = l.underlying[:underlyingLen-1]

	return ret, true
}
