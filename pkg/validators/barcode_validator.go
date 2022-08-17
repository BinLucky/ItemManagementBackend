package validators

import "regexp"

func IsBarcodeValid(barcode string) bool {

	var rxBarcode = regexp.MustCompile("^[0-9]+$")
	runeInBarcode := []rune(barcode)

	if len(runeInBarcode) != 8 || rxBarcode.MatchString(barcode) {
		return false
	}
	return true
}
