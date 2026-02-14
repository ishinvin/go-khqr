package main

import (
	"fmt"
	"log"

	khqr "github.com/ishinvin/go-khqr"
)

func main() {
	// First, generate a QR to decode
	data, err := khqr.GenerateIndividual(khqr.IndividualInfo{
		BakongAccountID: "ishin_vin@bkrt",
		MerchantName:    "Ishin Vin",
		MerchantCity:    "Phnom Penh",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Decode the QR string
	decoded, err := khqr.Decode(data.QR)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Decoded KHQR ===")
	fmt.Println("Bakong Account ID:", decoded.BakongAccountID)
	fmt.Println("Merchant Name:    ", decoded.MerchantName)
	fmt.Println("Merchant City:    ", decoded.MerchantCity)
	fmt.Println("Merchant Type:    ", decoded.MerchantType)
	fmt.Println("Currency:         ", decoded.TransactionCurrency)
	fmt.Println("Amount:           ", decoded.TransactionAmount)
	fmt.Println("CRC:              ", decoded.CRC)
}
