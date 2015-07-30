package main

import (
	"flag"
	"fmt"
	"listen"
	"log"
	"net/http"
	"os"
	"os/signal"
	"speak"
	"syscall"
)

func main() {
	mode := flag.String("mode", "both", "mode|speak")
	secrets := flag.String("secrets", "secrets.yml", "path to secrets")
	flag.Parse()

	read_signal := make(chan os.Signal, 2)
	signal.Notify(read_signal, os.Interrupt, syscall.SIGTERM)
	quit := make(chan bool)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if *mode == "listen" || *mode == "both" {
		fmt.Print("listen")
		go listen.Listen(secrets, quit)
	}

	if *mode == "speak" || *mode == "both" {
		fmt.Print("speak")
		go speak.Speak(secrets, quit)
	}

	<-read_signal
	fmt.Print("received a signal, stopping on quit\n")
	quit <- true
}
