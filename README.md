# yaju-bot

Discordの会話にたまに返信で乱入するBotです。特定キーワードやメンションでも呼べます。

## インストール

```bash
pip install yaju-bot
```

## 使い方

```python
from yaju_bot import run

run("DISCORD_BOT_TOKEN")
```

## 主なコマンド

- `/yaju on` / `/yaju off`
- `/yaju status` / `/yaju stats`
- `/yaju quote` / `/yaju test`
- `/yaju mode` / `/yaju chance` / `/yaju cooldown`
- `/yaju keyword add|remove|list|reset`
- `/yaju channel add|remove|list|clear`

初期呼び出しワード: `野獣先輩` `やじゅ` `やじゅせん` `yaju` `yajuu` `114514` `810` `淫夢`
