package lexer

// TokenKind identifies the lexical category of a token.
// Each category represents a keyword, literal, identifier, or structural token.
type TokenKind uint8

// Token kinds currently supported by Eloque's lexer.
// Each kind has a compact integer representation.
// Additional kinds will be introduced as the language's syntax is implemented.
const (
	TokenInvalid TokenKind = iota
	TokenLet
	TokenBe
	TokenSay
	TokenIdent
	TokenInt
	TokenString
	TokenNewline
	TokenEOF
	TokenIndent
	TokenDedent
	TokenColon
	TokenLeftParen
	TokenRightParen
	TokenPlus
	TokenMinus
	TokenStar
	TokenPower
	TokenSlash
	TokenFloorDiv
	TokenPercent
)

// String returns the human-readable name of a token kind.
// Unknown values are represented as UNKNOWN.
func (kind TokenKind) String() string {
	switch kind {
	case TokenInvalid:
		return "INVALID"
	case TokenLet:
		return "LET"
	case TokenBe:
		return "BE"
	case TokenSay:
		return "SAY"
	case TokenIdent:
		return "IDENT"
	case TokenInt:
		return "INT"
	case TokenString:
		return "STRING"
	case TokenNewline:
		return "NEWLINE"
	case TokenEOF:
		return "EOF"
	case TokenIndent:
		return "INDENT"
	case TokenDedent:
		return "DEDENT"
	case TokenColon:
		return "COLON"
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenStar:
		return "STAR"
	case TokenPower:
		return "POWER"
	case TokenSlash:
		return "SLASH"
	case TokenFloorDiv:
		return "FLOOR_DIV"
	case TokenPercent:
		return "PERCENT"
	default:
		return "UNKNOWN"
	}
}

// Position identifies a location in Eloque source code.
// Line and column numbers are 1-based.
type Position struct {
	Line   int
	Column int
}

// Span identifies a half-open range in Eloque source code.
// Start is inclusive and End is exclusive.
type Span struct {
	Start Position
	End   Position
}

// Token represents a recognized unit of Eloque source code.
// Kind identifies its lexical category, and Lexeme preserves the
// original source text. Span records its half-open source range.
type Token struct {
	Kind   TokenKind
	Lexeme string
	Span   Span
}

// lookupIdentifier determines whether a scanned word is a reserved keyword.
// Words that do not exactly match a recognized keyword are identifiers.
func lookupIdentifier(word string) TokenKind {
	switch word {
	case "let":
		return TokenLet
	case "be":
		return TokenBe
	case "say":
		return TokenSay
	default:
		return TokenIdent
	}
}
