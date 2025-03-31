package utils

var brackets = map[byte]byte{
	')': '(',
	']': '[',
	'}': '{',
}

func ValidParenthesesMap(s string) bool {
	stack := []byte{}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		switch ch {
		case '(', '[', '{':
			stack = append(stack, ch)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != brackets[ch] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func ValidParentheses(s string) bool {
	// Функція повертає `true` якщо у вхідній стрічці дотримані усі правила високристання дужок
	// Правила:
	// 1. Допустимі дужки `(`, `[`, `{`, `)`, `]`, `}`
	// 2. У кожної відкритої дужки є відповідна закриваюча дужка того ж типу
	// 3. Закриваючі дужки стоять у правильному порядку
	//    "[{}]" - правильно
	//    "[{]}" - не правильно
	// 4. Кожна закриваюча дужка має відповідну відкриваючу дужку
	stack := []byte{}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		switch ch {
		case '(', '[', '{':
			stack = append(stack, ch)
		case ')', ']', '}':
			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]

			if ch-top != 1 && ch-top != 2 {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
