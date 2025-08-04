package models

import (
	"context"
)

type Manga struct {
	ID             string
	CanonicalTitle string //EN if available if not jp or ID
	Chapters       []MangadexChapterData
	CoverURL       string
}

type MangaRepo interface {
	InsertMangaWithID(ctx context.Context, mangaID string, manga *Manga) error
	InsertAllChapters(ctx context.Context, mangaID string, manga *Manga) error
	UpdateMangaChapters(ctx context.Context, mangaID string, chapters []FeedChapter) ([]FeedChapter, error)
	CacheMangaIDList(ctx context.Context, mangaID []MangadexMangaData) (int, error)

	GetMangaIDList(ctx context.Context) ([]string, error)
	GetMangaTitle(ctx context.Context, mangaID string) (string, error)
}
