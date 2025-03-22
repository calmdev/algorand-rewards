package format

import (
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// printer is a message printer for formatting numbers.
	printer = message.NewPrinter(message.MatchLanguage(language.Make(os.Getenv("LANG")).String()))
)

// AddressShort returns a shortened version of the address.
func AddressShort(address string) string {
	if len(address) < 8 {
		return address
	}
	return address[:4] + "..." + address[len(address)-4:]
}

// Float formats a float64 as a string.
func Float(f float64) string {
	return printer.Sprintf("%.6f", f)
}

// FloatShort formats a float64 as a string with 3 decimal places.
func FloatShort(f float64) string {
	return printer.Sprintf("%.3f", f)
}

// Int formats an int64 as a string.
func Int(i int64) string {
	return printer.Sprintf("%d", i)
}
