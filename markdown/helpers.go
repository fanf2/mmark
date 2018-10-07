package markdown

import (
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
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }
