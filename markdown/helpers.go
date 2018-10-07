package markdown

import (
	"bytes"
	"io"
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
