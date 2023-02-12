docker build -t go-telegram-bot .
docker run -e TELEGRAM_BOT_TOKEN={token} go-telegram-bot