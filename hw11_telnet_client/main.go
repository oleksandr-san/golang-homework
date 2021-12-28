package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

func receiveRoutine(ctx context.Context, client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if err := client.Receive(); err != nil {
				log.Print(err)
				// cancel()
				break OUTER
			}
		}
	}

	log.Printf("Finished receiveRoutine")
}

func sendRoutine(ctx context.Context, client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	// scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if err := client.Send(); err != nil {
				log.Println(err)
				break OUTER
			}
		}
	}
	log.Printf("Finished sendRoutine")
}

func main() {
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("Usage: program [-timeout t] host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	log.Println("Connecting to", address)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go receiveRoutine(ctx, client, wg)

	wg.Add(1)
	go sendRoutine(ctx, client, wg)

	wg.Wait()
	cancel()
	client.Close()
}
