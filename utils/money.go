package utils

import (
	"errors"
	"math/big"
	"strings"
)

// maxAmount mirrors NUMERIC(15,2): 13 integer digits + 2 decimals.
var maxAmount = new(big.Rat).SetFrac64(9999999999999_99, 100)

// ParseAmount validates a monetary string and returns it normalized to two
// decimal places, so the value handed to a NUMERIC(15,2) column is exact.
// NUMERIC columns are mapped to Go strings by sqlc, which is why amounts are
// kept as strings end to end instead of going through float64.
func ParseAmount(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("jumlah harus diisi")
	}

	rat, ok := new(big.Rat).SetString(raw)
	if !ok {
		return "", errors.New("jumlah tidak valid")
	}
	if rat.Sign() <= 0 {
		return "", errors.New("jumlah harus lebih besar dari 0")
	}
	if rat.Cmp(maxAmount) > 0 {
		return "", errors.New("jumlah terlalu besar")
	}

	// Reject sub-cent precision rather than silently rounding the user's input.
	cents := new(big.Rat).Mul(rat, new(big.Rat).SetInt64(100))
	if !cents.IsInt() {
		return "", errors.New("jumlah maksimal 2 angka desimal")
	}

	return rat.FloatString(2), nil
}

// ParseSignedAmount behaves like ParseAmount but allows zero and negative
// values, for fields such as an account's opening balance.
func ParseSignedAmount(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "0.00", nil
	}

	rat, ok := new(big.Rat).SetString(raw)
	if !ok {
		return "", errors.New("jumlah tidak valid")
	}
	if new(big.Rat).Abs(rat).Cmp(maxAmount) > 0 {
		return "", errors.New("jumlah terlalu besar")
	}

	cents := new(big.Rat).Mul(rat, new(big.Rat).SetInt64(100))
	if !cents.IsInt() {
		return "", errors.New("jumlah maksimal 2 angka desimal")
	}

	return rat.FloatString(2), nil
}

// NegateAmount flips the sign of a normalized amount string, used to reverse
// a balance effect without re-parsing the value as a number.
func NegateAmount(amount string) string {
	amount = strings.TrimSpace(amount)
	if amount == "" || amount == "0" || amount == "0.00" {
		return "0.00"
	}
	if strings.HasPrefix(amount, "-") {
		return strings.TrimPrefix(amount, "-")
	}
	return "-" + amount
}

// SubAmount returns a - b for two decimal strings, formatted to two decimals.
// Invalid input is treated as zero, since both operands come from NUMERIC
// columns rather than from user input.
func SubAmount(a, b string) string {
	ra, okA := new(big.Rat).SetString(strings.TrimSpace(a))
	if !okA {
		ra = new(big.Rat)
	}
	rb, okB := new(big.Rat).SetString(strings.TrimSpace(b))
	if !okB {
		rb = new(big.Rat)
	}
	return new(big.Rat).Sub(ra, rb).FloatString(2)
}
