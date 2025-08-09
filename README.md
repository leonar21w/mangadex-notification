# MangaDex Notifications

A backend service that monitors [MangaDex](https://mangadex.org) for new chapter releases and sends email notifications when updates are available.  
Built with **Go**, uses **Redis** for caching, and supports multiple environments (test & production).

---

## Features
- Track manga updates for your MangaDex account
- Send email alerts for new chapter releases
- Multiple environment support (Test / Production)
---

## Requirements
- [Go 1.20+](https://go.dev/dl/)
- [Redis](https://redis.io/download) (local or remote)
- Email service that supports app passwords (e.g., Gmail)
- MangaDex API credentials ([how to get them](https://api.mangadex.org/docs/))

---

## Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/leonar21w/mangadex-notifications.git
cd mangadex-notifications
go mod tidy
```
 ---

## Configuration (.env fields)

```env
# Current environment: "test" or "production"
CURRENT_ENV=

# Redis (Test Environment)
REDIS_URL_TEST=
REDIS_TOKEN_TEST=

# Redis (Production)
REDIS_TOKEN=
REDIS_URL=

# Emails
EMAIL_APP_PASSWORD=
SENDER_EMAIL=
RECIPIENT_EMAIL=

# MangaDex API credentials
MGDEX_USERNAME=
MGDEX_CLIENT=
MGDEX_PASSWORD=
MGDEX_SECRET=