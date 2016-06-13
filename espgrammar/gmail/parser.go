package gmail

import (
	"strings"
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) (stateFn, error)

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// lexer holds the state of the scanner.
type lexer struct {
	input   string  // the string being scanned
	state   stateFn // the next lexing function to enter
	pos     Pos     // current position in the input
	start   Pos     // start position of this item
	width   Pos     // width of last rune read from input
	lastPos Pos     // position of most recent item returned by nextItem
}

const eof = -1

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// lex creates a new scanner for the input string.
func lex(mail string) bool {
	l := &lexer{
		input: mail,
	}
	err := l.run()
	return err != nil
}

// run runs the state machine for the lexer.
func (l *lexer) run() error {
	var err error
	for l.state = start; l.state != nil; {
		l.state, err = l.state(l)
		if err != nil {
			break
		}
	}
	return err
}

func start(l *lexer) (stateFn, error) {
	return nil, nil
}
