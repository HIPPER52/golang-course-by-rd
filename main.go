package main

import (
	"fmt"
	"lesson-02/utils"
	"time"
)

func main() {
	fmt.Println(utils.FibonacciIterative(6))  // 8
	fmt.Println(utils.FibonacciIterative(-3)) // -1
	fmt.Println(utils.FibonacciRecursive(7))  // 13
	fmt.Println()

	start := time.Now()
	fmt.Println(utils.IsPrime(3))                       // true
	fmt.Println(utils.IsPrime(4))                       // false
	fmt.Println(utils.IsPrime(17))                      // true
	fmt.Println(utils.IsPrime(97))                      // true
	fmt.Println(utils.IsPrime(100))                     // false
	fmt.Println("PrimeOld-version:", time.Since(start)) // 46ms for the first time, 14.792ms for the second time
	fmt.Println()

	start = time.Now()
	fmt.Println(utils.IsPrimeNew(3))                    // true
	fmt.Println(utils.IsPrimeNew(4))                    // false
	fmt.Println(utils.IsPrimeNew(17))                   // true
	fmt.Println(utils.IsPrimeNew(97))                   // true
	fmt.Println(utils.IsPrimeNew(100))                  // false
	fmt.Println("PrimeNew-version:", time.Since(start)) // 9.292ms for the first time, 2.917ms for the second time
	fmt.Println()

	fmt.Println(utils.IsBinaryPalindrome(7)) // true
	fmt.Println(utils.IsBinaryPalindrome(5)) // true
	fmt.Println(utils.IsBinaryPalindrome(6)) // false
	fmt.Println()

	start = time.Now()
	fmt.Println(utils.ValidParentheses("{}"))        // true
	fmt.Println(utils.ValidParentheses("[{}]"))      // true
	fmt.Println(utils.ValidParentheses("[{]}"))      // false
	fmt.Println(utils.ValidParentheses("]}{["))      // false
	fmt.Println(utils.ValidParentheses("[{}"))       // false
	fmt.Println(utils.ValidParentheses("{[()]}"))    // true
	fmt.Println("ASCII-version:", time.Since(start)) // 22.834ms for the first time, 3.958ms for the second time
	fmt.Println()

	start = time.Now()
	fmt.Println(utils.ValidParenthesesMap("{}"))     // true
	fmt.Println(utils.ValidParenthesesMap("[{}]"))   // true
	fmt.Println(utils.ValidParenthesesMap("[{]}"))   // false
	fmt.Println(utils.ValidParenthesesMap("]}{["))   // false
	fmt.Println(utils.ValidParenthesesMap("[{}"))    // false
	fmt.Println(utils.ValidParenthesesMap("{[()]}")) // true
	fmt.Println("Map-version:", time.Since(start))   // 6.625ms for the first time, 4.334ms for the second time
	fmt.Println()

	fmt.Println(utils.Increment("101"))  // 6
	fmt.Println(utils.Increment("111"))  // 8
	fmt.Println(utils.Increment("0001")) // 2
	fmt.Println(utils.Increment("1010")) // 11
	fmt.Println()

}
