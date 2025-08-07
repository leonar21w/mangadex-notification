package constants

// Internal unexported struct with private fields
type mdConstants struct{}

var MD = mdConstants{}

// Getter methods
func (mdConstants) APIBaseURL() string            { return "https://api.mangadex.org" }
func (mdConstants) AuthBaseURL() string           { return "https://auth.mangadex.org" }
func (mdConstants) AccessTokenPrefix() string     { return "accessToken:" }
func (mdConstants) RefreshTokenPrefix() string    { return "refreshToken:" }
func (mdConstants) DefaultPageLimit() int         { return 25 }
func (mdConstants) TimeLayout() string            { return "2006-01-02T15:04:05" }
func (mdConstants) UserAuthGrantType() string     { return "password" }
func (mdConstants) FollowedMangaEndpoint() string { return MD.APIBaseURL() + "/user/follows/manga" }
func (mdConstants) FeedMangaEndpoint() string     { return MD.APIBaseURL() + "/user/follows/manga/feed" }
func (mdConstants) ChapterMangaEndpoint() string  { return MD.APIBaseURL() + "/chapter" }
func (mdConstants) AuthEndpoint() string {
	return MD.AuthBaseURL() + "/realms/mangadex/protocol/openid-connect/token"
}
