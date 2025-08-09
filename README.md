# MangaDex Notifications

A backend service that monitors [MangaDex](https://mangadex.org) for new chapter releases and sends email notifications when updates are available.  
Built with **Go**, uses **Redis** for caching, and supports multiple environments (test & production).

---

## Features
- Track manga updates for your MangaDex account
- Cache API results in Redis to minimize requests
- Send email alerts for new chapter releases
- Multiple environment support (Test / Production)
- Secure credential management via `.env` file

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
git clone https://github.com/<your-username>/<your-repo>.git
cd <your-repo>
