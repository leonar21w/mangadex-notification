package models

type FollowedMangaCollection struct {
	ClientID        string
	MangaCollection []MangadexMangaData
}

type ClientFollowedMangaCollectionResponse struct {
	Data   []MangadexMangaData `json:"data"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
	Total  int                 `json:"total"`
}

type MangadexMangaData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes MangaAttributes `json:"attributes"`
}

type MangaAttributes struct {
	Title     map[string]string   `json:"title"`
	AltTitles []map[string]string `json:"altTitles"`
	Links     map[string]string   `json:"links"`
	Tags      []MangaTag          `json:"tags"`
}

type MangaTag struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
	Group       string            `json:"group"`
	Version     int               `json:"version"`
}
