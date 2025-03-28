package utils

func IsPrime(n int) bool {
	// Функція повертає `true` якщо число `n` - просте.
	// Інакше функція повертає `false`
	if n < 2 {
		return false
	}

	for i := 2; i <= n; i++ {
		if n%i == 0 && n != i {
			return false
		}
	}

	return true
}
