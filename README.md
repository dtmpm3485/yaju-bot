# yaju-bot — Discordの会話に乱入するGo製ネタBot

[![PyPI](https://img.shields.io/pypi/v/yaju-bot)](https://pypi.org/project/yaju-bot/)
[![Python](https://img.shields.io/pypi/pyversions/yaju-bot)](https://pypi.org/project/yaju-bot/)
[![License](https://img.shields.io/github/license/dtmpm3485/yaju-bot)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/dtmpm3485/yaju-bot?style=social)](https://github.com/dtmpm3485/yaju-bot/stargazers)

**yaju-bot** は、Discordの会話にたまに返信で乱入する **Go製DiscordネタBot** です。

通常会話への低確率乱入に加えて、特定キーワードやメンションでも呼び出せます。Pythonから簡単に起動でき、サーバーごとの出現率・クールダウン・キーワード・対象チャンネル設定にも対応しています。

> 💬 普通にDiscordで会話中…
>
> 🤖 **突然Botが乱入**

## ✨ 特徴

- 💬 **Discordの通常会話へ低確率で乱入**
- 🔑 **キーワードで呼び出し可能**
- @️ **メンションでも呼び出し可能**
- 🐹 **Go製Bot本体**
- 🐍 **Pythonから簡単起動** — `run(token)` だけ
- 🎚️ **サーバー別に出現率を変更可能**
- ⏱️ **クールダウン設定**
- 📝 **呼び出しキーワードを追加・削除可能**
- 📺 **対象チャンネルを指定可能**
- 📊 **ステータス・統計確認**

## 🚀 すぐに使う

### 1. インストール

```bash
pip install -U yaju-bot
```

### 2. Botを起動

```python
from yaju_bot import run

run("DISCORD_BOT_TOKEN")
```

## 📖 主なコマンド

| コマンド | 内容 |
|---|---|
| `/yaju on` | Botを有効化 |
| `/yaju off` | Botを無効化 |
| `/yaju status` | 現在の設定を確認 |
| `/yaju stats` | 統計を確認 |
| `/yaju quote` | Botのセリフを表示 |
| `/yaju test` | 動作テスト |
| `/yaju mode` | 動作モードを変更 |
| `/yaju chance` | 乱入確率を設定 |
| `/yaju cooldown` | クールダウンを設定 |
| `/yaju keyword add\|remove\|list\|reset` | 呼び出しキーワードを管理 |
| `/yaju channel add\|remove\|list\|clear` | 対象チャンネルを管理 |

## 🔑 初期呼び出しワード

`野獣先輩` `やじゅ` `やじゅせん` `yaju` `yajuu` `114514` `810` `淫夢`

キーワードはサーバー側で追加・削除できます。

## 🔎 こんな人向け

- Discordサーバーに**ネタBot / meme bot**を入れたい
- 会話中にランダムで反応するBotが欲しい
- キーワード反応型のDiscord Botを使いたい
- Go製BotをPythonから簡単に起動したい
- サーバーごとに出現率やチャンネルを調整したい

## ❓ FAQ

### Discordの普通の会話にも反応しますか？

はい。設定された確率や条件に応じて、通常の会話へ低確率で乱入します。

### キーワードを自分で追加できますか？

できます。`/yaju keyword` 系のコマンドから追加・削除・一覧確認・リセットができます。

### 乱入しすぎないようにできますか？

できます。`/yaju chance` と `/yaju cooldown` で出現率や間隔を調整できます。

## 🤖 DiscordネタBotシリーズ

| Bot | 内容 |
|---|---|
| [meigen-bot](https://github.com/dtmpm3485/meigen-bot) | Discordの会話から名言・迷言を自動検出 |
| [senryu-bot](https://github.com/dtmpm3485/senryu-bot) | Discordの会話から川柳を自動検出 |
| **yaju-bot** | Discordの会話に低確率で乱入するネタBot |

## ⭐ 気に入ったら

面白かった・役に立った場合は、GitHubの **Star ⭐** を付けてもらえると開発の励みになります。

Issue・改善案・バグ報告も歓迎です。

## License

MIT License
