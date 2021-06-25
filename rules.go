/********************************/
/*  This file contains all of   */
/*  the rule declarations       */
/*  for the final ast           */
/********************************/

package main

func makeDeclaration() (VariableDeclaration, error) {
	eac := errorAcc()

	// const, let, var
	tk, e := or_types([]int{const_var, let_var, var_var})
	eac.accumulate(e)

	// const identifier =
	identifier, e := compareNext(identifier)
	eac.accumulate(e)

	// =
	compareNext(equals)

	// const identifier = value
	value, e := compareNext(number)
	eac.accumulate(e)

	return var_declaration(identifier.Value, tk.Value, value.Value), eac.compileFirst()
}
