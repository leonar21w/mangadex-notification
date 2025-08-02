package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"gopkg.in/gomail.v2"
)

type MangaUpdateEmailData struct {
	MangaTitle string
	MangaURL   string
	Chapters   []models.FeedChapter
}

const mangaUpdateTemplate = `
<html>
  <body>
    <h2>New chapters for “{{ .MangaTitle }}”</h2>
    <ul>
      {{- range .Chapters }}
      <li style="margin-bottom:1em;">
        <strong>Chapter {{ .Attributes.Chapter }}</strong>{{ if .Attributes.Title }}: {{ .Attributes.Title }}{{ end }}<br>
        <a href="{{ .Attributes.ExternalURL }}">Read now</a>
      </li>
      {{- end }}
    </ul>
    <p><a href="{{ .MangaURL }}">See all chapters on MangaDex</a></p>
  </body>
</html>
`

// SendMangaUpdateEmail renders your feed into HTML + text and sends it via SMTP.
// Requires these env vars: SENDER_EMAIL, RECIPIENT_EMAIL, EMAIL_APP_PASSWORD
func SendMangaUpdateEmail(data MangaUpdateEmailData) error {
	sender := os.Getenv("SENDER_EMAIL")
	recipient := os.Getenv("RECIPIENT_EMAIL")
	password := os.Getenv("EMAIL_APP_PASSWORD")

	// 1) Render HTML template (no .Format anywhere)
	tmpl, err := template.New("mangaUpdate").Parse(mangaUpdateTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	// 2) Build the email
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "📚 Manga Update: "+data.MangaTitle)
	m.SetBody("text/html", buf.String())

	// 3) Plain-text fallback
	chapters := data.Chapters
	buildMessage := ""
	for _, chapter := range chapters {
		buildMessage += fmt.Sprintf("Volume: %s, Chapter: %s\n", chapter.Attributes.Volume, chapter.Attributes.Chapter)
	}
	plain := fmt.Sprintf("New chapters for “%s”\nVisit: %s\n\n%s", data.MangaTitle, data.MangaURL, buildMessage)

	m.AddAlternative("text/plain", plain)

	// 4) Send via Gmail SMTP (or swap host/port)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, sender, password)
	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
