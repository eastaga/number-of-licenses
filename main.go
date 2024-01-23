package main

import (
	"log"

	"number-of-licenses/helpers"
)

func main() {

	// read os env
	err := helpers.Validate()
	if err != nil {
		log.Fatal("Error encountered ", err)
	}

	log.Printf("Total number of licenses required: %d", helpers.GetTotalLicenses())
}
