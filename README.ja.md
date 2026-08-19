[English](README.md) | [日本語](README.ja.md)

# Cafe Inventory Management API

**Go**（`net/http`）と **PostgreSQL** で構築された、**複数店舗展開のカフェチェーン**を管理するための本番想定のRESTfulバックエンドAPI。

このシステムは、単一ブランドの下で複数のカフェ店舗を運営する企業向けに設計されています。本部はカフェ、商品、ユーザー、在庫を一元管理でき、各店舗のスタッフはリアルタイムで在庫数を更新できます。また、顧客が近くのカフェや在庫のある商品を見つけられるよう、位置情報ベースの検索機能も提供します。

---

# ビジネスシナリオ

東京で複数のカフェ店舗を展開する **ABC Coffee** という企業を想定してください。

事業が成長するにつれて、いくつかの運営上の課題が生じます。

- 店舗ごとに在庫レベルが異なる。
- 本部は全店舗の在庫状況を把握する必要がある。
- 商品カタログは全店舗で統一されている必要がある。
- マネージャーは誰が在庫を更新したかの監査証跡を必要とする。
- 顧客は特定の商品を在庫している最寄りの店舗を見つけたい。

このAPIは、一元管理、ロールベースのアクセス制御、在庫追跡、位置情報ベースの検索機能を提供することで、これらの課題に対応します。

---

# 機能

## 認証

- JWT認証
- ロールベースの認可（Admin / Staff）

## カフェ管理

- カフェの作成
- カフェ情報の更新
- カフェの削除
- カフェ一覧・詳細の取得

## 商品管理

- 商品の作成
- 商品情報の更新
- 商品の削除
- 商品カタログの取得

## カテゴリ管理

- カテゴリの作成
- カテゴリ名の更新
- カテゴリの削除
- カテゴリ一覧・詳細の取得

## 在庫管理

- 在庫レベルの確認
- 在庫数量の更新
- 在庫変動履歴の記録
- 各在庫操作を行ったユーザーの追跡

## ユーザー管理

- ユーザー登録
- ユーザー情報の更新
- ユーザーの削除
- ユーザーロールの管理

## 位置情報サービス

- カフェの座標（緯度・経度）の保存
- 最寄り駅によるカフェ検索
- 顧客の座標からの距離ベース検索（Haversine公式）、半径による絞り込みに対応
- 取得件数の制限

## 今後の機能

- 商品在庫検索（特定の商品を在庫している近隣のカフェを検索）
- Google Maps API連携
- 店舗パフォーマンスダッシュボード
- 売上分析

---

# 技術スタック

| カテゴリ | 技術 |
|----------|------------|
| 言語 | Go |
| HTTPサーバー | net/http |
| データベース | PostgreSQL |
| DBドライバ | pgx |
| 認証 | JWT |
| パスワードハッシュ化 | bcrypt |
| コンテナ | Docker |
| アーキテクチャ | レイヤードアーキテクチャ |
| APIスタイル | RESTful API |
| インフラ | AWS ECS |
| バージョン管理 | Git & GitHub |

---

# プロジェクト構成

```text
.github/
└── workflows/
    └── ci.yml

cmd/
└── server/
    └── main.go

internal/
├── database/
├── handler/
├── middleware/
├── model/
├── repository/
├── router/
├── service/
└── utils/

migrations/

Dockerfile
docker-compose.yml
go.mod
go.sum
README.md
```

---

# システムアーキテクチャ

```text
                Client
                   │
                   ▼
              HTTP Request
                   │
                   ▼
               Router
                   │
                   ▼
             Middleware
        (JWT Authentication,
           Request Logging)
                   │
                   ▼
               Handler
                   │
                   ▼
               Service
                   │
                   ▼
             Repository
                   │
                   ▼
             PostgreSQL
```

---

# データベース設計

## テーブル

- cafes
- products
- inventory
- inventory_logs
- users

### Entity Relationship

```text
cafes
   │
   │ 1:N
   ▼
products
   │
   │ 1:1
   ▼
inventory
   │
   │ 1:N
   ▼
inventory_logs
        ▲
        │
        │ N:1
      users
```

---

# API 概要

## 認証

| Method | Endpoint |
|---------|----------|
| POST | /api/v1/auth/login |

---

## カフェ

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/cafes |
| GET | /api/v1/cafes/{id} |
| POST | /api/v1/cafes |
| PATCH | /api/v1/cafes/{id} |
| DELETE | /api/v1/cafes/{id} |

---

## カテゴリ

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/categories |
| GET | /api/v1/categories/{id} |
| POST | /api/v1/categories |
| PATCH | /api/v1/categories/{id} |
| DELETE | /api/v1/categories/{id} |

---

## 商品

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/products |
| GET | /api/v1/products/{id} |
| POST | /api/v1/products |
| PATCH | /api/v1/products/{id} |
| DELETE | /api/v1/products/{id} |

---

## 在庫

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/inventory |
| GET | /api/v1/inventory/{id} |
| PATCH | /api/v1/inventory/{id} |
| GET | /api/v1/inventory/logs |

---

## ユーザー

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/users |
| GET | /api/v1/users/{id} |
| POST | /api/v1/users |
| PATCH | /api/v1/users/{id} |
| DELETE | /api/v1/users/{id} |

---

# 認証

このAPIはJWT（JSON Web Token）認証を使用します。

ログイン成功後、クライアントはリクエストヘッダーにアクセストークンを含める必要があります。

```http
Authorization: Bearer <JWT Token>
```

---

# 使い方

## 前提条件

- Docker & Docker Compose
- Go 1.26以上（コンテナ外でAPIを実行する場合のみ必要）

## 1. 環境変数を設定する

プロジェクトルートに `.env` ファイルを作成します。

```env
POSTGRES_USER=cafe
POSTGRES_PASSWORD=cafe
POSTGRES_DB=cafe_inventory
POSTGRES_PORT=5432

DATABASE_URL=postgres://cafe:cafe@localhost:5432/cafe_inventory?sslmode=disable

JWT_SECRET=<十分に長いランダムな文字列>
```

## 2. データベースを起動しマイグレーションを実行する

```bash
docker-compose up -d
```

これにより PostgreSQL が起動し、シードデータとして登録される管理者ユーザー（`admin@cafe.local`）を含む全マイグレーションが実行されます。

## 3. APIサーバーを起動する

```bash
go run cmd/server/main.go
```

これでAPIは `http://localhost:8080` で利用可能になります。

## 4. 一連の流れ：ログイン → データ作成 → 在庫更新 → 履歴確認

**シードされた管理者としてログイン:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@cafe.local","password":"admin123"}'
```

```json
{ "token": "<jwt>" }
```

このトークンを保存してください。以降の書き込み操作にはすべて `Authorization: Bearer <jwt>` が必要です。

**カテゴリを作成:**

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Coffee"}'
```

```json
{ "id": 1, "message": "Category created successfully" }
```

**カフェを作成:**

```bash
curl -X POST http://localhost:8080/api/v1/cafes \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ABC Coffee Shinjuku",
    "address": "1-1-1 Shinjuku, Tokyo",
    "latitude": 35.6895,
    "longitude": 139.6917,
    "nearest_station": "Shinjuku",
    "opening_time": "07:00",
    "closing_time": "22:00"
  }'
```

```json
{ "id": 1, "message": "Cafe created successfully" }
```

**商品を作成**（このとき `stock_quantity: 0` の在庫レコードも自動的に作成されます）:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "cafe_id": 1,
    "category_id": 1,
    "name": "Blend Coffee",
    "description": "House blend, medium roast",
    "price": "350.00"
  }'
```

```json
{ "id": 1, "message": "Product created successfully" }
```

**在庫を追加**（`operation` は `IN`、`OUT`、`ADJUST` のいずれか。ここでの `{id}` は商品のidを指します）:

```bash
curl -X PATCH http://localhost:8080/api/v1/inventory/1 \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"operation":"IN","change_quantity":50}'
```

```json
{ "message": "Inventory updated successfully", "stock_quantity": 50 }
```

**現在の在庫を確認:**

```bash
curl http://localhost:8080/api/v1/inventory
```

**在庫変動履歴を確認:**

```bash
curl http://localhost:8080/api/v1/inventory/logs \
  -H "Authorization: Bearer <jwt>"
```

## 5. APIの全体像を確認する（Swagger UI）

```
http://localhost:8080/swagger/index.html
```

**Authorize** をクリックし、手順4で取得した `Bearer <jwt>` を貼り付けることで、保護されたエンドポイントもブラウザから直接試すことができます。なお、生のスペックファイルは `docs/swagger.json` / `docs/swagger.yaml` からも参照できます。

---

# ユーザーロール

## Admin

- カフェの管理
- カテゴリの管理
- 商品の管理
- ユーザーの管理

## Staff

- ログインして公開エンドポイント（カフェ、商品、カテゴリ）を閲覧できる
- 現時点ではスタッフ専用の権限は実装されておらず、作成・更新・削除系のエンドポイントはすべてAdminロールが必要

---

# 検索機能

このAPIは以下をサポートします。

現在実装済み:

 - 最寄り駅によるカフェ検索（部分一致・大文字小文字を区別しない）
 - 顧客の座標からの距離ベース検索（Haversine公式）、近い順にソート、半径による絞り込みに対応
 - 取得件数の制限

例（駅名による検索）:

```http
GET /api/v1/cafes?station=Shinjuku%20station&limit=10
```

例（距離による検索）:

```http
GET /api/v1/cafes?lat=35.6895&lng=139.6917&radius=5&limit=10
```

カテゴリによる絞り込み、距離以外の並べ替え、ページ単位のページネーションは今後実装予定です。

---

# 学習目標

このプロジェクトは以下の知識を示すものです。

- RESTful API設計
- Go標準ライブラリ (`net/http`)
- レイヤードアーキテクチャ
- リポジトリパターン
- JWT認証
- ロールベースアクセス制御 (RBAC)
- PostgreSQLデータベース設計
- データベースマイグレーション
- 在庫監査ログ
- 位置情報と距離計算
- 絞り込み、並べ替え、ページネーション
- 本番想定のバックエンド開発
- コンテナ化
- 単体テスト
- 結合テスト
- Swagger / OpenAPI ドキュメント
- GitHub Actions CI/CD

---

# 今後の改善予定

- Google Maps API連携
- リフレッシュトークン認証
- 商品画像アップロード
- Redisキャッシュ
- 本部ダッシュボード

---

# ライセンス

このプロジェクトは教育目的、およびバックエンドエンジニアリングのポートフォリオとして作成されています。
