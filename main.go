package main

import "github.com/xuxusheng/time-frequency-be/internal/app"

func main() {
	a := app.New()

	a.Listen(":8080")
}
