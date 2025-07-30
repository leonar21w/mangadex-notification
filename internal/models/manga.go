package models

type Manga struct {
	ID             string
	CanonicalTitle string //EN if available if not jp or ID
	Chapters       []MangadexChapterData
	CoverURL       string
}
