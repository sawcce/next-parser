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
	return Literal{"Number", "", num}
}

func crt(ruleset []Rule) (int, bool, []interface{}, string) {
	breakPoint := 0
	atLeastMatch := true
	final := []interface{}{}
	err := ""

	tk_IDX := 0

	for _, rule := range ruleset {
		matches := false
		var toCompare Token

		if rule.Raw {
			_matches, toCompare := compareNext(rule.Type)
			matches = _matches
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
			atLeastMatch = false
			err = fmt.Sprintf("Unexpected token '%s' at column %d, expected : '%s' \n", toCompare.Value, toCompare.start, Tokens[rule.Type])
			fmt.Println("no match")
			break
		}
	}

	//if breakPoint > 0 {
	//	fmt.Println(breakPoint)
	//	atLeastMatch = true
	//}

	return breakPoint, atLeastMatch, final, err
}

func consume(rootRule Rule) (interface{}, string) {
	var err string
	if rootRule.Raw {

		_, matches, final, _err := crt(rootRule.Subrules)
		err = _err
		if matches {
			return rootRule.Callback(final), err
		} else {
			return "", err
		}
	} else {
		final, _err := consume_type(rootRule.Type, rootRule)
		return final, _err
	}
}

func consume_type(_type int, rootRule Rule) (interface{}, string) {
	var err string
	switch _type {
	case multiple:
		x := 0
		running := true
		final := []interface{}{}

		for running {
			bp, matches, _final, _err := crt(rootRule.Subrules)
			final = append(final, _final...)
			running = matches
			x += bp
			err = _err
		}
		return rootRule.Callback(final), err
	case or:
		match := false
		final := []interface{}{}
		tempErr := ""

		for _, rule := range rootRule.Subrules {
			_, matches, _final, _err := crt([]Rule{rule})
			fmt.Println("Rule :", rule, "Match", matches)
			match = matches
			tempErr = _err
			if match {
				final = append(final, _final...)
				break
			}
		}
		if !match {
			err = tempErr
		}
		return rootRule.Callback(final), err
		break
	case and:
		_, _, final, err := crt(rootRule.Subrules)
		return rootRule.Callback(final), err
		break
	}
	return "", err
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

	ToParse = `54602`

	num := rawRule(number)
	s := rawRule(str)

	dataTypes := []Rule{num, s}

	exp := rule_cb(or, dataTypes, func(args interface{}) interface{} {
		first := args.([]interface{})[0].(string)[0]
		if first == '"' || first == '\'' || first == '`' {
			return literal_str(interfaces_to_str(args))
		} else {
			return literal_num(interfaces_to_str(args))
		}
	})

	final, err := consume(exp)

	fmt.Println("Final", final.(Literal).valueType, ":", final, err)
}
