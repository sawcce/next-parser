package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

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

func def(args interface{}) interface{} {
	return args.([]interface{})[0]
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

func multiple_rule(rules []Rule) Rule {
	return Rule{false, multiple, rules, func(args interface{}) interface{} { return args.([]interface{})[0] }}
}

func or_rules(rules []Rule) Rule {
	return Rule{false, or, rules, func(args interface{}) interface{} { return args.([]interface{})[0] }}
}

func tokens_to_strings(tokens []Token) string {
	res := ""

	for _, token := range tokens {
		res += token.Value
	}

	return res
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
			atLeastMatch = _matches
			if matches {
				final = append(final, toCompare.Value)
			}
			tk_IDX += 1
		} else {
			_final, err := consume_type(rule.Type, rule)
			matches = err == nil
			atLeastMatch = err == nil
			final = append(final, _final)
		}

		breakPoint = tk_IDX
		if !matches {
			atLeastMatch = false
			val := strings.ReplaceAll(toCompare.Value, "\n", "\\n")
			plural := ""

			if len(val) <= 1 {
				plural = "token"
			} else {
				plural = "tokens"
			}

			str, _ := reggen.Generate(Matches[rule.Type], 10)
			err = fmt.Errorf("Unexpected %s '%s' at %s:%d:%d  expected something like: %s", plural, val, filePath, toCompare.start[1], toCompare.start[0], str)

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
		fmt.Println("CONSUME")
		final, _err := consume_type(rootRule.Type, rootRule)
		fmt.Println("END CONSUME", _err)
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
			fmt.Println("Running", matches, currentIndex, final)
			final = append(final, _final...)
			running = matches
			tempErr = _err
			if !matches {
				if x == 0 {
					err = tempErr
					return final, err
				}
			}
			x += 1
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
		err = tempErr
		return "", tempErr
	case and:
		_, _, final, err := crt(rootRule.Subrules)
		if err == nil {
			return rootRule.Callback(final), err
		} else {
			return "", err
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

	start := time.Now()

	filePath = "./code.js"
	body, err := ioutil.ReadFile("./code.js")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	ToParse = string(body)

	s := rawRule(str)
	num := rawRule(number)

	stringRule := rule_cb(and, []Rule{s}, func(args interface{}) interface{} {
		return literal_str(interfaces_to_str(args))
	})

	numRule := rule_cb(and, []Rule{num}, func(args interface{}) interface{} {
		return literal_num(interfaces_to_str(args))
	})

	//fp := rule(multiple, []Rule{rawRule(identifier), rawRule(dot)})
	//sp := rule(and, []Rule{fp, rawRule(identifier)})
	//accessor := or_rules([]Rule{rawRule(identifier)})

	expressions := rule(or, []Rule{stringRule, numRule, rawRule(identifier)})

	assignementStructure := []Rule{rawRule(const_var), rawRule(identifier), rawRule(equals), expressions}

	assignement := rule_cb(and, assignementStructure, func(args interface{}) interface{} {
		params := args.([]interface{})
		fmt.Println("Assignement :", params, params[1].(string))
		return var_declaration(params[1].(string), params[0].(string), params[3])
	})

	instructions := []Rule{assignement}

	instruction := rule_cb(and, instructions, func(args interface{}) interface{} {
		return args
	})

	program := rule_cb(multiple, []Rule{instruction}, func(args interface{}) interface{} {
		return args
	})

	final, err := consume(program)

	if err != nil {
		log.Fatalln(err)
	} else {
		//d, _ := json.Marshal(final)
		fmt.Println("fn :", final)
	}

	elapsed := time.Since(start)
	log.Printf("Reading, Parsing / Lexing took %s", elapsed)
}
