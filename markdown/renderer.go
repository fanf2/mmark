// The package markdown outputs normalized mmark markdown. It useful to have as a mmarkfmt.
package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/kr/text"
	"github.com/mmarkdown/mmark/mast"
)

// Flags control optional behavior of Markdown renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone Flags = 0

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of Markdown renderer.
type RendererOptions struct {
	Flags Flags // Flags allow customizing this renderer's behavior

	TextWidth int

	// if set, called at the start of RenderNode(). Allows replacing rendering of some nodes
	RenderNodeHook html.RenderNodeFunc
}

// Renderer implements Renderer interface for Markdown output.
type Renderer struct {
	opts RendererOptions

	paraStart int
	indent    int
	listLevel int
}

// NewRenderer creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	if opts.TextWidth == 0 {
		opts.TextWidth = 100
	}
	return &Renderer{opts: opts, paraStart: -1}
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
}

func (r *Renderer) matter(w io.Writer, node *ast.DocumentMatter, entering bool) {
	if !entering {
		return
	}
	switch node.Matter {
	case ast.DocumentMatterFront:
		r.outs(w, "{frontmatter}\n")
	case ast.DocumentMatterMain:
		r.outs(w, "{mainmatter}\n")
	case ast.DocumentMatterBack:
		r.outs(w, "{backmatter}\n")
	}
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if !entering {
		r.cr(w)
		r.cr(w)
		return
	}
	hashes := strings.Repeat("#", node.Level)
	r.outs(w, hashes)
	r.outs(w, " ")
}

func (r *Renderer) horizontalRule(w io.Writer, node *ast.HorizontalRule) {
	r.cr(w)
	r.outs(w, "******")
	r.cr(w)
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
}

func (r *Renderer) paragraph(w io.Writer, para *ast.Paragraph, entering bool) {
	if entering {
		buf, ok := w.(*bytes.Buffer)
		if ok {
			r.paraStart = buf.Len()
		}
		return
	}

	if r.paraStart == -1 {
		return
	}

	buf, ok := w.(*bytes.Buffer)
	end := 0
	if ok {
		end = buf.Len()
	}

	// Reformat the entire buffer and rewrite to the writer.
	b := buf.Bytes()[r.paraStart:end]
	wrapped := text.WrapBytes(b, r.opts.TextWidth-r.indent)
	prefix := bytes.Repeat([]byte(" "), r.indent)
	indented := text.IndentBytes(wrapped, prefix)

	buf.Truncate(r.paraStart)

	// If in list, start at the 3rd item to print.
	_, inList := para.Parent.(*ast.ListItem)
	if inList {
		r.out(w, indented[r.indent:])
	} else {
		r.out(w, indented)
	}

	if ast.GetNextNode(para) != nil {
		r.cr(w)
		r.cr(w)
	}
}

func (r *Renderer) listEnter(w io.Writer, nodeData *ast.List) {
	r.indent += 3
}

func (r *Renderer) listExit(w io.Writer, list *ast.List) {
	r.indent -= 3
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		r.listEnter(w, list)
	} else {
		r.listExit(w, list)
	}
}

func (r *Renderer) listItemEnter(w io.Writer, listItem *ast.ListItem) {
	indent := r.indent - 3
	if indent < 0 {
		indent = 0
	}
	prefix := bytes.Repeat([]byte(" "), indent)

	switch x := listItem.ListFlags; {
	case x&ast.ListTypeOrdered != 0:
		r.out(w, prefix)
		r.outs(w, "1. ")
	default:
		r.out(w, prefix)
		r.outs(w, "*  ")
	}
}

func (r *Renderer) listItemExit(w io.Writer, listItem *ast.ListItem) {
	r.cr(w)
	r.cr(w)
}

func (r *Renderer) listItem(w io.Writer, listItem *ast.ListItem, entering bool) {
	if entering {
		r.listItemEnter(w, listItem)
	} else {
		r.listItemExit(w, listItem)
	}
}

func (r *Renderer) codeBlock(w io.Writer, codeBlock *ast.CodeBlock) {
}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {
}

func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody, entering bool) {
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {
}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
}

func (r *Renderer) index(w io.Writer, index *ast.Index, entering bool) {
	if !entering {
		return
	}

	r.outs(w, "(!")
	if index.Primary {
		r.outs(w, "!")
	}
	r.out(w, index.Item)
	if len(index.Subitem) > 0 {
		r.outs(w, ", ")
		r.out(w, index.Subitem)
	}
	r.outs(w, ")")
}

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {
	if entering {
		r.imageEnter(w, node)
	} else {
		r.imageExit(w, node)
	}
}

func (r *Renderer) imageEnter(w io.Writer, image *ast.Image) {
}

func (r *Renderer) imageExit(w io.Writer, image *ast.Image) {
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock) {
}

func (r *Renderer) captionFigure(w io.Writer, captionFigure *ast.CaptionFigure, entering bool) {
}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
}

// RenderNode renders a markdown node to markdown.
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if r.opts.RenderNodeHook != nil {
		status, didHandle := r.opts.RenderNodeHook(w, node, entering)
		if didHandle {
			return status
		}
	}

	switch node := node.(type) {
	case *ast.Document:
		// do nothing
	case *mast.Title:
	case *mast.Bibliography:
	case *mast.BibliographyItem:
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
	case *ast.Text:
		r.text(w, node, entering)

	case *ast.Softbreak:
	case *ast.Hardbreak:
	case *ast.Callout:
		r.outOneOf(w, entering, "<<", ">>")

	case *ast.Emph:
		r.outOneOf(w, entering, "*", "*")

	case *ast.Strong:
		r.outOneOf(w, entering, "**", "**")

	case *ast.Del:
	case *ast.Citation:
	case *ast.DocumentMatter:
		r.matter(w, node, entering)
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.HorizontalRule:
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
	case *ast.HTMLBlock:
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.CodeBlock:
	case *ast.Caption:
	case *ast.CaptionFigure:
	case *ast.Table:
	case *ast.TableCell:
	case *ast.TableHeader:
	case *ast.TableBody:
	case *ast.TableRow:
	case *ast.TableFooter:
	case *ast.BlockQuote:
	case *ast.Aside:
	case *ast.CrossReference:
	case *ast.Index:
		r.index(w, node, entering)
	case *ast.Link:
	case *ast.Math:
	case *ast.Image:
	case *ast.Code:
		r.outs(w, "`")
		r.out(w, node.Literal)
		r.outs(w, "`")

	case *ast.MathBlock:
	case *ast.Subscript:
	case *ast.Superscript:
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (r *Renderer) text(w io.Writer, node *ast.Text, entering bool) {
	if !entering {
		return
	}
	r.out(w, node.Literal)
}

func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
}

func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {
}

func (r *Renderer) writeDocumentHeader(w io.Writer) {
}
