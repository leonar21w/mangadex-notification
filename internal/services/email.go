package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"gopkg.in/gomail.v2"
)

type MangaUpdateEmailData struct {
	MangaTitle string
	MangaURL   string
	Chapters   []models.FeedChapter
	SentAt     time.Time // optional; if zero, we'll set time.Now()
}

const mangaUpdateHTML = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Manga Update: {{ .MangaTitle }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      /* Client-friendly, inlined-safe styles (kept simple for Gmail/Yahoo/Outlook) */
      body { margin:0; padding:0; background:#0b0b0c; color:#e5e7eb; -webkit-font-smoothing:antialiased; }
      .preheader { display:none !important; visibility:hidden; opacity:0; color:transparent; height:0; width:0; overflow:hidden; mso-hide:all; }
      .container { max-width:640px; margin:0 auto; padding:24px 16px; font-family:ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, "Apple Color Emoji","Segoe UI Emoji"; }
      .card { background:#111214; border:1px solid #1f2227; border-radius:12px; padding:20px; }
      h1 { margin:0 0 8px 0; font-size:20px; line-height:1.25; color:#f9fafb; }
      .sub { color:#9ca3af; font-size:12px; margin-bottom:16px; }
      .list { list-style:none; padding:0; margin:0; }
      .item { padding:12px 0; border-bottom:1px solid #1f2227; }
      .item:last-child { border-bottom:none; }
      .ch-head { font-weight:600; color:#f3f4f6; }
      .ch-title { color:#d1d5db; }
      .btn { display:inline-block; margin-top:8px; padding:10px 14px; border-radius:10px; text-decoration:none; background:#2563eb; color:#ffffff !important; font-weight:600; }
      .btn:hover { filter:brightness(1.05); }
      .footer { color:#6b7280; font-size:12px; margin-top:18px; text-align:center; }
      a { color:#93c5fd; }
    </style>
  </head>
  <body>
    <div class="preheader">New chapters for {{ .MangaTitle }} are out.</div>
    <div class="container">
      <div class="card">
        <h1>New chapters for “{{ .MangaTitle }}”</h1>
        <div class="sub">Sent {{ .SentAt.Format "Jan 2, 2006 3:04 PM MST" }}</div>

        <ul class="list">
          {{- range .Chapters }}
          <li class="item">
            <div class="ch-head">Chapter {{ .Attributes.Chapter }}{{ if .Attributes.Volume }} · Vol. {{ .Attributes.Volume }}{{ end }}</div>
            {{- if .Attributes.Title }}
              <div class="ch-title">{{ .Attributes.Title }}</div>
            {{- end }}
            <a class="btn" href="{{ chapterURL . }}">Read chapter</a>
          </li>
          {{- end }}
        </ul>

        <p style="margin-top:18px;">
          <a class="btn" href="{{ .MangaURL }}">See all chapters on MangaDex</a>
        </p>
      </div>

      <div class="footer">
        You’re getting this because you follow updates for {{ .MangaTitle }}.
      </div>
    </div>
  </body>
</html>
`

func SendMangaUpdateEmail(data MangaUpdateEmailData) error {
	sender := os.Getenv("SENDER_EMAIL")
	recipient := os.Getenv("RECIPIENT_EMAIL")
	password := os.Getenv("EMAIL_APP_PASSWORD")

	if sender == "" || recipient == "" || password == "" {
		return fmt.Errorf("missing email env vars: SENDER_EMAIL / RECIPIENT_EMAIL / EMAIL_APP_PASSWORD")
	}

	if data.SentAt.IsZero() {
		data.SentAt = time.Now()
	}

	// FuncMap lets the template compute a link per chapter, with fallback.
	funcs := template.FuncMap{
		"chapterURL": func(ch models.FeedChapter) string {
			if ch.Attributes.ExternalURL != "" {
				return ch.Attributes.ExternalURL
			}
			// Fallback to MangaDex chapter page using ID
			// (safe guess for a usable link structure)
			return fmt.Sprintf("https://mangadex.org/chapter/%s", ch.ID)
		},
	}

	// Render HTML
	tmpl, err := template.New("mangaUpdateHTML").Funcs(funcs).Parse(mangaUpdateHTML)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	// Build plain-text alternative (readable in CLI/older clients)
	textBody := buildPlainTextBody(data)

	// Build and send
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "📚 Manga Update — "+data.MangaTitle)
	m.SetBody("text/html", htmlBuf.String())
	m.AddAlternative("text/plain", textBody)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, sender, password)
	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func buildPlainTextBody(data MangaUpdateEmailData) string {
	var b strings.Builder
	fmt.Fprintf(&b, "New chapters for %q\nSent %s\n\n", data.MangaTitle, data.SentAt.Format(time.RFC1123))
	for _, ch := range data.Chapters {
		chURL := ch.Attributes.ExternalURL
		if chURL == "" {
			chURL = fmt.Sprintf("https://mangadex.org/chapter/%s", ch.ID)
		}
		if ch.Attributes.Title != "" {
			fmt.Fprintf(&b, "• Ch. %s", ch.Attributes.Chapter)
			if ch.Attributes.Volume != "" {
				fmt.Fprintf(&b, " (Vol. %s)", ch.Attributes.Volume)
			}
			fmt.Fprintf(&b, ": %s\n  %s\n\n", ch.Attributes.Title, chURL)
		} else {
			title := fmt.Sprintf("Ch. %s", ch.Attributes.Chapter)
			if ch.Attributes.Volume != "" {
				title += fmt.Sprintf(" (Vol. %s)", ch.Attributes.Volume)
			}
			fmt.Fprintf(&b, "• %s\n  %s\n\n", title, chURL)
		}
	}
	if data.MangaURL != "" {
		fmt.Fprintf(&b, "All chapters: %s\n", data.MangaURL)
	}
	return b.String()
}
