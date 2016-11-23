package templates

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"strconv"
	"time"

	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/lang"
)

type htmlWriter struct {
	bytes.Buffer
}

// Allows passing additional information to thread-related templates
type postContext struct {
	common.Post
	OP             int64
	Board, Subject string
	Lang           lang.Pack
}

// Write an element attribute to the buffer
func (w *htmlWriter) attr(key, val string) {
	w.WriteByte(' ')
	w.WriteString(key)
	if val != "" {
		w.WriteString(`="`)
		w.WriteString(val)
		w.WriteByte('"')
	}
}

// Write an HTML-escaped string to buffer
func (w *htmlWriter) escape(s string) {
	w.WriteString(html.EscapeString(s))
}

// Write an URL-query-escaped string to buffer
func (w *htmlWriter) queryEscape(s string) {
	w.WriteString(url.QueryEscape(s))
}

// Outputs the buffer contents as a HTML string
func (w *htmlWriter) HTML() template.HTML {
	return template.HTML(w.String())
}

func wrapPost(
	p common.Post,
	op int64,
	board, subject string,
	lang lang.Pack,
) postContext {
	return postContext{
		Post:    p,
		OP:      op,
		Board:   board,
		Subject: subject,
		Lang:    lang,
	}
}

// Returns the HTTP path to the thumbnail of an image
func thumbPath(img common.Image) string {
	var ext string
	if img.FileType == common.JPEG {
		ext = "jpg"
	} else {
		ext = "png"
	}
	return fmt.Sprintf("/images/thumb/%s.%s", img.SHA1, ext)
}

// Returns the HTTP path to the source file
func sourcePath(img common.Image) string {
	return fmt.Sprintf(
		"/images/src/%s.%s",
		img.SHA1,
		common.Extensions[img.FileType],
	)
}

func extension(fileType uint8) string {
	return common.Extensions[fileType]
}

// Renders the post creation time field
func renderTime(sec int64) string {
	return time.Unix(sec, 0).Format("2 Jan 2006 (Mon) 15:04")
}

// Renders a human-readable representation video/audio length
func readableLength(l uint32) string {
	if l < 60 {
		return fmt.Sprintf("0:%02d", l)
	}
	min := l / 60
	return fmt.Sprintf("%02d:%02d", min, l-min)
}

// Renders a human-readable representation of file size
func readableFileSize(s int) string {
	if s < (1 << 10) {
		return fmt.Sprintf("%d B", s)
	}
	if s < (1 << 20) {
		return fmt.Sprintf("%d KB", s/(1<<10))
	}
	return fmt.Sprintf("%.1f MB", float32(s)/(1<<20))
}

// Render a link to another post. Can optionally be cross-thread.
func renderPostLink(id, op int64, board string, cross bool) template.HTML {
	var text, url string
	if !cross {
		text = strconv.FormatInt(id, 10)
		url = fmt.Sprintf("#p%d", id)
	} else {
		text = fmt.Sprintf(">/%s/%d", board, id)
		url = fmt.Sprintf("/%s/%d#p%d", board, op, id)
	}
	return template.HTML(fmt.Sprintf("<a href=\"%s\">>>%s</a>", url, text))
}
