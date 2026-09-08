# yaju-bot

## インストール

```bash
pip install yaju-bot
```

## 使い方

```python
from yaju_bot import run

run("DISCORD_BOT_TOKEN")
```

## コマンド一覧

- `/yaju on` - Botを有効化
- `/yaju off` - Botを無効化
- `/yaju status` - 現在の設定を確認
- `/yaju stats` - 返信統計を表示
- `/yaju quote` - ランダムで語録を表示
- `/yaju test` - 文章判定をテスト
- `/yaju mode` - quiet / normal / chaos を変更
- `/yaju chance` - 通常の乱入率を変更
- `/yaju cooldown` - 乱入クールダウンを変更
- `/yaju keyword add` - 呼び出しワードを追加
- `/yaju keyword remove` - 呼び出しワードを削除
- `/yaju keyword list` - 呼び出しワード一覧
- `/yaju keyword reset` - 呼び出しワードを初期化
- `/yaju channel add` - 乱入対象チャンネルを追加
- `/yaju channel remove` - 乱入対象チャンネルを削除
- `/yaju channel list` - 乱入対象チャンネル一覧
- `/yaju channel clear` - 全チャンネルを対象に戻す

通常の会話には低確率で返信し、`野獣先輩`、`やじゅ`、`yaju`、`114514` などで呼ぶと返信します。
