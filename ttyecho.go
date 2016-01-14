package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"unsafe"
)

var newline = flag.Bool("n", false, "append newline")

func insert(fd uintptr, s string) error {
	for _, c := range s {
		ptr := uintptr(unsafe.Pointer(&c))
		if _, _, e := unix.Syscall(unix.SYS_IOCTL, fd, unix.TIOCSTI, ptr); e != 0 {
			return fmt.Errorf("syscall errno: %d\n", e)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ttyecho: [-n] DEVNAME COMMAND\n")
		flag.PrintDefaults()
	}
	if flag.Arg(0) == "" || flag.Arg(1) == "" {
		flag.Usage()
		os.Exit(0)
	}
	f, err := os.OpenFile(flag.Arg(0), os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	if err := insert(f.Fd(), flag.Arg(1)); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	if *newline {
		if err := insert(f.Fd(), "\n"); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
