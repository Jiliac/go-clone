package main

import (
	"fmt"
	"log"

	"os"
	"sync"
	"syscall"
)

//#include <unistd.h>
//
//int cFork() {
//    pid_t pid;
//    pid = fork();
//    return ((int)pid);
//}
import "C"

const (
	msg     = "Hello dad, it's your son :-)."
	bufSize = 2000
)

func main() {
	forkComTest()
}

func forkComTest() {
	r, w, err := pipe()
	if err != nil {
		log.Printf("Failed to create pipes: %v.\n", err)
		return
	}

	pid := int(C.cFork())
	if pid == 0 { // Child
		fmt.Printf("(from child)\tHello.\n")
		err = r.Close()
		if err != nil {
			log.Printf("Child could not close reading pipe: %v.\n", err)
			return
		}

		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Printf("Could not write in pipe: %v.\n", err)
		}

	} else { // Parent
		fmt.Printf("(from parent)\tChild pid = %d.\n", pid)
		err = w.Close()
		if err != nil {
			log.Printf("Parent could not close writing pipe: %v\n.", err)
			return
		}

		buf := make([]byte, bufSize)

		_, err := r.Read(buf)
		if err != nil {
			log.Printf("Could not reat from pipe: %v.\n", err)
		}

		msg := string(buf)
		fmt.Printf("(from parent)\tMessage: '%s'\n", msg)
	}
}

func forkTest() {
	pid := int(C.cFork())
	if pid == 0 {
		fmt.Println("Child")
	} else {
		fmt.Println("Parent")
	}
}

func pipeTest() {
	var wg sync.WaitGroup
	wg.Add(2)

	r, w, err := pipe()
	if err != nil {
		log.Printf("Failed to create pipes: %v.\n", err)
		return
	}

	go func() {
		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Printf("Could not write in pipe: %v.\n", err)
		}

		wg.Done()
	}()

	go func() {
		buf := make([]byte, bufSize)
		_, err := r.Read(buf)
		if err != nil {
			log.Printf("Could not reat from pipe: %v.\n", err)
		}

		msg := string(buf)
		fmt.Printf("buf: %+v\n", msg)

		wg.Done()
	}()

	wg.Wait()
}

// I guess go dev really didn't want us to do this.
// In the os package, Pipe() mark the pipes
func pipe() (r, w *os.File, err error) {
	var p [2]int
	err = syscall.Pipe(p[0:])
	if err != nil {
		return r, w, err
	}

	r = os.NewFile(uintptr(p[0]), "|0")
	w = os.NewFile(uintptr(p[1]), "|1")
	return r, w, err
}
