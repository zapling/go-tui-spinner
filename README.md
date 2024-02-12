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
    ctx, cancel := context.WithCancel(context.Background())

    s := spinner.New(os.Stdout)
        WithText("Applying some things").
        WithFaces([]string{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"})

    s.Run(ctx)

    s.Println("Some progress is being made")

    cancel() // Stop the spinner

    s.Println("Message after the spinner is stopped")

    ctx, cancel = context.WithCancel(context.Background())

    s.Run(ctx) // Restart spinner

    s.Print("Spinner can be restarted")

    cancel()

    time.Sleep(100 * time.Millisecond) // Give the cleanup time to delete the spinner fully
}
```
