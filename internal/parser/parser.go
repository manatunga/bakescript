package parser

import "github.com/manatunga/eloque/internal/lexer"

// Parser maintains the state required to construct an Eloque
// abstract syntax tree from a sequence of lexical tokens.
// Tokens remain in source order, including structural tokens
// such as NEWLINE, INDENT, DEDENT, and the terminal EOF.
// Current identifies the next token to be consumed.
type Parser struct {
	tokens  []lexer.Token
	current int
}

// New creates a parser for an already-tokenized Eloque source.
// Tokens must be ordered, nonempty, and terminated by TokenEOF.
// TokenEOF must be the final token in the slice.
// The parser does not perform lexical analysis or modify the tokens.
// Passing a token stream that violates these preconditions is unsupported.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// peek returns the current token without consuming it.
// When the parser reaches the end of its token stream,
// it continues returning the terminal EOF token.
func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

// advance consumes and returns the current token.
// The cursor advances unless the current token is EOF,
// which remains available for subsequent inspection.
func (p *Parser) advance() lexer.Token {
	currentToken := p.peek()
	if currentToken.Kind != lexer.TokenEOF {
		p.current++
	}
	return currentToken
}

// check reports whether the current token has the specified kind.
// It does not consume or modify the token stream.
func (p *Parser) check(kind lexer.TokenKind) bool {
	return p.peek().Kind == kind
}

// match consumes the current token if its kind matches any of
// the supplied kinds. It returns true when a match is found
// and false otherwise. A failed match leaves the cursor unchanged.
func (p *Parser) match(kinds ...lexer.TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes and returns the current token when its kind matches
// the required kind. Otherwise, it returns a positioned syntax
// diagnostic describing the mismatch and leaves the cursor unchanged.
// The supplied message provides context for the expected syntax.
func (p *Parser) expect(
	kind lexer.TokenKind,
	message string,
) (lexer.Token, *Diagnostic) {
	currentToken := p.peek()
	if currentToken.Kind == kind {
		return p.advance(), nil
	}
	return lexer.Token{}, &Diagnostic{
		Message:  message,
		Location: currentToken.Span,
		Expected: []lexer.TokenKind{kind},
		Actual:   &currentToken.Kind,
	}
}
