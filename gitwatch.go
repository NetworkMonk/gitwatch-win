package main

import (
	"math/rand"
	"time"

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

}
