package main

import "go/format"

func formatSource(src []byte) ([]byte, error) {
	return format.Source(src)
}
