package spinner

import (
	"context"
	"fmt"
	"io"
	"time"
)

const clearLineAnsiSeq = "\033[2K\r"

var defaultFaces = []string{"|", "/", "—", "\\", "|", "/", "—", "\\"}

func New(out io.Writer) *Spinner {
	return &Spinner{Out: out, Faces: defaultFaces}
}

type Spinner struct {
	Out   io.Writer
	Faces []string

	isDone  bool
	printCh chan []any
}

func (s *Spinner) Run(ctx context.Context) {
	s.isDone = false

	doneCh := make(chan struct{})
	printCh := make(chan []any)

	go func() {
		<-doneCh
		s.isDone = true
	}()

	go s.run(ctx, printCh, doneCh)
	s.printCh = printCh
}

func (s *Spinner) Println(a ...any) {
	if s.isDone {
		s.renderPrintln(a...)
		return
	}

	s.printCh <- a
}

func (s *Spinner) run(ctx context.Context, printCh chan []any, doneCh chan struct{}) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	faces := s.Faces
	faceIndex := s.renderFace(0, faces)

	for {
		select {
		case <-ctx.Done():
			s.clearLine()
			doneCh <- struct{}{}
			return
		case <-ticker.C:
			s.clearLine()
			faceIndex = s.renderFace(faceIndex, faces)
		case values := <-printCh:
			s.clearLine()
			s.renderPrintln(values...)
		}
	}
}

func (s *Spinner) renderPrintln(a ...any) {
	fmt.Fprintln(s.Out, a...)
}

func (s *Spinner) renderFace(index int, faces []string) int {
	fmt.Fprint(s.Out, faces[index])
	index++
	if index == len(faces) {
		return 0
	}
	return index
}

func (s *Spinner) clearLine() {
	fmt.Fprint(s.Out, clearLineAnsiSeq)
}
