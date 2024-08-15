package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ItamarYuran/Chat/server" // Adjust this import path if necessary

	"golang.org/x/text/language"
)

func main() {
	// Replace with your service account key path
	keyPath := "../../server/keys.json"

	// Create a Translation client
	client, err := server.NewTranslationClient(keyPath)
	if err != nil {
		log.Fatalf("Failed to create translation client: %v", err)
	}
	defer client.Close()

	// Define your text and target language
	text := "Hello, world!"
	targetLanguage := language.French // Use the `golang.org/x/text/language` package to specify languages

	// Translate the text
	translatedTexts, err := server.TranslateText(context.Background(), client, []string{text}, targetLanguage)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	// Print the translated text
	fmt.Println("Original text:", text)
	fmt.Println("Translated text:", translatedTexts[0])
}
