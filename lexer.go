package main

import (
	"fmt"
	"regexp"

	"github.com/lucasjones/reggen"
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

	// Data Types
	identifier = 30

	// Variables
	const_var = 40
	let_var   = 41
	var_var   = 42

	//
	eof       = 1000
	unmatched = 1050
)

var (
	ToParse         = ""
	currentIndex    = 0 // Current of the lexer
	rememberedIndex = 0 // Used on set / pop
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

	Matches[identifier] = `([a-zA-Z])+`

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

func generateExample(_type int) string {
	str, _ := reggen.Generate(Matches[_type], 10)
	return str
}

func setState() {
	rememberedIndex = currentIndex
}

func popState() {
	currentIndex = rememberedIndex
}

func compareNext(expected int) (Token, error) {
	var err error
	s := ToParse[currentIndex:]
	r, _ := regexp.Compile(Matches[expected])
	l := r.FindSubmatchIndex([]byte(s))

	if len(l) == 0 {
		return token(unmatched, s[currentIndex:currentIndex+1], [2]int{0, 0}, [2]int{0, 0}), fmt.Errorf("%s", generateExample(expected))
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
		return token(expected, match, startPos, endPos), err
	} else {
		startPos[0] -= 1
		return token(unmatched, s[:l[0]], startPos, startPos), fmt.Errorf("")
	}
}

func test(str string, reg string, opts ...string) (string, bool) {
	r, _ := regexp.Compile(reg)
	l := r.FindIndex([]byte(str))

	fmt.Println(l)

	return str[l[0]:l[1]], false

	return "", true
}
