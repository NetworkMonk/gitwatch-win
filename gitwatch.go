package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/NetworkMonk/gitwatch/config"
	"github.com/NetworkMonk/gitwatch/watch"
	"github.com/NetworkMonk/service"
)

const svcName = "gitwatch"
const svcTitle = "gitwatch"
const version = "1.0.0"

func main() {
	rand.Seed(time.Now().UnixNano())
	service.Handle(svcName, svcTitle, execute)
}

func execute() {
	// Get the current exe path
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Load configuration
	configPath := strings.Replace(exePath, "gitwatch.exe", "", 1) + "gitwatch.json"
	configuration, configErr := config.Load(configPath)
	if configErr != nil {
		log.Fatal(configErr)
	}

	// Start the watch process
	watch.Start(configuration)
}
