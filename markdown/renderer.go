// The package markdown outputs normalized mmark markdown. It useful to have as a mmarkfmt.
package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/kr/text"
	"github.com/miekg/mmark/mast"
)

// Flags control optional behavior of Markdown renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone Flags = 0

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of XML renderer.
type RendererOptions struct {
	Flags Flags // Flags allow customizing this renderer's behavior

	TextWidth int

	// if set, called at the start of RenderNode(). Allows replacing
	// rendering of some nodes
	RenderNodeHook html.RenderNodeFunc
}

// Renderer implements Renderer interface for Markdown output.
type Renderer struct {
	opts RendererOptions
}

// NewRenderer creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	if opts.TextWidth == 0 {
		opts.TextWidth = 80
	}
	return &Renderer{opts: opts}
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
}

func (r *Renderer) matter(w io.Writer, node *ast.DocumentMatter) {
}

func (r *Renderer) headingEnter(w io.Writer, heading *ast.Heading) {
}

func (r *Renderer) headingExit(w io.Writer, heading *ast.Heading) {
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
}

var rule = strings.Repeat("-", 60)

func (r *Renderer) horizontalRule(w io.Writer, node *ast.HorizontalRule) {
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
}

func (r *Renderer) paragraph(w io.Writer, para *ast.Paragraph, entering bool) {
	if !entering {
		r.cr(w)
		r.cr(w)
	}
}

func (r *Renderer) listEnter(w io.Writer, nodeData *ast.List) {
}

func (r *Renderer) listExit(w io.Writer, list *ast.List) {
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		r.listEnter(w, list)
	} else {
		r.listExit(w, list)
	}
}

func (r *Renderer) listItemEnter(w io.Writer, listItem *ast.ListItem) {
}

func (r *Renderer) listItemExit(w io.Writer, listItem *ast.ListItem) {
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

func (r *Renderer) callout(w io.Writer, callout *ast.Callout) {
}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
}

func (r *Renderer) index(w io.Writer, index *ast.Index) {
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

func (r *Renderer) code(w io.Writer, node *ast.Code) {
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock) {
}

func (r *Renderer) captionFigure(w io.Writer, captionFigure *ast.CaptionFigure, entering bool) {
}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
}

// RenderNode renders a markdown node to XML.
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
		r.Text(w, node, entering)

	case *ast.Softbreak:
	case *ast.Hardbreak:
	case *ast.Callout:
	case *ast.Emph:
		r.outOneOf(w, entering, "*", "*")

	case *ast.Strong:
		r.outOneOf(w, entering, "**", "**")

	case *ast.Del:
	case *ast.Citation:
	case *ast.DocumentMatter:
	case *ast.Heading:
	case *ast.HorizontalRule:
	case *ast.Paragraph:
		r.paragraph(w, node, entering)

	case *ast.HTMLSpan:
	case *ast.HTMLBlock:
	case *ast.List:
	case *ast.ListItem:
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
	case *ast.Link:
	case *ast.Math:
	case *ast.Image:
	case *ast.Code:
	case *ast.MathBlock:
	case *ast.Subscript:
	case *ast.Superscript:
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (r *Renderer) Text(w io.Writer, node *ast.Text, entering bool) {
	r.out(w, text.WrapBytes(node.Literal, r.opts.TextWidth))
}

// RenderHeader writes HTML document preamble and TOC if requested.
func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
}

// RenderFooter writes HTML document footer.
func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {
}

func (r *Renderer) writeDocumentHeader(w io.Writer) {
}
