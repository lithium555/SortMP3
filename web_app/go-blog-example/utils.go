package main

import (
	"crypto/rand"
	"fmt"

	"github.com/russross/blackfriday"
)

func GenerateID() string { // method generate random ID
	b := make([]byte, 18)
	rand.Read(b)                // читаем рандомные числа в наш массив байт
	return fmt.Sprintf("%x", b) // потом печатаем и создаем из этих байт строку
}

func ConvertMarkdownToHtml(mardown string) string {
	return string(blackfriday.Run([]byte(mardown)))
}
