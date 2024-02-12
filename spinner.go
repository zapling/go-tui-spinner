package spinner

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

const clearLineAnsiSeq = "\033[2K\r"

var defaultFaces = []string{"|", "/", "—", "\\", "|", "/", "—", "\\"}

func New(out io.Writer) *Spinner {
	return &Spinner{Out: out, faces: defaultFaces}
}

func (s *Spinner) WithFaces(f []string) *Spinner {
	s.faces = f
	return s
}

func (s *Spinner) WithText(t string) *Spinner {
	s.text = t
	return s
}

type Spinner struct {
	Out   io.Writer
	faces []string
	text  string

	isDone  atomic.Bool
	printCh chan []any
	textCh  chan string
}

func (s *Spinner) Run(ctx context.Context) {
	s.isDone.Swap(false)

	printCh := make(chan []any)
	textCh := make(chan string)

	go s.run(ctx, printCh, textCh)
	s.printCh = printCh
	s.textCh = textCh
}

func (s *Spinner) Println(a ...any) {
	if s.isDone.Load() {
		s.renderPrintln(a...)
		return
	}
	s.printCh <- a
}

func (s *Spinner) SetText(t string) {
	if s.isDone.Load() {
		s.text = t
		return
	}
	s.textCh <- t
}

func (s *Spinner) run(ctx context.Context, printCh chan []any, textCh chan string) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer close(printCh)
	defer close(textCh)

	text := s.text
	faces := s.faces
	faceIndex := s.renderFace(0, faces, text)

	for {
		select {
		case <-ctx.Done():
			s.isDone.Swap(true)
			s.clearLine()
			return
		case <-ticker.C:
			s.clearLine()
			faceIndex = s.renderFace(faceIndex, faces, text)
		case newText := <-textCh:
			text = newText
		case values := <-printCh:
			s.clearLine()
			s.renderPrintln(values...)
			s.renderFace(faceIndex, faces, text)
		}
	}
}

func (s *Spinner) renderPrintln(a ...any) {
	fmt.Fprintln(s.Out, a...)
}

func (s *Spinner) renderFace(index int, faces []string, text string) int {
	str := faces[index]
	if text != "" {
		str += " " + text
	}
	fmt.Fprint(s.Out, str)
	index++
	if index == len(faces) {
		return 0
	}
	return index
}

func (s *Spinner) clearLine() {
	fmt.Fprint(s.Out, clearLineAnsiSeq)
}
