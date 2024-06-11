package main

import (
	"blog-system/internal/app"
	"fmt"
	"os"
)

func main() {
	if err := app.Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot start the application: %v", err)
		os.Exit(1)
	}
}
