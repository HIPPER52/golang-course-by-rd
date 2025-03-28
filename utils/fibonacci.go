package utils

func FibonacciIterative(n int) int {
	// Функція вираховує і повертає n-не число фібоначчі
	// Імплементація без використання рекурсії
	if n == 0 {
		return 0
	}

	temp1, temp2 := 0, 1
	var next int

	for i := 2; i <= n; i++ {
		next = temp1 + temp2
		temp1 = temp2
		temp2 = next
	}

	return temp2
}

func FibonacciRecursive(n int) int {
	// Функція вираховує і повертає n-не число фібоначчі
	// Імплементація з використанням рекурсії
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
	}
}
