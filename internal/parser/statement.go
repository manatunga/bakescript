package parser

import (
	"github.com/manatunga/eloque/internal/ast"
	"github.com/manatunga/eloque/internal/lexer"
	"github.com/manatunga/eloque/internal/source"
)

// parseLetStatement parses a mutable variable declaration in Eloque.
// It expects the "let" keyword, a binding identifier, the "be"
// keyword, and an initializer expression.
//
// The resulting LetStatement spans the complete declaration,
// from the opening keyword through the initializer's syntactic extent.
// The statement terminator (NEWLINE or EOF) is left unconsumed
// for the enclosing program parser to handle.
//
// Malformed declarations produce a positioned syntax diagnostic.
// The current token must be TokenLet when this function is called.
func (p *Parser) parseLetStatement() (ast.Statement, *Diagnostic) {
	letToken := p.advance()
	identifier, diagnostic := p.expect(
		lexer.TokenIdent,
		"expected a variable name after 'let'",
	)
	if diagnostic != nil {
		return nil, diagnostic
	}
	_, diagnostic = p.expect(
		lexer.TokenBe,
		"expected 'be' after variable name to initialize a variable value",
	)
	if diagnostic != nil {
		return nil, diagnostic
	}
	initializer, diagnostic := p.parseExpression(0)
	if diagnostic != nil {
		return nil, diagnostic
	}
	letStatement := &ast.LetStatement{
		Name:  identifier.Lexeme,
		Value: initializer.expression,
		Location: source.Span{
			Start: letToken.Span.Start,
			End:   initializer.extent.End,
		},
	}
	return letStatement, nil
}

// parseChangeStatement parses a variable reassignment in Eloque.
// The current token must be TokenChange when this function is called.
// It expects the "change" keyword, a binding identifier, the "to"
// keyword, and an assignment expression.
//
// The resulting ChangeStatement spans the complete assignment,
// from the opening keyword through the expression's syntactic extent.
// The statement terminator (NEWLINE or EOF) is left unconsumed
// for the enclosing program parser to handle.
//
// Malformed assignments produce a positioned syntax diagnostic.
func (p *Parser) parseChangeStatement() (ast.Statement, *Diagnostic) {
	changeToken := p.advance()
	identifier, diagnostic := p.expect(
		lexer.TokenIdent,
		"expected a variable name after 'change'",
	)
	if diagnostic != nil {
		return nil, diagnostic
	}
	_, diagnostic = p.expect(
		lexer.TokenTo,
		"expected 'to' after variable name to assign a new variable value",
	)
	if diagnostic != nil {
		return nil, diagnostic
	}
	assigned, diagnostic := p.parseExpression(0)
	if diagnostic != nil {
		return nil, diagnostic
	}
	changeStatement := &ast.ChangeStatement{
		Name:  identifier.Lexeme,
		Value: assigned.expression,
		Location: source.Span{
			Start: changeToken.Span.Start,
			End:   assigned.extent.End,
		},
	}
	return changeStatement, nil
}

// parseSayStatement parses an output statement in Eloque.
// The current token must be TokenSay when this function is called.
// It expects the "say" keyword followed by an expression.
//
// The resulting SayStatement spans the complete statement,
// from the opening keyword through the expression's syntactic extent.
// The statement terminator (NEWLINE or EOF) is left unconsumed
// for the enclosing program parser to handle.
//
// Malformed output statements produce a positioned syntax diagnostic.
func (p *Parser) parseSayStatement() (ast.Statement, *Diagnostic) {
	sayToken := p.advance()
	value, diagnostic := p.parseExpression(0)
	if diagnostic != nil {
		return nil, diagnostic
	}
	sayStatement := &ast.SayStatement{
		Value: value.expression,
		Location: source.Span{
			Start: sayToken.Span.Start,
			End:   value.extent.End,
		},
	}
	return sayStatement, nil
}

// parseStatement parses a single supported Eloque statement.
// It dispatches based on the current token without consuming
// statement terminators.
//
// Unsupported statement starters produce a positioned syntax diagnostic.
func (p *Parser) parseStatement() (ast.Statement, *Diagnostic) {
	currentToken := p.peek()
	switch currentToken.Kind {
	case lexer.TokenLet:
		return p.parseLetStatement()
	case lexer.TokenChange:
		return p.parseChangeStatement()
	case lexer.TokenSay:
		return p.parseSayStatement()
	default:
		return nil, &Diagnostic{
			Message:  "expected a statement ('let', 'change' or 'say')",
			Location: currentToken.Span,
			Expected: []lexer.TokenKind{
				lexer.TokenLet,
				lexer.TokenChange,
				lexer.TokenSay,
			},
			Actual: &currentToken.Kind,
		}
	}
}
