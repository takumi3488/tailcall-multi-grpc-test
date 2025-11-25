# Tailcall Multi-gRPC Test

Tailcall GraphQL Gatewayを使用して複数のgRPCマイクロサービスを集約するテストプロジェクト。

## アーキテクチャ

```
GraphQL Client
      │
      ▼
┌─────────────┐
│  Tailcall   │  (port 8000)
│  Gateway    │
└─────┬───────┘
      │
      ├──────────────┬──────────────┐
      ▼              ▼              │
┌───────────┐  ┌───────────┐        │
│NameService│  │ AgeService│        │
│ (50051)   │  │ (50052)   │        │
└───────────┘  └───────────┘        │
      │              │              │
      └──────────────┴──────────────┘
                     │
                     ▼
              users.json (データソース)
```

### サービス構成

| サービス | ポート | 説明 |
|---------|-------|------|
| NameService | 50051 | IDからユーザー名を返すgRPCサービス |
| AgeService | 50052 | IDからユーザー年齢を返すgRPCサービス |
| Tailcall Gateway | 8000 | 両gRPCサービスを1つのGraphQL APIとして公開 |

## 必要環境

- Go 1.21+
- [Tailcall](https://tailcall.run/)
- [buf](https://buf.build/) (Protobufコード生成用)
- [goreman](https://github.com/mattn/goreman) (プロセス管理)

## セットアップ

```bash
# 依存関係のインストール
go mod download

# Protobufコード生成
buf generate
```

## 起動方法

### 全サービス一括起動（推奨）

```bash
goreman start
```

### 個別起動

```bash
# NameService
go run cmd/name/main.go

# AgeService
go run cmd/age/main.go

# Tailcall Gateway
tailcall start main.yml
```

## 使用方法

### GraphQLクエリ例

```bash
curl -X POST http://localhost:8000/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ user(id: 1) { id name age } }"}'
```

### レスポンス例

```json
{
  "data": {
    "user": {
      "id": 1,
      "name": "田中太郎",
      "age": 28
    }
  }
}
```

## プロジェクト構成

```
.
├── cmd/
│   ├── name/main.go    # NameService実装
│   └── age/main.go     # AgeService実装
├── gen/go/             # buf生成コード（手動編集禁止）
├── proto/
│   ├── name/name.proto # NameService定義
│   └── age/age.proto   # AgeService定義
├── main.yml            # Tailcall設定ファイル
├── main.graphql        # GraphQLスキーマ
├── users.json          # データソース
├── Procfile            # goreman用プロセス定義
├── buf.yaml            # buf設定
└── buf.gen.yaml        # bufコード生成設定
```

## ライセンス

MIT
