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

	// Native Data Types
	number = 10
	str    = 11

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

	Matches map[int]string = make(map[int]string)
)

func initMap() {
	Matches[lparen] = "\\("
	Matches[rparen] = "\\)"
	Matches[dot] = "\\."

	Matches[number] = `([1-9]\d*(\.\d*[1-9])?|0\.\d*[1-9]+)|\d+(\.\d*[1-9])?`
	Matches[str] = `\".*?\"`

	Tokens[lparen] = "("
	Tokens[rparen] = ")"
	Tokens[dot] = "."

	Tokens[number] = `a number`
	Tokens[str] = `a "String"`
}

type Token struct {
	Type  int
	Value string
	start int
	end   int
}

func token(_type int, value string, start int, end int) Token {
	return Token{_type, value, start, end}
}

func token_type_val(_type int, value string) Token {
	return Token{_type, value, 0, 0}
}

func TK(_type int) Token {
	return Token{_type, "0", 0, 0}
}

func compareNext(expected int) (bool, Token) {
	s := ToParse[currentIndex:]
	r, _ := regexp.Compile(Matches[expected])
	l := r.FindSubmatchIndex([]byte(s))

	none, _ := regexp.Compile(`((\ )|(\n))+`)
	m := l[0] == 0 || none.MatchString(s[0:l[0]])

	fmt.Println(l[0], l[1])

	if m {
		match := s[l[0]:l[1]]
		currentIndex += l[1]
		return true, token(expected, match, currentIndex+l[0], currentIndex+l[1])
	} else {
		return false, token(eof, s[:l[0]], currentIndex, currentIndex)
	}
}

func test(str string, reg string, opts ...string) (string, bool) {
	r, _ := regexp.Compile(reg)
	l := r.FindIndex([]byte(str))

	fmt.Println(l)

	return str[l[0]:l[1]], false

	return "", true
}

func lex(s string) Token {
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
}
