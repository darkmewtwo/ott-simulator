package main

import (
	"log"
	"simulator/internal/orchastrator"
	"time"
)

func main() {
	o, err := orchastrator.NewOrchastrator(12345, 539)

	if err != nil {
		log.Fatal(err)
	}
	for {
		if err = o.Run(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(10 * time.Second)
	}
}
