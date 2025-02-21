package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cirdartcounter "github.com/One-Hundred-Eighty/Circle/backend/cir-dartcounter"
	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
)

func main() {
	// create loggers
	mainLogger := dartmasterlogger.NewDartmasterLogger("[main] ")
	dartcounterServerLogger := dartmasterlogger.NewDartmasterLogger("[dartcounter-server] ")

	mainLogger.Println("boot servers...")
	fmt.Println()

	// create servers
	dartcounterServer := cirdartcounter.NewServer(dartcounterServerLogger, "8888")

	// run boot the servers
	go func() {
		err := dartcounterServer.ListenAndServe()
		if err != nil {
			mainLogger.PrintfErr("error occurred starting server: %v", err)
		}
	}()

	// create a channel to listen for a signal that shuts down the running program
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for a termination signal --> allows the function go clean up itself with the pre-configured defer functionalities.
	sig := <-sigChan
	fmt.Println()
	mainLogger.Printf("received signal: %v. clean up...", sig)
}
