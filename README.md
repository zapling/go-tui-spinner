# go-tui-spinner

A dead simple TUI spinner

# Usage

```go
package main

import (
	"os"
	"time"

	spinner "github.com/zapling/go-tui-spinner"
)

func main() {
	s := spinner.New(os.Stdout)
	stop := s.RunAsync()

	s.Println("Some text while the spinner is spinning")

	// Some slow process
	for i := 0; i < 3; i++ {
		s.Println("Processing slow thing")
		time.Sleep(1 * time.Second)
	}

	stop()

	s.Println("We are done processing that thing")

	stop = s.RunAsync()

	s.Println("Processing one last thing")

	time.Sleep(1 * time.Second)
	stop()

	s.Println("We are fully done!")
}
```
