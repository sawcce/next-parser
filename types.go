/********************************/
/*  This file contains all of   */
/*  the struct declarations     */
/*  for the final ast           */
/********************************/

package main

import "strconv"

// LITERAL

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

// Variable Declaration

type VariableDeclaration struct {
	declaredVariables []string
	values            []interface{}
	kind              string
}

func var_declaration(id string, kind string, val interface{}) VariableDeclaration {
	return VariableDeclaration{[]string{id}, []interface{}{val}, kind}
}
