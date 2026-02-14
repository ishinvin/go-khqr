package main

import (
	"fmt"
	"log"

	khqr "github.com/ishinvin/go-khqr"
)

func main() {
	// First, generate a QR to verify
	data, err := khqr.GenerateIndividual(khqr.IndividualInfo{
		BakongAccountID: "ishin_vin@bkrt",
		MerchantName:    "Ishin Vin",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Verify valid QR
	fmt.Println("=== Verify Valid QR ===")
	if err := khqr.Verify(data.QR); err != nil {
		fmt.Println("Invalid:", err)
	} else {
		fmt.Println("QR is valid")
	}

	// Verify invalid QR
	fmt.Println("\n=== Verify Invalid QR ===")
	if err := khqr.Verify("invalid-qr-string"); err != nil {
		fmt.Println("Invalid:", err)
	} else {
		fmt.Println("QR is valid")
	}
}
