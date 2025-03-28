package main

import (
	"fmt"
	"lesson-02/utils"
)

func main() {
	fmt.Println(utils.FibonacciIterative(6)) // 8
	fmt.Println(utils.FibonacciRecursive(7)) // 13
	fmt.Println()

	fmt.Println(utils.IsPrime(3))  // true
	fmt.Println(utils.IsPrime(4))  // false
	fmt.Println(utils.IsPrime(7))  // true
	fmt.Println(utils.IsPrime(13)) // true
	fmt.Println(utils.IsPrime(15)) // false
	fmt.Println()

	fmt.Println(utils.IsBinaryPalindrome(7)) // true
	fmt.Println(utils.IsBinaryPalindrome(5)) // true
	fmt.Println(utils.IsBinaryPalindrome(6)) // false
	fmt.Println()

	fmt.Println(utils.ValidParentheses("{}"))     // true
	fmt.Println(utils.ValidParentheses("[{}]"))   // true
	fmt.Println(utils.ValidParentheses("[{]}"))   // false
	fmt.Println(utils.ValidParentheses("]}{["))   // false
	fmt.Println(utils.ValidParentheses("[{}"))    // false
	fmt.Println(utils.ValidParentheses("{[()]}")) // true
	fmt.Println()

	fmt.Println(utils.Increment("101"))  // 6
	fmt.Println(utils.Increment("111"))  // 8
	fmt.Println(utils.Increment("0001")) // 2
	fmt.Println(utils.Increment("1010")) // 11
	fmt.Println()

}
