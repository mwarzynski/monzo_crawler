package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Fatalf("HTTP Server is not implemented.")
	// Should use the internal/transport/http package.
}
