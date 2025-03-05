package algo

import (
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// Address is the Algorand address to fetch rewards for.
	Address = ""
	// printer is a message printer for formatting numbers.
	printer = message.NewPrinter(message.MatchLanguage(language.Make(os.Getenv("LANG")).String()))
)

// ShortAddress returns a shortened version of the address.
func ShortAddress() string {
	if len(Address) < 8 {
		return Address
	}
	return Address[:4] + "..." + Address[len(Address)-4:]
}

// ClearCache clears the cache files.
func ClearCache() {
	ClearRewardsCache()
}

// FormatFloat formats a float64 as a string.
func FormatFloat(f float64) string {
	return printer.Sprintf("%.6f", f)
}

// FormatInt formats an int64 as a string.
func FormatInt(i int64) string {
	return printer.Sprintf("%d", i)
}
