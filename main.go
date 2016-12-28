package main

import (
	"cwmetric2"
	"fmt"
	"log"
)

func init() {
	log.SetFlags(0)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	cwm2, err := cwmetric2.ParseFlag()

	if err != nil {
		log.Fatal(err)
	}

	value, err := cwm2.GetMetricStatistics()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)
}
