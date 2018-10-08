package markdown

import (
	"bytes"
	"io"

	"github.com/kr/text"
)

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }

func (r *Renderer) cr(w io.Writer) {
	// suppress multiple newlines
	if buf, ok := w.(*bytes.Buffer); ok {
		b := buf.Bytes()
		if len(b) > 2 && b[len(b)-1] == '\n' && b[len(b)-2] == '\n' {
			return
		}
	}

	r.outs(w, "\n")
}

// wrapText wraps the text in data, taking r.indent into account.
func (r *Renderer) wrapText(data []byte) []byte {
	wrapped := text.WrapBytes(data, r.opts.TextWidth-r.indent)
	prefix := bytes.Repeat([]byte(" "), r.indent)
	indented := text.IndentBytes(wrapped, prefix)
	return indented
}

func (r *Renderer) indentText(data []byte) []byte {
	prefix := bytes.Repeat([]byte(" "), r.indent)
	indented := text.IndentBytes(data, prefix)
	return indented
}

func (r *Renderer) escapeText(data []byte) []byte {

	// TODO, escape < \ and _, but only if there are not escaped yet.
	return nil
}
