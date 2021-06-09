package calc

import (
	"fmt"
	"strconv"
	"strings"
)

// These constants are used to return the type of
// token which has been lexed.
const (
	// Basic token-types
	EOF    = "EOF"
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	ERROR  = "ERROR"

	// Assignment-magic
	LET    = "LET"
	ASSIGN = "="

	// Paren
	LPAREN = "("
	RPAREN = ")"

	// Operations
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
)

// Token holds a lexed token from our input.
type Token struct {

	// The type of the token.
	Type string

	// The value of the token.
	//
	// If the type of the token is NUMBER then this
	// will be stored as a float64.  Otherwise the
	// value will be a string representation of the token.
	//
	Value interface{}
}

// Lexer holds our lexer state.
type Lexer struct {

	// input is the string we're lexing.
	input string

	// position is the current position within the input-string.
	position int

	// simple map of single-character tokens to their type
	known map[string]string
}

// NewLexer creates a new lexer, for the given input.
func NewLexer(input string) *Lexer {

	// Create the lexer object.
	l := &Lexer{input: input}

	// Populate the simple token-types in a map for
	// later use.
	l.known = make(map[string]string)
	l.known["*"] = MULTIPLY
	l.known["+"] = PLUS
	l.known["-"] = MINUS
	l.known["/"] = DIVIDE
	l.known["="] = ASSIGN
	l.known["("] = LPAREN
	l.known[")"] = RPAREN

	return l
}

// Next returns the next token from our input stream.
//
// This is pretty naive lexer, however it is sufficient to
// recognize numbers, identifiers, and our small set of
// operators.
func (l *Lexer) Next() *Token {

	// Loop until we've exhausted our input.
	for l.position < len(l.input) {

		// Get the next character
		char := string(l.input[l.position])

		// Is this a known character/token?
		t, ok := l.known[char]
		if ok {
			// skip the character, and return the token
			l.position++
			return &Token{Value: char, Type: t}
		}

		// If we reach here it is something more complex.
		switch char {

		// Skip whitespace
		case " ", "\n", "\r", "\t", ";":
			l.position++
			continue

			// Is it a digit?
		case "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".":

			//
			// Loop for more digits
			//

			// Starting offset of our number
			start := l.position

			// ending offset of our number.
			end := l.position

			// keep walking forward, minding we don't wander
			// out of our input.
			for end < len(l.input) {

				if !l.isNumberComponent(l.input[end], end == start) {
					break
				}
				end++
			}

			l.position = end

			// Here we have the number
			token := l.input[start:end]

			// too many periods?
			bits := strings.Split(token, ".")
			if len(bits) > 2 {
				return &Token{Type: ERROR, Value: fmt.Sprintf("too many periods in '%s'", token)}
			}

			// Convert to float64
			number, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return &Token{Value: fmt.Sprintf("failed to parse number: %s", err.Error()), Type: ERROR}
			}

			return &Token{Value: number, Type: NUMBER}
		}

		//
		// We'll assume we have an identifier at this point.
		//

		// Starting offset of our ident
		start := l.position

		// ending offset of our ident.
		end := l.position

		// keep walking forward, minding we don't wander
		// out of our input.
		for end < len(l.input) {

			// Build up identifiers from any permitted
			// character - which is just a-zA-Z
			if l.isIdentifierCharacter(l.input[end]) {
				end++
			} else {
				break
			}
		}

		// Change the position to be after the end of the identifier
		// we found - if we didn't find one then that results in no
		// change.
		l.position = end

		// Now record the text of the token (i.e. identifier).
		token := l.input[start:end]

		//
		// In a real language/lexer we might have
		// keywords/reserved-words to handle.
		//
		// We only need to cope with "let".
		//
		// If the identifier was LET then return that
		// token instead.
		//
		if strings.ToLower(token) == "let" {
			return &Token{Value: "let", Type: LET}
		}

		//
		// So we handled the easy cases, and then defaulted
		// to looking for an identifier.
		//
		// If we failed to find one that means that we've got
		// to skip the unknown character - to avoid an infinite
		// the next time we try to move forwards - and report
		// the error immediately.
		//
		if token == "" {
			l.position++
			return &Token{Value: fmt.Sprintf("unknown character %c", l.input[end]), Type: ERROR}
		}

		//
		// We found a non-empty identifier, which
		// wasn't converted into a  `let` keyword.
		//
		// Return it.
		//
		return &Token{Value: token, Type: IDENT}

	}

	//
	// If we get here then we've walked past the end of
	// our input-string.
	//
	return &Token{Value: "", Type: EOF}
}

// isIdentifierCharacter tests whether the given character is
// valid for use in an identifier.
func (l *Lexer) isIdentifierCharacter(d byte) bool {

	if (d >= 'a' && d <= 'z') ||
		(d >= 'A' && d <= 'Z') {
		return true
	}

	return false
}

// isNumberComponent looks for characters that can make up integers/floats
//
// We handle the first-character specially, which is why that's an argument
func (l *Lexer) isNumberComponent(d byte, first bool) bool {

	// digits
	if d >= '0' && d <= '9' {
		return true
	}

	// floating-point numbers
	if d == '.' {
		return true
	}

	// negative sign can only occur at the start of the input
	if d == '-' && first {
		return true
	}
	return false
}
