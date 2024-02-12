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
	return &Spinner{Out: out, Faces: defaultFaces}
}

func (s *Spinner) WithText(t string) *Spinner {
	s.Text = t
	return s
}

type Spinner struct {
	Out   io.Writer
	Faces []string
	Text  string

	isDone  atomic.Bool
	printCh chan []any
}

func (s *Spinner) Run(ctx context.Context) {
	s.isDone.Swap(false)

	printCh := make(chan []any)

	go s.run(ctx, printCh)
	s.printCh = printCh
}

func (s *Spinner) Println(a ...any) {
	if s.isDone.Load() {
		s.renderPrintln(a...)
		return
	}

	s.printCh <- a
}

func (s *Spinner) run(ctx context.Context, printCh chan []any) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	faces := s.Faces
	faceIndex := s.renderFace(0, faces)

	for {
		select {
		case <-ctx.Done():
			s.isDone.Swap(true)
			s.clearLine()
			return
		case <-ticker.C:
			s.clearLine()
			faceIndex = s.renderFace(faceIndex, faces)
		case values := <-printCh:
			s.clearLine()
			s.renderPrintln(values...)
			s.renderFace(faceIndex, faces)
		}
	}
}

func (s *Spinner) renderPrintln(a ...any) {
	fmt.Fprintln(s.Out, a...)
}

func (s *Spinner) renderFace(index int, faces []string) int {
	str := faces[index]
	if s.Text != "" {
		str += " " + s.Text
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
