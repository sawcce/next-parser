# Next Parser

⚠️This project is in very early development⚠️
⚠️So this readme will for now explain how I'm tackling this pretty tideous task⚠️

## Idea

What was the idea behind all this ? 🤔

I wanted to write my own parser for a long time and since well, I love go and JS I thought : Why not make a JS parser in Go ?

I'll try to make the code and phylosphy as reusable as possible and not super scoped like I may have done in the past when trying to build parsers 🤓

This will allow better performance as well as more readability in general.

## Logic

The parser will be divided into two categories : The actual ast-builder and the lexer

The lexer transforms a string into a list of tokens and the ast-builder will take these tokens and transform them into a list of instructions

Right now how I've been working on it, I plan on making the project "rule based" 📕

#### What does that even mean ?

In my logic "rule based" means that the parser will be fed "rules" that are then checked on a list of tokens. 

These tokens will be matched onto different patterns to build an ast from the tokens fed into it.