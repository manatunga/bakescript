package parser

import (
	"strings"

	"github.com/manatunga/eloque/internal/ast"
	"github.com/manatunga/eloque/internal/lexer"
	"github.com/manatunga/eloque/internal/source"
)

// expressionResult represents the result of parsing an Eloque expression.
//
// The expression field contains the semantic AST node, while extent
// records the complete source range consumed by the expression,
// including grouping parentheses that are not represented in the AST.
//
// Keeping these separately allows the parser to construct accurate
// enclosing expression spans without modifying child node locations.
type expressionResult struct {
	expression ast.Expression
	extent     source.Span
}

// Expression binding powers define the precedence of Eloque operators.
// Higher values indicate tighter binding. These constants follow the
// precedence rules established in the language specification and are
// used by the Pratt parser to determine expression grouping.
const (
	bindingPowerAdditive       = 50
	bindingPowerMultiplicative = 60
	bindingPowerUnary          = 70
	bindingPowerExponentiation = 80
)

// binaryBindingPower returns the binding power of a supported
// binary arithmetic operator.
//
// The boolean result indicates whether the token kind represents
// an operator supported by the current expression parser.
//
// Operators with higher binding powers bind more tightly.
// Addition, subtraction, multiplication, division, floor division,
// and modulo are left-associative. Exponentiation is right-associative.
func binaryBindingPower(kind lexer.TokenKind) (int, bool) {
	switch kind {
	case lexer.TokenPlus, lexer.TokenMinus:
		return bindingPowerAdditive, true
	case lexer.TokenStar, lexer.TokenSlash, lexer.TokenFloorDiv, lexer.TokenPercent:
		return bindingPowerMultiplicative, true
	case lexer.TokenPower:
		return bindingPowerExponentiation, true
	default:
		return 0, false
	}
}

// parsePrefix parses a prefix unary expression or a primary expression.
//
// Unary plus and minus are right-nesting prefix operators.
// Their operands are parsed using the unary binding power, allowing
// higher-precedence operators to bind inside the operand.
//
// The returned expression result preserves both the semantic AST
// and the full syntactic source extent.
func (p *Parser) parsePrefix() (expressionResult, *Diagnostic) {
	operator := p.peek()
	if operator.Kind != lexer.TokenPlus && operator.Kind != lexer.TokenMinus {
		return p.parsePrimary()
	}
	p.advance()
	operand, diagnostic := p.parseExpression(bindingPowerUnary)
	if diagnostic != nil {
		return expressionResult{}, diagnostic
	}
	unaryExpression := &ast.UnaryExpression{
		Operator: operator.Kind,
		Operand:  operand.expression,
		Location: source.Span{
			Start: operator.Span.Start,
			End:   operand.extent.End,
		},
	}
	return expressionResult{
		expression: unaryExpression,
		extent: source.Span{
			Start: operator.Span.Start,
			End:   operand.extent.End,
		},
	}, nil
}

// decodeStringLiteral converts a validated Eloque string token's
// lexeme into its decoded string value.
//
// The input includes the original surrounding quotation marks.
// Escape sequences are converted to their corresponding characters.
// The lexer is responsible for rejecting malformed string literals
// before this function is called.
func decodeStringLiteral(lexeme string) string {
	innerString := lexeme[1 : len(lexeme)-1]
	var result strings.Builder
	result.Grow(len(innerString))

	for i := 0; i < len(innerString); i++ {
		char := innerString[i]
		if char == '\\' && i+1 < len(innerString) {
			nextChar := innerString[i+1]
			switch nextChar {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '\'':
				result.WriteByte('\'')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte('\\')
				result.WriteByte(nextChar)
			}
			i++
		} else {
			result.WriteByte(char)
		}
	}
	return result.String()
}

// parsePrimary parses a literal, identifier, or parenthesized
// Eloque expression.
//
// Parenthesized expressions reuse the inner semantic AST node
// rather than introducing a grouping node. The returned extent
// includes the surrounding parentheses, while the inner node's
// source span remains unchanged.
//
// Unexpected tokens produce a positioned syntax diagnostic
// without consuming the offending token.
func (p *Parser) parsePrimary() (expressionResult, *Diagnostic) {
	currentToken := p.peek()
	switch currentToken.Kind {
	case lexer.TokenInt:
		p.advance()
		expression := &ast.IntegerLiteral{
			Value:    currentToken.Lexeme,
			Location: currentToken.Span,
		}
		return expressionResult{
			expression: expression,
			extent:     currentToken.Span,
		}, nil
	case lexer.TokenString:
		p.advance()
		expression := &ast.StringLiteral{
			Value:    decodeStringLiteral(currentToken.Lexeme),
			Location: currentToken.Span,
		}
		return expressionResult{
			expression: expression,
			extent:     currentToken.Span,
		}, nil
	case lexer.TokenIdent:
		p.advance()
		expression := &ast.Identifier{
			Name:     currentToken.Lexeme,
			Location: currentToken.Span,
		}
		return expressionResult{
			expression: expression,
			extent:     currentToken.Span,
		}, nil
	case lexer.TokenLeftParen:
		openingToken := p.advance()
		innerResult, diagnostic := p.parseExpression(0)
		if diagnostic != nil {
			return expressionResult{}, diagnostic
		}
		closingToken, diagnostic := p.expect(
			lexer.TokenRightParen,
			"expected ')' to close parenthesized expression",
		)
		if diagnostic != nil {
			return expressionResult{}, diagnostic
		}
		return expressionResult{
			expression: innerResult.expression,
			extent: source.Span{
				Start: openingToken.Span.Start,
				End:   closingToken.Span.End,
			},
		}, nil
	default:
		return expressionResult{}, &Diagnostic{
			Message:  "expected an expression",
			Location: currentToken.Span,
			Expected: []lexer.TokenKind{
				lexer.TokenInt,
				lexer.TokenString,
				lexer.TokenIdent,
				lexer.TokenLeftParen,
				lexer.TokenPlus,
				lexer.TokenMinus,
			},
			Actual: &currentToken.Kind,
		}
	}
}

// parseExpression parses an Eloque expression using Pratt parsing.
//
// minBindingPower determines the minimum binding power required
// for a binary operator to be consumed at the current parsing level.
//
// The parser first reads a prefix unary or primary expression,
// then repeatedly consumes eligible binary arithmetic operators
// and recursively parses their right operands.
//
// Addition, subtraction, multiplication, division, floor division,
// and modulo are left-associative. Exponentiation is right-associative.
//
// Exponentiation binds more tightly than prefix unary operators,
// allowing expressions such as -2 ** 2 and 2 ** -3 to be parsed
// according to Eloque's precedence rules.
//
// The returned result preserves the semantic AST node separately
// from the complete syntactic source extent.
//
// Syntax errors produce positioned diagnostics.
func (p *Parser) parseExpression(minBindingPower int) (expressionResult, *Diagnostic) {
	left, diagnostic := p.parsePrefix()
	if diagnostic != nil {
		return expressionResult{}, diagnostic
	}
	for {
		bindingPower, isBinaryOperator := binaryBindingPower(p.peek().Kind)
		if !isBinaryOperator || bindingPower < minBindingPower {
			break
		}
		operatorToken := p.advance()
		rightMinBindingPower := bindingPower + 1
		if operatorToken.Kind == lexer.TokenPower {
			rightMinBindingPower = bindingPower
		}
		right, diagnostic := p.parseExpression(rightMinBindingPower)
		if diagnostic != nil {
			return expressionResult{}, diagnostic
		}
		binaryExpression := &ast.BinaryExpression{
			Left:     left.expression,
			Operator: operatorToken.Kind,
			Right:    right.expression,
			Location: source.Span{
				Start: left.extent.Start,
				End:   right.extent.End,
			},
		}
		left = expressionResult{
			expression: binaryExpression,
			extent: source.Span{
				Start: left.extent.Start,
				End:   right.extent.End,
			},
		}
	}
	return left, nil
}
