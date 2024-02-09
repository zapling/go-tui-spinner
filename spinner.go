package spinner

import (
	"context"
	"fmt"
	"io"
	"time"
)

var defaultFaces = []string{"|", "/", "—", "\\", "|", "/", "—", "\\"}

func New(out io.Writer) *Spinner {
	return &Spinner{Out: out, Faces: defaultFaces}
}

type Spinner struct {
	Out   io.Writer
	Faces []string

	isPrintingText bool
	waitCtx        context.Context
	waitCtxCancel  context.CancelFunc
}

// Run starts the spinner in blocking way and runs until the context is cancelled.
func (s *Spinner) Run(ctx context.Context) {
	faceIndex := 0
	for {
		if ctx.Err() != nil {
			return
		}

		if s.isPrintingText {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// TODO: ctx should be parent of waitCtx?
		s.waitCtx, s.waitCtxCancel = context.WithCancel(context.Background())

		fmt.Fprint(s.Out, s.Faces[faceIndex])

		select {
		case <-s.waitCtx.Done():
			break
		case <-time.After(200 * time.Millisecond):
			break
		}

		faceIndex++
		if faceIndex == len(s.Faces) {
			faceIndex = 0
		}
		fmt.Fprint(s.Out, "\033[2K\r")
	}
}

// RunAsync starts the spinner within a gorutine. Returns a cancelFunc that stops the spinner.
func (s *Spinner) RunAsync() context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go s.Run(ctx)
	return cancelFunc
}

// Println can be used to print any text while the spinner is running. The spinner will be
// temporarily cleared while text is being printed, and resume its state when done.
// If the spinner is stopped text will be printed as expected.
func (s *Spinner) Println(a ...any) {
	s.isPrintingText = true
	s.waitCtxCancel()
	fmt.Fprint(s.Out, "\033[2K\r")

	fmt.Fprintln(s.Out, a...)
	s.isPrintingText = false
}
