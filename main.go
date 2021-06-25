package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
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

func interfaces_to_str(i interface{}) string {
	fn := []string{}

	for _, arg := range i.([]interface{}) {
		fn = append(fn, arg.(string))
	}

	return strings.Join(fn, "")
}

func or_types(types []int) (Token, error) {
	eac := errorAcc()
	var token Token
	for _, tp := range types {
		tk, _err := compareNext(tp)
		if _err == nil {
			return tk, _err
		}
		eac.accumulate(_err)
	}

	return token, eac.compileOR()
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

	//if err != nil {
	//	log.Fatalln(err)
	//} else {
	//	//d, _ := json.Marshal(final)
	//	fmt.Println("fn :", final)
	//}

	fmt.Println(makeDeclaration())

	elapsed := time.Since(start)
	log.Printf("Reading, Parsing / Lexing took %s", elapsed)
}
