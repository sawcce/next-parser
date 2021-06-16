package main

import "fmt"

var (
	Tokens = []string{"(", ")", "dig", "."}
)

const (
	lparen = 0
	rparen = 1
	digit  = 2
	dot    = 3

	or       = 0
	multiple = 1
)

type Token struct {
	Type  int
	Value string
	start int
	end   int
}

func token(_type int, value string, start int, end int) Token {
	return Token{_type, value, start, end}
}

func TK(_type int) Token {
	return Token{_type, "0", 0, 0}
}

type Rule struct {
	Raw      bool
	Type     int
	Subrules []Rule
}

func rawRule(_type int) Rule {
	return Rule{true, _type, []Rule{}}
}

func rule(_type int, rules []Rule) Rule {
	return Rule{false, _type, rules}
}

func consume_ruleset_token(ruleset []Rule, tokens []Token) (int, bool) {
	matches := false
	breakPoint := 0

	x := 0

	fmt.Println(">> Entering with the ruleset :", ruleset)
	fmt.Println(">> Entering with the tokens :", tokens)

	for i, rule := range ruleset {
		fmt.Println("[", i, "]", "[", x, "]")
		breakPoint = i
		compare := tokens[x]

		if rule.Raw == true {
			fmt.Println("[Raw] with the rule being :", Tokens[rule.Type], "and the compared :", Tokens[compare.Type])
			x += 1
			matches = rule.Type == compare.Type
		} else {
			switch rule.Type {
			case or:
				break
			case multiple:
				fmt.Println("[SubRule] Multiple x:", x)
				running := true
				fbp := 0
				bp := x

				for running == true {
					_breakPoint, match := consume_ruleset_token(rule.Subrules, tokens[bp:])
					running = match
					bp += _breakPoint
					fbp += _breakPoint
					fmt.Print("|", _breakPoint, "|")
				}

				fmt.Print("\n", fbp)

				matches = (fbp != 0)

				x += fbp
				breakPoint += fbp

				fmt.Println("[Mutiple : Finished executing, current token idx:", x, ", current rule idx:", i)
				fmt.Println("[Mutiple : did it match?", matches)
				fmt.Println("[Mutiple : first break point", fbp)
				break
			}
		}

		if !matches {
			fmt.Println(">> Exit because not matching.", matches)
			fmt.Println("[Exit]", x, i)
			break
		}

	}
	if len(ruleset) == 1 && matches {
		breakPoint = 1
	}

	fmt.Println(">> Exit")
	fmt.Println(x, len(ruleset))
	return breakPoint, matches
}

func main() {
	//expression := [] Rule{rawRule(lparen), rawRule(rparen)}
	//toMatch := [] Token{token(lparen, "(", 0, 0), token(rparen, ")", 1, 1)}

	multipleDigits := rule(multiple, []Rule{rawRule(digit)})

	toTest := []Token{TK(digit), TK(digit), TK(dot), TK(digit), TK(digit), TK(dot)}
	number := []Rule{multipleDigits, rawRule(dot), multipleDigits}

	breakPoint, matches := consume_ruleset_token(number, toTest)

	fmt.Println("[Result] Does match ?", matches, ", breakPoint :", breakPoint)
}
