package models

type Language string

const (
	EN Language = "en"
	JP Language = "jp"
	KR Language = "kr"
	ID Language = "id"
)

type Platform string

const (
	Mangadex Platform = "mangadex"
)

type Manga struct {
	ID             string
	CanonicalTitle string //EN if available if not jp or ID
	AltTitles      map[Language][]string
	Description    string
	Authors        []string
	Artists        []string
	Genres         []string
	Status         string
	CoverURL       string
	Source         map[Platform]SourceID
	Chapters       []Chapter
}

// per‑platform metadata so you never lose traceability
type SourceID string

// a normalized chapter entry
type Chapter struct {
	ID         string // your own UUID
	MangaID    string
	Volume     int
	Number     float64
	Title      string              // canonical chapter title
	AltTitles  map[Language]string // per‑lang overrides
	Scanlator  string
	SourceMeta map[Platform]SourceID
}
