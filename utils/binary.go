package utils

import "strconv"

func IsBinaryPalindrome(n int) bool {
	// Функція повертає `true` якщо число `n` у бінарному вигляді є паліндромом
	// Інакше функція повертає `false`
	//
	// Приклади:
	// Число 7 (111) - паліндром, повертаємо `true`
	// Число 5 (101) - паліндром, повертаємо `true`
	// Число 6 (110) - не є паліндромом, повертаємо `false`
	var binary = strconv.FormatInt(int64(n), 2)
	var leftNumber = 0
	var rightNumber = len(binary) - 1

	for leftNumber < rightNumber {
		if binary[leftNumber] != binary[rightNumber] {
			return false
		}
		leftNumber++
		rightNumber--
	}

	return true
}

func Increment(num string) int {
	// Функція на вхід отримує стрічку яка складається лише з символів `0` та `1`
	// Тобто стрічка містить певне число у бінарному вигляді
	// Потрібно повернути число на один більше
	var result int = 0
	var temp int = 1

	for i := len(num) - 1; i >= 0; i-- {
		switch num[i] {
		case '1':
			result += temp
			temp *= 2
		case '0':
			temp *= 2
		}
	}

	return result + 1
}
