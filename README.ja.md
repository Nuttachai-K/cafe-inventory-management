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
- 取得件数の制限

## 今後の機能

- 在庫追跡（在庫数、変動履歴、監査証跡）
- 距離ベースのカフェ検索（Haversine公式）
- 顧客向けカフェ検索
- 商品在庫検索
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
| インフラ | AWS EC2 |
| バージョン管理 | Git & GitHub |

---

# プロジェクト構成

```text
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
| GET | /api/v1/inventory/{productId} |
| PATCH | /api/v1/inventory/{productId} |
| GET | /api/v1/inventory/{productId}/logs |

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

 - 最寄り駅によるカフェ検索
 - 取得件数の制限

例:

```http
GET /api/v1/cafes?station=Shinjuku%20station&limit=10
```
距離ベースの検索、カテゴリによる絞り込み、並べ替え、ページ単位のページネーションは今後実装予定です。

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

---

# 今後の改善予定

- Google Maps API連携
- リフレッシュトークン認証
- GitHub Actions CI/CD
- Swagger / OpenAPI ドキュメント
- 商品画像アップロード
- Redisキャッシュ
- 本部ダッシュボード

---

# ライセンス

このプロジェクトは教育目的、およびバックエンドエンジニアリングのポートフォリオとして作成されています。
