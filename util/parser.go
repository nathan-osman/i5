package util

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

var (
	ErrEOF              = errors.New("end of file")
	ErrInvalidCharacter = errors.New("invalid character")
	ErrEqualsExpected   = errors.New("'=' expected")
	ErrQuoteExpected    = errors.New("'\"' expected")
)

// Parser provides a parser for data in the "k1=v1 k2=v2 ..." format.
type Parser struct {
	reader *bufio.Reader
	r      rune
}

func (p *Parser) isAl() bool {
	return p.r >= 'A' && p.r <= 'Z' || p.r >= 'a' && p.r <= 'z' || p.r == '_'
}

func (p *Parser) isAlNum() bool {
	return p.isAl() || p.r >= '0' && p.r <= '9'
}

func (p *Parser) isEOF() bool {
	return p.r == 0
}

func (p *Parser) readNextChar() error {
	r, _, err := p.reader.ReadRune()
	if err != nil && err != io.EOF {
		return err
	}
	p.r = r
	return nil
}

func (p *Parser) skipSpaces() error {
	for unicode.IsSpace(p.r) {
		if err := p.readNextChar(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseKey() (string, error) {
	if !p.isAl() {
		return "", ErrInvalidCharacter
	}
	var key string
	for p.isAlNum() {
		key += string(p.r)
		if err := p.readNextChar(); err != nil {
			return "", err
		}
	}
	return key, nil
}

func (p *Parser) parseValue() (string, error) {
	var isQuoted bool
	if p.r == '"' {
		isQuoted = true
		if err := p.readNextChar(); err != nil {
			return "", err
		}
	}
	var val string
	for !p.isEOF() {
		if isQuoted {
			if p.r == '\\' {
				if err := p.readNextChar(); err != nil {
					return "", err
				}
				if p.isEOF() {
					break
				}
			} else if p.r == '"' {
				break
			}
		} else if unicode.IsSpace(p.r) {
			break
		}
		val += string(p.r)
		if err := p.readNextChar(); err != nil {
			return "", err
		}
	}
	if isQuoted {
		if p.r != '"' {
			return "", ErrQuoteExpected
		}
		if err := p.readNextChar(); err != nil {
			return "", err
		}
	}
	return val, nil
}

// NewParser creates a new parser for the provided reader.
func NewParser(r io.Reader) (*Parser, error) {
	p := &Parser{
		reader: bufio.NewReader(r),
	}
	if err := p.readNextChar(); err != nil {
		return nil, err
	}
	return p, nil
}

// ParseNextEntry parses the next key=value entry.
func (p *Parser) ParseNextEntry() (string, string, error) {
	if err := p.skipSpaces(); err != nil {
		return "", "", err
	}
	if p.isEOF() {
		return "", "", ErrEOF
	}
	k, err := p.parseKey()
	if err != nil {
		return "", "", err
	}
	if err := p.skipSpaces(); err != nil {
		return "", "", err
	}
	if p.r != '=' {
		return "", "", ErrEqualsExpected
	}
	if err := p.readNextChar(); err != nil {
		return "", "", err
	}
	if err := p.skipSpaces(); err != nil {
		return "", "", err
	}
	v, err := p.parseValue()
	if err != nil {
		return "", "", err
	}
	return k, v, nil
}
