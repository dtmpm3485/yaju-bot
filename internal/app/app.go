package app

import (
    "embed"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "slices"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/dtmpm3485/yaju-bot/internal/engine"
    "github.com/dtmpm3485/yaju-bot/internal/store"
)

//go:embed quotes.json
var quoteFS embed.FS

type Bot struct {
    session *discordgo.Session
    store *store.Store
    quotes []engine.Quote
    lastMu sync.Mutex
    lastReply map[string]time.Time
}

func Run(token, dataDir string) error {
    if token == "" { return fmt.Errorf("DISCORD_TOKEN is empty") }
    if dataDir == "" { dataDir = os.Getenv("YAJU_DATA_DIR") }
    if dataDir == "" { dataDir = ".yaju-bot" }

    st, err := store.Open(filepath.Join(dataDir, "config.json")); if err != nil { return err }
    raw, err := quoteFS.ReadFile("quotes.json"); if err != nil { return err }
    var quotes []engine.Quote
    if err := json.Unmarshal(raw, &quotes); err != nil { return fmt.Errorf("load quotes: %w", err) }

    s, err := discordgo.New("Bot " + token); if err != nil { return err }
    s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
    b := &Bot{session:s, store:st, quotes:quotes, lastReply:map[string]time.Time{}}
    s.AddHandler(b.onReady)
    s.AddHandler(b.onMessage)
    s.AddHandler(b.onInteraction)
    if err := s.Open(); err != nil { return err }
    defer s.Close()

    if err := b.registerCommands(); err != nil { return err }
    log.Printf("yaju-bot connected as %s", s.State.User.String())
    select {}
}

func (b *Bot) onReady(s *discordgo.Session, _ *discordgo.Ready) {
    _ = s.UpdateGameStatus(0, "会話にたまに乱入中")
}

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author == nil || m.Author.Bot || m.GuildID == "" { return }
    cfg := b.store.Get(m.GuildID)
    _ = b.store.Update(m.GuildID, func(c *store.GuildConfig){ c.Seen++ })
    if !cfg.Enabled || !channelAllowed(cfg.Channels, m.ChannelID) { return }

    content := m.Content
    for _, u := range m.Mentions {
        if s.State.User != nil && u.ID == s.State.User.ID { content += " yaju"; break }
    }

    r := engine.Evaluate(content, cfg.SummonWords, cfg.Chance, cfg.Mode, b.quotes)
    if !r.Reply { return }
    if !r.Summoned && !b.cooldownReady(m.ChannelID, cfg.CooldownSeconds) { return }

    _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
        Content:r.Text,
        Reference:m.Reference(),
        AllowedMentions:&discordgo.MessageAllowedMentions{RepliedUser:false},
    })
    if err != nil { log.Printf("reply failed: %v", err); return }
    b.markReply(m.ChannelID)
    _ = b.store.Update(m.GuildID, func(c *store.GuildConfig){
        c.Replies++
        if r.Summoned { c.Summons++ }
    })
}

func channelAllowed(channels []string, channelID string) bool {
    return len(channels) == 0 || slices.Contains(channels, channelID)
}

func (b *Bot) cooldownReady(channelID string, seconds int) bool {
    if seconds <= 0 { return true }
    b.lastMu.Lock(); defer b.lastMu.Unlock()
    t, ok := b.lastReply[channelID]
    return !ok || time.Since(t) >= time.Duration(seconds)*time.Second
}
func (b *Bot) markReply(channelID string) { b.lastMu.Lock(); b.lastReply[channelID]=time.Now(); b.lastMu.Unlock() }

func (b *Bot) registerCommands() error {
    cmd := &discordgo.ApplicationCommand{
        Name:"yaju", Description:"野獣先輩Botの設定・テスト",
        Options:[]*discordgo.ApplicationCommandOption{
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"on", Description:"Botを有効化"},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"off", Description:"Botを無効化"},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"status", Description:"現在の設定を表示"},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"stats", Description:"統計を表示"},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"quote", Description:"ランダム語録を1つ表示"},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"test", Description:"文章判定をテスト", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionString, Name:"text", Description:"判定する文章", Required:true}}},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"mode", Description:"乱入モード", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionString, Name:"value", Description:"quiet / normal / chaos", Required:true, Choices:[]*discordgo.ApplicationCommandOptionChoice{{Name:"quiet",Value:"quiet"},{Name:"normal",Value:"normal"},{Name:"chaos",Value:"chaos"}}}}},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"chance", Description:"通常乱入率(%)", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionNumber, Name:"percent", Description:"0〜100", Required:true}}},
            {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"cooldown", Description:"通常乱入のクールダウン秒", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionInteger, Name:"seconds", Description:"0〜3600", Required:true}}},
            {Type:discordgo.ApplicationCommandOptionSubCommandGroup, Name:"keyword", Description:"呼び出しワード設定", Options:[]*discordgo.ApplicationCommandOption{
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"add", Description:"呼び出しワード追加", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionString,Name:"word",Description:"追加する語句",Required:true}}},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"remove", Description:"呼び出しワード削除", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionString,Name:"word",Description:"削除する語句",Required:true}}},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"list", Description:"一覧表示"},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"reset", Description:"初期値へ戻す"},
            }},
            {Type:discordgo.ApplicationCommandOptionSubCommandGroup, Name:"channel", Description:"乱入チャンネル設定", Options:[]*discordgo.ApplicationCommandOption{
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"add", Description:"対象チャンネル追加", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionChannel,Name:"target",Description:"チャンネル",Required:true}}},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"remove", Description:"対象チャンネル削除", Options:[]*discordgo.ApplicationCommandOption{{Type:discordgo.ApplicationCommandOptionChannel,Name:"target",Description:"チャンネル",Required:true}}},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"list", Description:"対象一覧"},
                {Type:discordgo.ApplicationCommandOptionSubCommand, Name:"clear", Description:"制限解除（全チャンネル）"},
            }},
        },
    }
    _, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
    return err
}

func (b *Bot) onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type != discordgo.InteractionApplicationCommand { return }
    data := i.ApplicationCommandData(); if data.Name != "yaju" || i.GuildID == "" { return }
    opts := data.Options; if len(opts)==0 { b.respond(i,"サブコマンドを指定してください。",true); return }
    root := opts[0]
    cfg := b.store.Get(i.GuildID)

    if root.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
        if len(root.Options)==0 { return }
        b.handleGroup(i, root.Name, root.Options[0], cfg); return
    }

    public := root.Name=="quote" || root.Name=="test"
    if !public && root.Name!="status" && root.Name!="stats" && !isAdmin(i) {
        b.respond(i,"この設定は「サーバー管理」権限が必要です。",true); return
    }

    switch root.Name {
    case "on":
        _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Enabled=true}); b.respond(i,"✅ yaju-bot を有効にしました。",true)
    case "off":
        _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Enabled=false}); b.respond(i,"⏹️ yaju-bot を無効にしました。",true)
    case "status":
        channels := "全チャンネル"; if len(cfg.Channels)>0 { channels = fmt.Sprintf("%dチャンネル",len(cfg.Channels)) }
        b.respond(i,fmt.Sprintf("有効: %t\nモード: %s\n通常乱入率: %.1f%%\nクールダウン: %d秒\n呼び出し語: %d個\n対象: %s",cfg.Enabled,cfg.Mode,cfg.Chance,cfg.CooldownSeconds,len(cfg.SummonWords),channels),true)
    case "stats": b.respond(i,fmt.Sprintf("監視: %d\n返信: %d\nキーワード呼び出し: %d",cfg.Seen,cfg.Replies,cfg.Summons),true)
    case "quote": b.respond(i,engine.RandomQuote(b.quotes),false)
    case "test":
        text:=root.Options[0].StringValue(); r:=engine.Evaluate(text,cfg.SummonWords,100,cfg.Mode,b.quotes)
        b.respond(i,fmt.Sprintf("判定カテゴリ: %s\nスコア: %.2f\n候補返信: %s",r.Category,r.Score,r.Text),true)
    case "mode":
        v:=root.Options[0].StringValue(); _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Mode=v}); b.respond(i,"モードを `"+v+"` に変更しました。",true)
    case "chance":
        v:=root.Options[0].FloatValue(); if v<0 {v=0}; if v>100 {v=100}; _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Chance=v}); b.respond(i,fmt.Sprintf("通常乱入率を %.1f%% に変更しました。",v),true)
    case "cooldown":
        v:=root.Options[0].IntValue(); if v<0 {v=0}; if v>3600 {v=3600}; _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.CooldownSeconds=int(v)}); b.respond(i,fmt.Sprintf("クールダウンを %d秒に変更しました。",v),true)
    }
}

func (b *Bot) handleGroup(i *discordgo.InteractionCreate, group string, sub *discordgo.ApplicationCommandInteractionDataOption, cfg store.GuildConfig) {
    if sub.Name!="list" && !isAdmin(i) { b.respond(i,"この設定は「サーバー管理」権限が必要です。",true); return }
    switch group {
    case "keyword":
        switch sub.Name {
        case "add":
            w:=strings.TrimSpace(sub.Options[0].StringValue()); if w=="" { b.respond(i,"空の語句は追加できません。",true); return }
            _=b.store.Update(i.GuildID,func(c *store.GuildConfig){ if !slices.Contains(c.SummonWords,w){c.SummonWords=append(c.SummonWords,w)} })
            b.respond(i,"呼び出しワード `"+w+"` を追加しました。",true)
        case "remove":
            w:=strings.TrimSpace(sub.Options[0].StringValue()); _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.SummonWords=remove(c.SummonWords,w)}); b.respond(i,"呼び出しワード `"+w+"` を削除しました。",true)
        case "list": b.respond(i,"呼び出しワード: `"+strings.Join(cfg.SummonWords,"`, `")+"`",true)
        case "reset": _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.SummonWords=engine.DefaultSummonWords()}); b.respond(i,"呼び出しワードを初期値へ戻しました。",true)
        }
    case "channel":
        switch sub.Name {
        case "add": id:=sub.Options[0].ChannelValue(nil).ID; _=b.store.Update(i.GuildID,func(c *store.GuildConfig){if !slices.Contains(c.Channels,id){c.Channels=append(c.Channels,id)}}); b.respond(i,"対象に <#"+id+"> を追加しました。",true)
        case "remove": id:=sub.Options[0].ChannelValue(nil).ID; _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Channels=remove(c.Channels,id)}); b.respond(i,"対象から <#"+id+"> を外しました。",true)
        case "list":
            if len(cfg.Channels)==0 { b.respond(i,"対象: 全チャンネル",true) } else { parts:=make([]string,len(cfg.Channels)); for n,id:=range cfg.Channels {parts[n]="<#"+id+">"}; b.respond(i,"対象: "+strings.Join(parts," "),true) }
        case "clear": _=b.store.Update(i.GuildID,func(c *store.GuildConfig){c.Channels=nil}); b.respond(i,"チャンネル制限を解除しました。",true)
        }
    }
}

func isAdmin(i *discordgo.InteractionCreate) bool { return i.Member != nil && i.Member.Permissions&discordgo.PermissionManageServer != 0 }
func remove(xs []string, target string) []string { out:=xs[:0]; for _,x:=range xs {if x!=target {out=append(out,x)}}; return out }
func (b *Bot) respond(i *discordgo.InteractionCreate, content string, ephemeral bool) {
    flags:=discordgo.MessageFlags(0); if ephemeral {flags=discordgo.MessageFlagsEphemeral}
    if err:=b.session.InteractionRespond(i.Interaction,&discordgo.InteractionResponse{Type:discordgo.InteractionResponseChannelMessageWithSource,Data:&discordgo.InteractionResponseData{Content:content,Flags:flags}}); err!=nil {log.Printf("interaction response failed: %v",err)}
}

func EnvToken() string {
    if t:=os.Getenv("DISCORD_TOKEN"); t!="" {return t}
    if t:=os.Getenv("YAJU_BOT_TOKEN"); t!="" {return t}
    return ""
}

func ParsePort(s string, fallback int) int { v,err:=strconv.Atoi(s); if err!=nil || v<1 || v>65535 {return fallback}; return v }
