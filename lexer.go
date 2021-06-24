package main

import (
	"fmt"
	"regexp"
)

const (
	// Punctuation
	lparen = 0
	rparen = 1
	dot    = 2
	equals = 3

	// Native Data Types
	number = 10
	str    = 11

	// Variables
	const_var = 20
	let_var   = 21
	var_var   = 22

	//
	eof           = 100
	notRecognized = 150
)

var (
	ToParse      = ""
	currentIndex = 0 // Current of the lexer
)

var (
	Tokens map[int]string = make(map[int]string)

	Matches         map[int]string = make(map[int]string)
	defaultTokenPos                = [2]int{0, 0}
)

func initMap() {
	Matches[lparen] = "\\("
	Matches[rparen] = "\\)"
	Matches[dot] = "\\."
	Matches[equals] = `\=`

	Matches[number] = `([1-9]\d*(\.\d*[1-9])?|0\.\d*[1-9]+)|\d+(\.\d*[1-9])?`
	Matches[str] = `\".*?\"`

	Matches[const_var] = "const"
	Matches[let_var] = "let"
	Matches[var_var] = "var"

	// STRING REPS

	Tokens[lparen] = "("
	Tokens[rparen] = ")"
	Tokens[dot] = "."

	Tokens[number] = `a number`
	Tokens[str] = `a "String"`

}

type Token struct {
	Type  int
	Value string
	start [2]int
	end   [2]int
}

func token(_type int, value string, start [2]int, end [2]int) Token {
	return Token{_type, value, start, end}
}

func token_type_val(_type int, value string) Token {
	return Token{_type, value, defaultTokenPos, defaultTokenPos}
}

func TK(_type int) Token {
	return Token{_type, "0", defaultTokenPos, defaultTokenPos}
}

func countLines(str string) int {
	none, _ := regexp.Compile(`\n`)
	m := none.FindAllStringSubmatch(str, -1)
	return len(m)
}

func countNothing(str string) int {
	none, _ := regexp.Compile(`[\n ]`)
	m := none.FindAllStringSubmatch(str, -1)
	return len(m)
}

func compareNext(expected int) (bool, Token) {
	s := ToParse[currentIndex:]
	r, _ := regexp.Compile(Matches[expected])
	l := r.FindSubmatchIndex([]byte(s))

	fmt.Println("Testing for:", Tokens[expected], currentIndex, s)

	if len(l) == 0 {
		fmt.Println("NO MATCH")
		return false, TK(eof)
	}

	none, _ := regexp.Compile(`\S`)
	d := none.MatchString(s[0:l[0]])
	m := l[0] == 0 || !d

	startCol := currentIndex + l[0]

	startPos := [2]int{startCol, countLines(ToParse[:startCol])}

	if m {
		match := s[l[0]:l[1]]
		endCol := currentIndex + l[1] + 0
		endPos := [2]int{endCol, countLines(ToParse[:endCol])}
		currentIndex += l[1]
		return true, token(expected, match, startPos, endPos)
	} else {
		startPos[0] -= 1
		return false, token(eof, s[:l[0]], startPos, startPos)
	}
}

func test(str string, reg string, opts ...string) (string, bool) {
	r, _ := regexp.Compile(reg)
	l := r.FindIndex([]byte(str))

	fmt.Println(l)

	return str[l[0]:l[1]], false

	return "", true
}

/*func lex(s string) Token {
	for i, str := range s {
		fmt.Println("Rune:", i, string(str), string('1'))
		switch str {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			num, _ := test(s[i:], "\\d*\\.?\\d*", "Number")
			fmt.Println("Number :", num)
			i += len(num)
			return token(number, num, i, len(s))
		case '(':
			return token(lparen, string(str), i, i)
		case ')':
			return token(rparen, string(str), i, i)
		default:
			return TK(notRecognized)
		}
	}
	return TK(eof)
}*/
