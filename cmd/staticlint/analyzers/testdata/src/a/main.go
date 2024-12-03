package main

import "os"

func main() {
	os.Exit(1) // want "не допускается прямой вызов os.Exit в функции main"
}
