package main

import (
	"auth/internal/app"
    _ "github.com/lib/pq"
)

func main() {
	app := app.New()
    if err := app.Run("8000"); err != nil {
        panic(err)
    }
}