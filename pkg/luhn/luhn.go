package luhn

func ValidateString(val string) bool {
	sum := 0

	isSecondDigit := false

	for i := len(val) - 1; i >= 0; i-- {
		digit := int(val[i] - '0')

		if digit < 0 || digit > 9 {
			return false
		}

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		isSecondDigit = !isSecondDigit

		sum += digit
	}

	return sum%10 == 0 && sum > 0
}
