package main

import (
    "flag"
    "log"
    "os"

    "github.com/dtmpm3485/yaju-bot/internal/app"
)

func main() {
    dataDir := flag.String("data-dir", "", "directory for yaju-bot data")
    flag.Parse()

    token := app.EnvToken()
    if token == "" {
        log.Fatal("Discord bot token is missing. Set DISCORD_TOKEN or YAJU_BOT_TOKEN.")
    }
    if err := app.Run(token, *dataDir); err != nil {
        log.Printf("yaju-bot stopped: %v", err)
        os.Exit(1)
    }
}
