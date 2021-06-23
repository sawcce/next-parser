package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	or       = 0
	multiple = 1
	and      = 2
)

type Rule struct {
	Raw      bool
	Type     int
	Subrules []Rule
	Callback func(interface{}) interface{}
}

func def(interface{}) interface{} {
	return ""
}

func rawRule(_type int) Rule {
	return Rule{true, _type, []Rule{}, def}
}

func rule(_type int, rules []Rule) Rule {
	return Rule{false, _type, rules, def}
}

func rule_cb(_type int, rules []Rule, callback func(interface{}) interface{}) Rule {
	return Rule{false, _type, rules, callback}
}

func tokens_to_strings(tokens []Token) string {
	res := ""

	for _, token := range tokens {
		res += token.Value
	}

	return res
}

type Literal struct {
	valueType string
	raw       string
	value     interface{}
}

func literal_str(value string) Literal {
	return Literal{"string", value, 0}
}

func literal_num(value string) Literal {
	num, _ := strconv.ParseFloat(value, 64)
	return Literal{"string", "", num}
}

func crt(ruleset []Rule) (int, bool, []interface{}) {
	breakPoint := 0
	atLeastMatch := false
	final := []interface{}{}

	tk_IDX := 0

	for _, rule := range ruleset {
		matches := false
		var toCompare Token

		if rule.Raw {
			matches, toCompare = compareNext(rule.Type)
			fmt.Println("Val :", toCompare.Value, matches)
			fmt.Println(toCompare)
			if matches {
				final = append(final, toCompare.Value)
			}
			tk_IDX += 1
		} else {
			consume_type(rule.Type, rule)
		}

		breakPoint = tk_IDX
		if !matches {
			fmt.Printf("Unexpected token '%s' at column %d, expected : '%s' \n", toCompare.Value, toCompare.start, Tokens[rule.Type])
			fmt.Println("no match")
			break
		}
	}

	if breakPoint > 0 {
		atLeastMatch = true
	}

	return breakPoint, atLeastMatch, final
}

func consume(rootRule Rule) interface{} {
	if rootRule.Raw {

		_, matches, final := crt(rootRule.Subrules)
		if matches {
			return rootRule.Callback(final)
		} else {
			return ""
		}
	} else {
		final := consume_type(rootRule.Type, rootRule)
		return final
	}
}

func consume_type(_type int, rootRule Rule) interface{} {
	switch _type {
	case multiple:
		x := 0
		running := true
		final := []interface{}{}

		for running {
			bp, matches, _final := crt(rootRule.Subrules)
			final = append(final, _final...)
			running = matches
			x += bp
		}
		return rootRule.Callback(final)
	case or:
		break
	case and:
		_, _, final := crt(rootRule.Subrules)
		return rootRule.Callback(final)
		break
	}
	return ""
}

func interfaces_to_str(i interface{}) string {
	fn := []string{}

	for _, arg := range i.([]interface{}) {
		fn = append(fn, arg.(string))
	}

	return strings.Join(fn, "")
}

func main() {
	initMap()
	ToParse = "  1247854567"

	expression := []Rule{rawRule(number)}

	exp := rule_cb(and, expression, func(args interface{}) interface{} {
		fmt.Println(literal_num(interfaces_to_str(args)))
		return literal_num(interfaces_to_str(args))
	})

	//toMatch := []Token{token(number, "1", 0, 0), token(number, "7", 1, 1)}

	final := consume(exp)

	fmt.Println("Final number:", final)

	//fmt.Println("Lexed :", lex("1247854567"))
}
