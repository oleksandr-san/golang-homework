package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func receiveRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := client.Receive(); err == nil {
		log.Println("...Connection was closed by peer")
		if err := client.Close(); err != nil {
			log.Println("...Receive side closing error", err)
		}
	}
}

func sendRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := client.Send(); err == nil {
		log.Println("...EOF")
		if err := client.Close(); err != nil {
			log.Println("...Send side closing error", err)
		}
	}
}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("Usage: program [-timeout t] host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Println("...Connected to", address)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go receiveRoutine(client, wg)

	wg.Add(1)
	go sendRoutine(client, wg)

	wg.Wait()
	client.Close()
}
