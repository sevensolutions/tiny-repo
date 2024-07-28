package main

import (
	"github.com/joho/godotenv"
	"github.com/sevensolutions/tiny-repo/cmd"
)

func main() {
	godotenv.Load()

	cmd.Execute()
}
