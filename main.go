package main

// Tests whether io blocking will stall the process
import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

func main() {
	gomaxprocs := runtime.GOMAXPROCS(0)
	blockingCount := gomaxprocs * 2
	fmt.Printf("GOMAXPROCS=%d running %d blocking goroutines\n", gomaxprocs, blockingCount)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go handleSignals(cancel, os.Interrupt, os.Kill)
	for i := 0; i < blockingCount; i++ {
		go work(i)
	}
	<-ctx.Done()
}

func work(i int) {
	fmt.Printf("Started %d\n", i)
	name := fmt.Sprintf("./pipe%d", i)
	os.Remove(name)
	err := syscall.Mkfifo(name, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Calling mkfifo: %s\n", err)
	}
	file, err := os.OpenFile(name, os.O_WRONLY, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
	}
	defer file.Close()

	n, err := fmt.Fprintf(file, "Hello world!\n")
	if err != nil {
		fmt.Printf("Wrote %d, err: %s\n", n, err)
	} else {
		fmt.Printf("Wrote %d\n", n)
	}
}

func handleSignals(cancel context.CancelFunc, signals ...os.Signal) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	s := <-signalChan
	fmt.Printf("Got %s, shutting down...", s)
	cancel()

	time.Sleep(5 * time.Second)
	os.Exit(0)
}
