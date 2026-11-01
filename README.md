# Midblog

A satirical blog site where all of the posts are content-locked and AI generated.

Built using [Go](https://go.dev/), [Templ](https://github.com/a-h/templ), [Tailwind](https://tailwindcss.com/), [DaisyUI](https://daisyui.com/), and [HTMX](https://htmx.org/).

Large amounts of text have been generated using [OpenAI's GPT-3](https://openai.com/gpt-3/).

## Running with Docker

**Production** (app + MySQL; copy `.env.example` to `.env` and fill in the passwords first):

```sh
docker compose up -d --build        # serves on :4231 (override with APP_PORT)
```

Migrations run automatically on startup. `OPENAI_KEY` is optional; without it the site runs
but the AI generation endpoints in the admin panel will fail.

**Development** (live reload via air, Parcel asset watcher, MySQL exposed on localhost:3306; no `.env` needed):

```sh
docker compose -f docker-compose.dev.yaml up --build
```

On macOS, if the repo lives in an iCloud-synced folder (e.g. `~/Documents`) with "Optimize Mac Storage"
on, evicted files show up as `Input/output error` inside containers. Keep the folder downloaded
(Finder → right-click → Keep Downloaded) or move the repo out of iCloud.

## Running locally

Make sure you've got an MySQL backend up and running. I use docker, you can run it however you like.

Also make sure you've got all of the `.env` file filled out properly. You'll need an OpenAI key.

Requires Go 1.27+ and Node 20+.

```sh
npm ci                     # install frontend deps (Tailwind 4, daisyUI 5, htmx 2)
npm run build              # build assets/ into internal/templates/static/
go tool templ generate     # regenerate *_templ.go from *.templ (templ is pinned in go.mod)
go run ./cmd/web           # or `air` for live reload
```

Always edit the `.templ` sources and regenerate; never edit the generated `*_templ.go` files.
