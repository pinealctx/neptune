package tex

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	MaxDecimals = 15 // Maximum supported decimal places
)

var (
	powers = [16]float64{
		1,                // 10^0
		10,               // 10^1
		100,              // 10^2
		1000,             // 10^3
		10000,            // 10^4
		100000,           // 10^5
		1000000,          // 10^6
		10000000,         // 10^7
		100000000,        // 10^8
		1000000000,       // 10^9
		10000000000,      // 10^10
		100000000000,     // 10^11
		1000000000000,    // 10^12
		10000000000000,   // 10^13
		100000000000000,  // 10^14
		1000000000000000, // 10^15
	}

	powerDecimals = [16]decimal.Decimal{
		decimal.NewFromInt(1),                // 10^0
		decimal.NewFromInt(10),               // 10^1
		decimal.NewFromInt(100),              // 10^2
		decimal.NewFromInt(1000),             // 10^3
		decimal.NewFromInt(10000),            // 10^4
		decimal.NewFromInt(100000),           // 10^5
		decimal.NewFromInt(1000000),          // 10^6
		decimal.NewFromInt(10000000),         // 10^7
		decimal.NewFromInt(100000000),        // 10^8
		decimal.NewFromInt(1000000000),       // 10^9
		decimal.NewFromInt(10000000000),      // 10^10
		decimal.NewFromInt(100000000000),     // 10^11
		decimal.NewFromInt(1000000000000),    // 10^12
		decimal.NewFromInt(10000000000000),   // 10^13
		decimal.NewFromInt(100000000000000),  // 10^14
		decimal.NewFromInt(1000000000000000), // 10^15
	}

	strDecimals = [16]string{
		"%.0f",
		"%.1f",
		"%.2f",
		"%.3f",
		"%.4f",
		"%.5f",
		"%.6f",
		"%.7f",
		"%.8f",
		"%.9f",
		"%.10f",
		"%.11f",
		"%.12f",
		"%.13f",
		"%.14f",
		"%.15f",
	}
)

// Float642String converts a float64 value to a string representation
func Float642String(fv float64, decimalPlaces int) (string, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[string](decimalPlaces)
	}
	return fmt.Sprintf(strDecimals[decimalPlaces], fv), nil
}

// Float64ToDecimalV1 converts a float64 value to a decimal.Decimal representation
func Float64ToDecimalV1(fv float64, decimalPlaces int) (decimal.Decimal, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[decimal.Decimal](decimalPlaces)
	}
	return decimal.NewFromFloatWithExponent(fv, -int32(decimalPlaces)), nil
}

// Float64ToDecimalV2 converts a float64 value to a decimal.Decimal representation
func Float64ToDecimalV2(fv float64, decimalPlaces int) (decimal.Decimal, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[decimal.Decimal](decimalPlaces)
	}
	str := fmt.Sprintf(strDecimals[decimalPlaces], fv)
	return decimal.NewFromString(str)
}

// Float64ToIntDecimalV1 converts a float64 value to an int64 representation
// with the specified number of decimal places.
func Float64ToIntDecimalV1(fv float64, decimalPlaces int) (int64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[int64](decimalPlaces)
	}

	factor := powers[decimalPlaces]
	scaled := fv * factor

	// Proper rounding for both positive and negative numbers
	if scaled >= 0 {
		return int64(scaled + 0.5), nil
	}
	return int64(scaled - 0.5), nil
}

// Float64ToIntDecimalV2 converts a float64 value to an int64 representation
// with the specified number of decimal places.
func Float64ToIntDecimalV2(fv float64, decimalPlaces int) (int64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[int64](decimalPlaces)
	}
	decimalA := decimal.NewFromFloat(fv)
	decimalB := decimalA.Mul(powerDecimals[decimalPlaces])
	rounded := decimalB.Round(0)
	return rounded.IntPart(), nil
}

// Float64ToIntDecimalV3 converts a float64 value to an int64 representation
// with the specified number of decimal places.
func Float64ToIntDecimalV3(fv float64, decimalPlaces int) (int64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[int64](decimalPlaces)
	}
	str := fmt.Sprintf(strDecimals[decimalPlaces], fv)
	// remove the decimal point
	str = strings.ReplaceAll(str, ".", "")
	return strconv.ParseInt(str, 10, 64)
}

// IntDecimalToFloat64V1 converts an int64 representation back to a float64
func IntDecimalToFloat64V1(iv int64, decimalPlaces int) (float64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[float64](decimalPlaces)
	}

	factor := powers[decimalPlaces]
	return float64(iv) / factor, nil
}

// IntDecimalToFloat64V2 converts an int64 representation back to a float64
func IntDecimalToFloat64V2(iv int64, decimalPlaces int) (float64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[float64](decimalPlaces)
	}
	factor := decimal.New(iv, -int32(decimalPlaces))
	r, _ := factor.Float64() // 'exact' is not important here.
	return r, nil
}

// IntDecimalToFloat64V3 converts an int64 representation back to a float64
func IntDecimalToFloat64V3(iv int64, decimalPlaces int) (float64, error) {
	if decimalPlaces < 0 || decimalPlaces > MaxDecimals {
		return GenDecimalPlacesError[float64](decimalPlaces)
	}
	intStr := strconv.FormatInt(iv, 10)
	if decimalPlaces == 0 {
		return strconv.ParseFloat(intStr, 64)
	}
	// Insert decimal point
	intStr = intStr[:len(intStr)-decimalPlaces] + "." + intStr[len(intStr)-decimalPlaces:]
	return strconv.ParseFloat(intStr, 64)
}

// GenDecimalPlacesError generates an error for invalid decimal places
func GenDecimalPlacesError[T any](decimalPlaces int) (T, error) {
	var t T
	return t, fmt.Errorf("decimalPlaces:%d out of range [0, 15]", decimalPlaces)
}
