package server

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// NewTranslationClient creates a new Google Cloud Translation client
func NewTranslationClient(keyPath string) (*translate.Client, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx, option.WithCredentialsFile(keyPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create translation client: %w", err)
	}
	return client, nil
}

// TranslateText translates the given text to the specified target language
func TranslateText(ctx context.Context, client *translate.Client, text []string, targetLang language.Tag) ([]string, error) {
	resp, err := client.Translate(ctx, text, targetLang, nil)
	if err != nil {
		return nil, err
	}

	var translations []string
	for _, translation := range resp {
		translations = append(translations, translation.Text)
	}
	return translations, nil
}
