package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	"github.com/floreks/breathalyzer/service"
	"github.com/spf13/pflag"
)

var (
	argPort = pflag.Int("port", 3000, "The port to listen on for incoming HTTP requests")
)

func main() {
	// Set logging out to standard console out
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// Register handler
	mqService, err := service.NewMQ3Service()
	if err != nil {
		panic(err)
	}

	restful.Add(mqService.Handler())

	log.Printf("Listening on port: %d", *argPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *argPort), nil))
}
