# bfcc

[Brainf*ck](https://en.wikipedia.org/wiki/Brainfuck) compiler and interpreter
written in Go.

## Interpreter

The interpreter contained in the `cmd/bfi` directory is a basic brainf*ck
interpreter that does no optimizations.

## Compiler

The compiler contained in `cmd/bfcc` is a LLVM-based compiler for brainf*ck.
Currently, it emits [LLVM IR](https://llvm.org/docs/LangRef.html) on the
`stdout` which may be piped to clang for compilation.

The compiled program is expected to be linked with `libc` for `getch`, `putch`
and `memset`.

### Compiling a program

```shell
go run ./cmd/bfcc <path-to-brainf*ck-program> | clang -x ir -o <path-to-output> -
```

## Performance

All these tests are run on a R7-5700U running Windows 11 with the following
flags for `clang`: `-flto -O3 -fuse-ld=lld-link`

| Program            | `bfi`                   | `bfcc`  |
|--------------------|-------------------------|---------|
| Tower of hanoi     | 14.895s                 | 0.029s  |
| Mandelbrot         | 33.511s                 | 0.613s  |
| Mandelbrot Extreme | did not complete in 30m | 76.102s |

The compiled version here is anywhere from **50-500 times** faster.
