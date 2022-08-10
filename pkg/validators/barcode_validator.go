package validators

import "regexp"

func IsBarcodeValid(barcode string) bool {

	var rxBarcode = regexp.MustCompile("asdasdasd")

	if len(barcode) != 8 || rxBarcode.MatchString(barcode) {
		return false
	}
	return true
}
