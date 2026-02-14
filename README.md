# KHQR SDK for Go

Go SDK for generating, decoding, and verifying Bakong KHQR.

## Installation

```bash
go get github.com/ishinvin/go-khqr
```

## Usage

```go
import khqr "github.com/ishinvin/go-khqr"
```

### Generate Individual QR

**Static QR (no amount):**

```go
data, err := khqr.GenerateIndividual(khqr.IndividualInfo{
    BakongAccountID: "ishin_vin@bkrt",
    MerchantName:    "Ishin Vin",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.QR)
```

**Dynamic QR (with amount and expiration):**

```go
data, err := khqr.GenerateIndividual(khqr.IndividualInfo{
    BakongAccountID:     "ishin_vin@bkrt",
    MerchantName:        "Ishin Vin",
    Currency:            khqr.USD,
    Amount:              5.00,
    ExpirationTimestamp:  time.Now().Add(5 * time.Minute).UnixMilli(),
    BillNumber:          "INV-2026-001",
    MobileNumber:        "85512345678",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.QR)
fmt.Println(data.MD5())
```

**With alternate language:**

```go
data, err := khqr.GenerateIndividual(khqr.IndividualInfo{
    BakongAccountID:     "ishin_vin@bkrt",
    MerchantName:        "Ishin Vin",
    Currency:            khqr.KHR,
    Amount:              50000,
    ExpirationTimestamp:  time.Now().Add(10 * time.Minute).UnixMilli(),
    AltLanguagePreference: "km",
    AltMerchantName:     "វិន អ៊ីស៊ីន",
    AltMerchantCity:     "ភ្នំពេញ",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.QR)
```

### Generate Merchant QR

**Static QR:**

```go
data, err := khqr.GenerateMerchant(khqr.MerchantInfo{
    BakongAccountID: "ishin_vin@bkrt",
    MerchantName:    "Ishin Coffee",
    MerchantCity:    "Phnom Penh",
    MerchantID:      "123456",
    AcquiringBank:   "Bakong",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.QR)
```

**Dynamic QR (with amount and additional data):**

```go
data, err := khqr.GenerateMerchant(khqr.MerchantInfo{
    BakongAccountID:     "ishin_vin@bkrt",
    MerchantName:        "Ishin Coffee",
    MerchantCity:        "Phnom Penh",
    MerchantID:          "123456",
    AcquiringBank:       "Bakong",
    Currency:            khqr.USD,
    Amount:              10.50,
    ExpirationTimestamp:  time.Now().Add(5 * time.Minute).UnixMilli(),
    BillNumber:          "INV-2026-001",
    StoreLabel:          "Main Branch",
    TerminalLabel:       "Cashier_1",
    Purpose:             "Coffee",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.QR)
fmt.Println(data.MD5())
```

### Decode QR

```go
decoded, err := khqr.Decode(qrString)
if err != nil {
    log.Fatal(err)
}
fmt.Println(decoded.BakongAccountID)
fmt.Println(decoded.MerchantName)
fmt.Println(decoded.TransactionAmount)
```

### Verify QR

```go
if err := khqr.Verify(qrString); err != nil {
    log.Fatal(err)
}
```

## API

| Function                                            | Description                                      |
| --------------------------------------------------- | ------------------------------------------------ |
| `GenerateIndividual(IndividualInfo) (*Data, error)` | Generate a KHQR string for an individual payment |
| `GenerateMerchant(MerchantInfo) (*Data, error)`     | Generate a KHQR string for a merchant payment    |
| `Decode(string) (*DecodedData, error)`              | Parse a KHQR string into structured data         |
| `Verify(string) error`                              | Validate CRC and structure of a KHQR string      |

## IndividualInfo

### Required Fields

| Field             | Type     | Max Length | Description                           |
| ----------------- | -------- | ---------- | ------------------------------------- |
| `BakongAccountID` | `string` | 32         | Must contain `@` (e.g. `"user@bank"`) |
| `MerchantName`    | `string` | 25         | Name of the individual                |

### Optional Fields with Defaults

| Field                  | Type       | Default        | Description              |
| ---------------------- | ---------- | -------------- | ------------------------ |
| `Currency`             | `Currency` | `KHR`          | `khqr.KHR` or `khqr.USD` |
| `MerchantCity`         | `string`   | `"Phnom Penh"` | Max 15 characters        |
| `MerchantCategoryCode` | `string`   | `"5999"`       | 1-4 digit numeric code   |

### Dynamic QR Fields

| Field                 | Type      | Description                                                     |
| --------------------- | --------- | --------------------------------------------------------------- |
| `Amount`              | `float64` | `0` = static QR; KHR must be whole number, USD up to 2 decimals |
| `ExpirationTimestamp` | `int64`   | Unix milliseconds, required when `Amount > 0`                   |

### Optional Fields

| Field                   | Type     | Max Length | Description                                              |
| ----------------------- | -------- | ---------- | -------------------------------------------------------- |
| `AcquiringBank`         | `string` | 32         | Acquiring bank name                                      |
| `AccountInfo`           | `string` | 32         | Additional account information                           |
| `UPIAccountInfo`        | `string` | 99         | UPI account info (not supported with USD)                |
| `BillNumber`            | `string` | 25         | Bill/invoice number                                      |
| `StoreLabel`            | `string` | 25         | Store identifier                                         |
| `TerminalLabel`         | `string` | 25         | Terminal identifier                                      |
| `MobileNumber`          | `string` | 25         | Mobile number                                            |
| `Purpose`               | `string` | 25         | Purpose of transaction                                   |
| `AltLanguagePreference` | `string` | 2          | ISO 639-1 code (e.g. `"km"`); requires `AltMerchantName` |
| `AltMerchantName`       | `string` | 25         | Required when `AltLanguagePreference` is set             |
| `AltMerchantCity`       | `string` | 15         | Alternate language city name                             |

## MerchantInfo

### Required Fields

| Field             | Type     | Max Length | Description                               |
| ----------------- | -------- | ---------- | ----------------------------------------- |
| `BakongAccountID` | `string` | 32         | Must contain `@` (e.g. `"merchant@bank"`) |
| `MerchantName`    | `string` | 25         | Name of the merchant                      |
| `MerchantCity`    | `string` | 15         | City of the merchant                      |
| `MerchantID`      | `string` | 32         | Merchant identifier                       |
| `AcquiringBank`   | `string` | 32         | Acquiring bank name                       |

### Optional Fields with Defaults

| Field                  | Type       | Default  | Description              |
| ---------------------- | ---------- | -------- | ------------------------ |
| `Currency`             | `Currency` | `KHR`    | `khqr.KHR` or `khqr.USD` |
| `MerchantCategoryCode` | `string`   | `"5999"` | 1-4 digit numeric code   |

### Dynamic QR Fields

| Field                 | Type      | Description                                                     |
| --------------------- | --------- | --------------------------------------------------------------- |
| `Amount`              | `float64` | `0` = static QR; KHR must be whole number, USD up to 2 decimals |
| `ExpirationTimestamp` | `int64`   | Unix milliseconds, required when `Amount > 0`                   |

### Optional Fields

| Field                   | Type     | Max Length | Description                                              |
| ----------------------- | -------- | ---------- | -------------------------------------------------------- |
| `UPIAccountInfo`        | `string` | 99         | UPI account info (not supported with USD)                |
| `BillNumber`            | `string` | 25         | Bill/invoice number                                      |
| `StoreLabel`            | `string` | 25         | Store identifier                                         |
| `TerminalLabel`         | `string` | 25         | Terminal identifier                                      |
| `MobileNumber`          | `string` | 25         | Mobile number                                            |
| `Purpose`               | `string` | 25         | Purpose of transaction                                   |
| `AltLanguagePreference` | `string` | 2          | ISO 639-1 code (e.g. `"km"`); requires `AltMerchantName` |
| `AltMerchantName`       | `string` | 25         | Required when `AltLanguagePreference` is set             |
| `AltMerchantCity`       | `string` | 15         | Alternate language city name                             |

## Static vs Dynamic QR

- **Static QR** (`Amount = 0`): No amount encoded, reusable for multiple payments.
- **Dynamic QR** (`Amount > 0`): Includes a specific amount and requires an `ExpirationTimestamp` (unix milliseconds).

## Currency

| Constant   | ISO 4217 Code | Rules                  |
| ---------- | ------------- | ---------------------- |
| `khqr.KHR` | 116           | Whole numbers only     |
| `khqr.USD` | 840           | Up to 2 decimal places |

## Error Handling

All errors are of type `*khqr.Error` with a `Code` and `Message`. Use `errors.Is` for comparison:

```go
if errors.Is(err, khqr.ErrAccountIDRequired) {
    // handle missing account ID
}
```

## License

MIT
