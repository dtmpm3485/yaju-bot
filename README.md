# yaju-bot

[![PyPI](https://img.shields.io/pypi/v/yaju-bot)](https://pypi.org/project/yaju-bot/)
[![Python](https://img.shields.io/pypi/pyversions/yaju-bot)](https://pypi.org/project/yaju-bot/)
[![License](https://img.shields.io/github/license/dtmpm3485/yaju-bot)](LICENSE)

Discordの会話にたまに返信で乱入するBotです。特定キーワードやメンションでも呼び出せます。Bot本体はGoで実装されており、Pythonから起動できます。

![yaju-botの動作例](https://github.com/dtmpm3485/yaju-bot/blob/main/assets/demo.jpg?raw=true)

## 機能

- 通常の会話へ低確率で乱入
- キーワードで呼び出し
- メンションで呼び出し
- サーバーごとの出現率設定
- クールダウン設定
- 呼び出しキーワードの追加・削除
- 対象チャンネルの指定
- ステータス・統計表示

## インストール

```bash
pip install -U yaju-bot
```

## 起動

```python
from yaju_bot import run

run("DISCORD_BOT_TOKEN")
```

## コマンド

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

## 初期呼び出しワード

`野獣先輩` `やじゅ` `やじゅせん` `yaju` `yajuu` `114514` `810` `淫夢`

キーワードはサーバー側で追加・削除できます。

## 関連リポジトリ

- [meigen-bot](https://github.com/dtmpm3485/meigen-bot) - Discordの会話から名言・迷言を検出するBot
- [senryu-bot](https://github.com/dtmpm3485/senryu-bot) - Discordの会話から川柳を検出するBot

## License

MIT License
