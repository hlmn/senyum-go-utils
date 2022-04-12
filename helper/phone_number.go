package helper

import "strings"

func convertPhoneNumberTo08(phoneNumber *string) {
	if phoneNumber == nil {
		return
	}
	var phoneNumberTemp = *phoneNumber
	if len(phoneNumberTemp) == 0 {
		return
	}
	if phoneNumberTemp[0:2] == "62" {
		*phoneNumber = strings.Replace(phoneNumberTemp, "62", "0", 1)
		return
	}
	if phoneNumberTemp[0:3] == "+62" {
		*phoneNumber = strings.Replace(phoneNumberTemp, "+62", "0", 1)
		return
	}

	return
}
