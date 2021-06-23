
package main
	multipleDigits := rule(multiple, []Rule{rawRule(digit)})

	decimal := []Rule{rule_cb(and, []Rule{multipleDigits, rawRule(dot), multipleDigits}, func(args ...string) interface{} {
		numStr := strings.Join(args, "")
		numRep, _ := strconv.ParseFloat(numStr, 64)
		fmt.Println("Parsed num:", numRep)
		return 123
	})}

	dig1 := token_type_val(digit, "1")
	tk_dot := token_type_val(dot, ".")

	toTest := []Token{dig1, dig1, tk_dot, dig1, dig1}
	//number := []Rule{rule(and, decimal)}

	breakPoint, matches, final := consume_ruleset_token(decimal, toTest)

	fmt.Println("[Result] Does match ?", matches, ", breakPoint :", breakPoint, "Final:", final)

func consume_ruleset_token(ruleset []Rule, tokens []Token) (int, bool, []interface{}) {
	final := []interface{}{}
	matches := false
	breakPoint := 0

	x := 0

	fmt.Println(">> Entering with the ruleset :", ruleset)
	fmt.Println(">> Entering with the tokens :", tokens)

	for i, rule := range ruleset {
		matches := false
		if len(tokens) == 0 {
			break
		}

		fmt.Println("[", i, "]", "[", x, "]")
		breakPoint = i
		compare := tokens[x]

		if rule.Raw == true {
			fmt.Println("[Raw] with the rule being :", Tokens[rule.Type], "and the compared :", Tokens[compare.Type])
			x += 1
			final = append(final, compare.Value)
			matches = rule.Type == compare.Type
		} else {
			switch rule.Type {
			case and:
				_breakPoint, match, _ := consume_ruleset_token(rule.Subrules, tokens[x:])
				if breakPoint != 0 {
					final = append(final, rule.Callback(tokens_to_strings(tokens[x:_breakPoint])))
				}
				fmt.Println("[SubRule] And: bp:", _breakPoint, match)
				breakPoint += _breakPoint
				matches = match
				break
			case or:
				break
			case multiple:
				fmt.Println("[SubRule] Multiple x:", x)
				running := true
				fbp := 0
				bp := x
				start := x
				var fn []interface{}

				for running == true {
					_breakPoint, match, final := consume_ruleset_token(rule.Subrules, tokens[bp:])
					running = match
					bp += _breakPoint
					fbp += _breakPoint
					fmt.Print("|", _breakPoint, "|")
					fn = final
				}

				fmt.Println("All the tokens:", tokens[start:bp], "Final :", fn)

				final = append(final, rule.Callback(fn))

				fmt.Println("FBP :", fbp)

				matches = (fbp != 0)

				x += fbp
				breakPoint += fbp

				fmt.Println("[Mutiple : Finished executing, current token idx:", x, ", current rule idx:", i)
				fmt.Println("[Mutiple : did it match?", matches)
				break
			}
		}

		if !matches {
			fmt.Println("[Exit not matching]", x, i)
			break
		}

	}

	if len(ruleset) == 1 && matches {
		breakPoint = 1
	}

	if breakPoint > 0 && !matches {
		matches = true
	}

	fmt.Println(">> Exit", final)
	fmt.Println(x, len(ruleset))
	return breakPoint + 1, matches, final
}