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

func IsPrimeNew(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}
