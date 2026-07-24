# doty-lore-bot
A Discord Bot for Logging World Building Related Events in a DnD Campaign

## Build
- If you have a Golang compiler installed, you can run `go build -o dotty ./cmd/bot`
- IF you prefer using Docker, you can build the image and run it in a container

## Deployment
- Database for lore cards is hosted on supabase
- Server needs to run on laptop by running `go run /cmd/bot/main.go`
