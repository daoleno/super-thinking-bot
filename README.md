# Super Thinking Telegram Bot

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/mod/github.com/daoleno/super-thinking-bot)
![Build](https://github.com/daoleno/super-thinking-bot/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/daoleno/super-thinking-bot)](https://goreportcard.com/report/github.com/daoleno/super-thinking-bot)

Collect the mental models in Book [Super Thinking: The Big Book of Mental Models](https://www.amazon.com/dp/0525533583/ref=cm_sw_r_tw_dp_x_ADNIFb61NSMW9).

Push one mental model a day to the [@superthinking2u](https://t.me/superthinking2u) telegram channel.

## build

```sh
go build .

# For RaspberryPi
GOOS=linux GOARCH=arm GOARM=7 go build .
```

## run

```
./super-thinking-bot --start "FAT-TAILED DISTRIBUTIONS"
```

## License

[MIT](LICENSE) License
