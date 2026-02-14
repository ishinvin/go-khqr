package main

import (
	"fmt"
	"log"
	"time"

	khqr "github.com/ishinvin/go-khqr"
)

func main() {
	// Static QR (no amount, reusable)
	staticData, err := khqr.GenerateMerchant(khqr.MerchantInfo{
		BakongAccountID: "ishin_vin@bkrt",
		MerchantName:    "Ishin Coffee",
		MerchantCity:    "Phnom Penh",
		MerchantID:      "123456",
		AcquiringBank:   "Bakong",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Static Merchant QR ===")
	fmt.Println("QR:", staticData.QR)
	fmt.Println("MD5:", staticData.MD5())

	// Dynamic QR (with amount and additional data)
	dynamicData, err := khqr.GenerateMerchant(khqr.MerchantInfo{
		BakongAccountID:     "ishin_vin@bkrt",
		MerchantName:        "Ishin Coffee",
		MerchantCity:        "Phnom Penh",
		MerchantID:          "123456",
		AcquiringBank:       "Bakong",
		Currency:            khqr.USD,
		Amount:              10.50,
		ExpirationTimestamp: time.Now().Add(5 * time.Minute).UnixMilli(),
		BillNumber:          "INV-2026-001",
		StoreLabel:          "Main Branch",
		TerminalLabel:       "Cashier_1",
		Purpose:             "Coffee",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n=== Dynamic Merchant QR (USD) ===")
	fmt.Println("QR:", dynamicData.QR)
	fmt.Println("MD5:", dynamicData.MD5())
}
