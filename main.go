package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/lucasjones/reggen"
)

const (
	or       = 0
	multiple = 1
	and      = 2
)

var (
	filePath string
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

func crt(ruleset []Rule) (int, bool, []interface{}, error) {
	breakPoint := 0
	atLeastMatch := true
	final := []interface{}{}
	var err error

	tk_IDX := 0

	for _, rule := range ruleset {
		matches := false
		var toCompare Token

		if rule.Raw {
			_matches, _toCompare := compareNext(rule.Type)
			toCompare = _toCompare
			matches = _matches
			fmt.Println("Val :", toCompare.Value, matches)
			fmt.Println(toCompare)
			if matches {
				final = append(final, toCompare.Value)
			}
			tk_IDX += 1
		} else {
			_final, err := consume_type(rule.Type, rule)
			matches = err == nil
			fmt.Println("OR MATCH:", matches, err)
			final = append(final, _final)
		}

		breakPoint = tk_IDX
		if !matches {
			atLeastMatch = false
			fmt.Println("Error :", toCompare)
			val := strings.ReplaceAll(toCompare.Value, "\n", "\\n")
			plural := ""

			if len(val) == 1 {
				plural = "token"
			} else {
				plural = "tokens"
			}
			str, _ := reggen.Generate(Matches[rule.Type], 10)
			err = fmt.Errorf("Unexpected %s '%s' at %s:%d:%d  expected something like: %s", plural, val, filePath, toCompare.start[1], toCompare.start[0], str)
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

func consume(rootRule Rule) (interface{}, error) {
	var err error
	if rootRule.Raw {

		_, matches, final, _err := crt(rootRule.Subrules)
		if err != nil {
			err = _err
		}
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

func consume_type(_type int, rootRule Rule) (interface{}, error) {
	var err error
	switch _type {
	case multiple:
		x := 0
		running := true
		final := []interface{}{}

		var tempErr error
		for running {
			_, matches, _final, _err := crt(rootRule.Subrules)
			final = append(final, _final...)
			running = matches
			tempErr = _err
			if !matches {
				break
			}
			x += 1
		}

		if x == 0 {
			err = tempErr
		}
		return rootRule.Callback(final), err
	case or:
		match := false
		final := []interface{}{}
		var tempErr error

		for _, rule := range rootRule.Subrules {
			_, matches, _final, _err := crt([]Rule{rule})
			match = matches
			tempErr = _err
			if match {
				final = append(final, _final...)
				return rootRule.Callback(final), err
				break
			}
		}
		if !match {
			err = tempErr
		} else {
		}
	case and:
		_, _, final, err := crt(rootRule.Subrules)
		if err == nil {
			return rootRule.Callback(final), err
		}
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

	filePath = "./code.js"
	body, err := ioutil.ReadFile("./code.js")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	ToParse = string(body)

	s := rawRule(str)
	num := rawRule(number)

	dataTypes := []Rule{num, s}

	expression := rule_cb(or, dataTypes, func(args interface{}) interface{} {
		first := args.([]interface{})[0].(string)[0]
		fmt.Println("First : ", args)
		if first == '"' || first == '\'' || first == '`' {
			return literal_str(interfaces_to_str(args))
		} else {
			return literal_num(interfaces_to_str(args))
		}
	})

	assignementStructure := []Rule{rawRule(var_var), rawRule(equals), expression}

	assignement := rule_cb(and, assignementStructure, func(args interface{}) interface{} {
		return "Assignement"
	})

	instructions := []Rule{expression, assignement}

	instruction := rule_cb(multiple, instructions, func(args interface{}) interface{} {
		return args
	})

	final, err := consume(instruction)

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Final :", final)
	}
}
