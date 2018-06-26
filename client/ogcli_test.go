package client

import (
	"log"
	"testing"
)

// TestRunClient ...
func TestRunClient(t *testing.T) {
	og, err := RunClient("localhost:9091", "binary", false, false, false, "http://ogp.me/")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(og)
}
