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
	CoverURL  string
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

type MangadexChapterList struct {
	Chapters []MangadexChapterData
}

type MangadexChapterListResponse struct {
	Data   []MangadexChapterData `json:"data"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Total  int                   `json:"total"`
}

type MangadexChapterData struct {
	ID         string                   `json:"id"`
	Attributes MangadexChapterAttribute `json:"attributes"`
}

type MangadexChapterAttribute struct {
	Volume             string `json:"volume"`
	Chapter            string `json:"chapter"`
	Title              string `json:"title"`
	TranslatedLanguage string `json:"translatedLanguage"`
	ExternalURL        string `json:"externalUrl"`
	ReadableAt         string `json:"readableAt"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Version            int    `json:"version"`
}

type FeedResponse struct {
	Result   string        `json:"result"`
	Response string        `json:"response"`
	Data     []FeedChapter `json:"data"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
	Total    int           `json:"total"`
}

type FeedChapter struct {
	ID            string                   `json:"id"`
	Type          string                   `json:"type"`
	Attributes    MangadexChapterAttribute `json:"attributes"`
	Relationships []Relationship           `json:"relationships"`
}

// look for type == manga, this will give mangaID
type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
