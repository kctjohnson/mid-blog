# Midblog

A satirical blog site where all of the posts are content-locked and AI generated.

Built using [Go](https://go.dev/), [Templ](https://github.com/a-h/templ), [Tailwind](https://tailwindcss.com/), [DaisyUI](https://daisyui.com/), and [HTMX](https://htmx.org/).

Large amounts of text have been generated using [OpenAI's GPT-3](https://openai.com/gpt-3/).

## Running

Make sure you've got an MySQL backend up and running. I use docker, you can run it however you like.

Also make sure you've got all of the `.env` file filled out properly. You'll need an OpenAI key.

Run `go mod tidy`, then `npm i`, then `air`, and the blog should run.
