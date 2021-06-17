package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type jsonScanner struct {
	buf *bufio.Reader
	w   writer
	err error
	// used for parser data
	data interface{}
}

type writer interface {
	io.ByteWriter
	io.Writer
	io.StringWriter
}

// FixJSONSpacing modifies a valid JSON input to include
// spaces between elements for compatibility with
// Python's json.dumps method
func FixJSONSpacing(r io.Reader) ([]byte, error) {
	w := new(bytes.Buffer)
	s := newJSONScanner(r, w)
	yyParse(s)
	if s.err != nil {
		return nil, s.err
	}
	return w.Bytes(), nil
}

func newJSONScanner(r io.Reader, w writer) *jsonScanner {
	return &jsonScanner{
		buf: bufio.NewReader(r),
		w:   w,
	}
}

func (sc *jsonScanner) Error(s string) {
	sc.err = fmt.Errorf("syntax error: %s", s)
}

func (sc *jsonScanner) Reduced(rule, state int, lval *yySymType) bool {
	return false
}

func (s *jsonScanner) Lex(lval *yySymType) int {
	return s.lex(lval)
}

func (s *jsonScanner) lex(lval *yySymType) int {
	for {
		r := s.read()
		if r == 0 {
			return 0
		}
		if isWhitespace(r) {
			continue
		}

		if isDigit(r) {
			s.unread()
			lval.i = s.scanNumber()
			return NUMBER
		}

		switch r {
		case '[':
			s.w.WriteByte('[')
			return LS
		case ']':
			s.w.WriteByte(']')
			return RS
		case '{':
			s.w.WriteByte('{')
			return LC
		case '}':
			s.w.WriteByte('}')
			return RC
		case ',':
			s.w.WriteString(", ")
			return COMMA
		case ':':
			s.w.WriteString(": ")
			return COLON
		case '"':
			s.unread()
			lval.s = s.scanStr()
			return STRING
		case 't':
			s.unread()
			if s.scanTrue() {
				return TRUE
			}
		case 'f':
			s.unread()
			if s.scanFalse() {
				return FALSE
			}
		case 'n':
			s.unread()
			if s.scanNull() {
				return NULL
			}
		default:
			s.err = errors.New("error")
			return 0
		}
	}
}

func (s *jsonScanner) scanTrue() bool {
	t := []rune{'t', 'r', 'u', 'e'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("true is error")
			return false
		}
	}
	s.w.WriteString("true")
	return true
}

func (s *jsonScanner) scanFalse() bool {
	t := []rune{'f', 'a', 'l', 's', 'e'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("false is error")
			return false
		}
	}
	s.w.WriteString("false")
	return true
}

func (s *jsonScanner) scanNull() bool {
	t := []rune{'n', 'u', 'l', 'l'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("null is error")
			return false
		}
	}
	s.w.WriteString("null")
	return true
}

func (s *jsonScanner) scanStr() string {
	var str []rune
	//begin with ", end with "
	r := s.read()
	if r != '"' {
		os.Exit(1)
	}

	for {
		r := s.read()
		if r == '"' || r == 1 {
			break
		}
		str = append(str, r)
	}
	s.w.WriteString(`"` + string(str) + `"`)
	return string(str)
}

func (s *jsonScanner) scanNumber() interface{} {
	var number []rune
	var isFloat bool
	for {
		r := s.read()
		if r == '.' && len(number) > 0 && !isFloat {
			isFloat = true
			number = append(number, r)
			continue
		}

		if isWhitespace(r) || r == ',' || r == '}' || r == ']' {
			s.unread()
			break
		}
		if !isDigit(r) {
			return nil
		}
		number = append(number, r)
	}
	if isFloat {
		f, _ := strconv.ParseFloat(string(number), 64)
		return f
	}
	i, _ := strconv.Atoi(string(number))
	s.w.WriteString(string(number))
	return i
}

func (s *jsonScanner) read() rune {
	ch, _, _ := s.buf.ReadRune()
	return ch
}

func (s *jsonScanner) unread() { _ = s.buf.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
