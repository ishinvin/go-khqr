package main

import (
	"fmt"
	"log"
	"time"

	khqr "github.com/ishinvin/go-khqr"
)

func main() {
	// Static QR (no amount, reusable)
	staticData, err := khqr.GenerateIndividual(khqr.IndividualInfo{
		BakongAccountID: "ishin_vin@bkrt",
		MerchantName:    "Ishin Vin",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Static Individual QR ===")
	fmt.Println("QR:", staticData.QR)
	fmt.Println("MD5:", staticData.MD5())

	// Dynamic QR (with amount and expiration)
	dynamicData, err := khqr.GenerateIndividual(khqr.IndividualInfo{
		BakongAccountID:     "ishin_vin@bkrt",
		MerchantName:        "Ishin Vin",
		Currency:            khqr.USD,
		Amount:              5.00,
		ExpirationTimestamp: time.Now().Add(5 * time.Minute).UnixMilli(),
		BillNumber:          "INV-2026-001",
		MobileNumber:        "85512345678",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n=== Dynamic Individual QR (USD) ===")
	fmt.Println("QR:", dynamicData.QR)
	fmt.Println("MD5:", dynamicData.MD5())
}
