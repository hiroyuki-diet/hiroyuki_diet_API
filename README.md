# hiroyuki_diet_API

## プロジェクト概要

このプロジェクトは、ひろゆき氏のダイエットをテーマにしたアプリケーションのバックエンド API です。GraphQL を採用しており、ユーザー管理、食事記録、運動記録、アチーブメント、アイテム、スキン、ボイスなどの機能を提供します。

## 特徴

- **GraphQL API**: `gqlgen` を使用した強力な型付けされた API。
- **ユーザー認証**: JWT (JSON Web Token) を利用したセキュアな認証システム。
- **データベース**: PostgreSQL を利用し、`GORM` で ORM を実装。
- **Docker 対応**: 開発環境の構築が容易。
- **パスワードハッシュ化**: `bcrypt` による安全なパスワード保存。

## 使用技術

- **Go**: バックエンド言語
- **GraphQL**: API クエリ言語
- **gqlgen**: Go 言語用 GraphQL サーバーフレームワーク
- **GORM**: Go 言語用 ORM ライブラリ
- **PostgreSQL**: データベース
- **Docker / Docker Compose**: 開発環境構築
- **JWT**: 認証トークン
- **bcrypt**: パスワードハッシュ化

## セットアップ

### 前提条件

- [Go](https://golang.org/doc/install) (バージョン 1.16 以上推奨)
- [Docker](https://docs.docker.com/get-docker/) および [Docker Compose](https://docs.docker.com/compose/install/)

### 1. リポジトリのクローン

```bash
git clone https://github.com/moXXcha/hiroyuki_diet_API.git
cd hiroyuki_diet_API
```

### 2. 環境変数の設定

プロジェクトのルートディレクトリに `.env` ファイルを作成し、以下の内容を記述してください。

```dotenv
# PostgreSQL データベース設定
POSTGRES_HOST=db
POSTGRES_PORT=5432

# JWT シークレットキー (任意の強力な文字列を設定してください)
JWT_SECRET_KEY=your_super_secret_jwt_key_here

# アプリケーションポート
PORT=8080
```

### 3. データベースのセットアップ

Docker Compose を使用して PostgreSQL コンテナを起動し、データベースを初期化します。

```bash
docker-compose up -d db
```

データベースの初期化スクリプト (`init.sql`) が自動的に実行されます。

### 4. Go モジュールのインストール

`backend` ディレクトリに移動し、必要な Go モジュールをインストールします。

```bash
cd backend
go mod tidy
```

### 5. アプリケーションの実行

`backend` ディレクトリからアプリケーションを起動します。

```bash
go run server.go
```

または、Docker Compose を使用してすべてのサービスを起動することもできます。

```bash
docker-compose up -d
```

アプリケーションはデフォルトで `http://localhost:8080` で起動します。

#### Docker コンテナ内での作業

アプリケーションが Docker コンテナ内で実行されている場合、以下のコマンドでコンテナのシェルに入ることができます。

```bash
docker-compose exec container_name: hiroyuki_diet_app sh
```

コンテナ内で Go コマンドを実行したり、ログを確認したりする際に便利です。

## GraphQL スキーマ

GraphQL Playground にアクセスすると、API スキーマをインタラクティブに探索できます。

**GraphQL Playground URL**: `http://localhost:8080/`
**GraphQL API エンドポイント**: `http://localhost:8080/query`

### スキーマ定義 (`backend/graph/schema.graphqls`)

```graphql
# GraphQL schema example
#
# https://gqlgen.com/getting-started/
enum FieldEnum {
  login
  signin
  home
  meal
  meal_form
  meal_edit
  data
  profile
  profile_edit
  exercise
  eexercise_complete
  achievement
  achievement_complete
  chibi_hiroyuki
}
enum SkinPartEnum {
  head
  face
  body
}

enum MealTypeEnum {
  breakfast
  lunch
  dinner
  snacking
}

enum GenderEnum {
  man
  woman
}
type Query {
  user(id: ID!): User!
  foods: [Food!]!
}

type Mutation {
  signUp(input: Auth!): ID
  tokenAuth(input: InputTokenAuth!): JWTTokenResponse!
  login(input: Auth!): JWTTokenResponse!
  logout(input: ID!): ID
  createExercise(input: InputExercise!): ID
  editExercise(input: InputExercise!): ID
  receiptAchievement(input: InputAchievement!): ID
  createProfile(input: InputProfile!): ID
  editProfile(input: InputProfile!): ID
  createMeal(input: InputMeal!): ID
  editMeal(input: InputMeal!): ID
  deleteMeal(input: ID!): ID
  postSkin(input: InputPostSkin!): ID
  useItem(input: InputUseItem!): ID
}
input InputUseItem {
  userId: ID!
  itemId: ID!
  count: Int!
}

input InputTokenAuth {
  userId: ID!
  token: Int!
}

input InputPostSkin {
  userId: ID!
  skinId: ID!
}

input InputAchievement {
  userId: ID!
  achievementId: ID!
}

input Auth {
  email: String!
  password: String!
}
input InputExercise {
  userId: ID
  time: Int!
}
input InputProfile {
  userId: ID!
  userName: String!
  age: Int!
  gender: String!
  weight: Int!
  height: Int!
  targetWeight: Int!
  targetDailyExerciseTime: Int!
  targetDailyCarorie: Int!
}
input InputMeal {
  userId: ID
  mealId: ID
  mealType: MealTypeEnum!
  foods: [ID!]!
}

type User {
  id: ID!
  email: String!
  profile: Profile!
  level: Int!
  signUpToken: SignUpToken!
  isTokenAuthenticated: Boolean!
  experiencePoInt: Int!
  exercisies(offset: String!, limit: String!): [Exercise!]
  meals: [Meal!]
  meal(id: ID!): Meal!
  items: [ItemResponse!]!
  hiroyukiSkins(usingSkin: Boolean!): [SkinResponse!]!
  achievements: [AchievementResponse!]!
  hiroyukiVoicies(fields: [FieldEnum!]!): [HiroyukiVoiceResponse!]!
}

type HiroyukiVoiceResponse {
  id: ID!
  name: String!
  voiceUrl: String!
  releaseLevel: Int!
  fields: [FieldEnum!]!
  isHaving: Boolean!
}
type AchievementResponse {
  id: ID!
  name: String!
  isClear: Boolean!
}

type SkinResponse {
  id: ID!
  name: String!
  description: String!
  part: SkinPartEnum!
  skinImage: String!
  releaseLevel: Int!
  isUsing: Boolean!
  isHaving: Boolean!
}

type ItemResponse {
  id: ID!
  name: String!
  description: String!
  itemImage: String!
  count: Int!
}

type Meal {
  id: ID!
  mealType: MealTypeEnum!
  totalCalorie: Int!
  foods: [Food!]!
}

type SignUpToken {
  id: ID!
  token: Int!
  surviveTime: Int!
}

type Exercise {
  id: ID!
  time: Int!
  date: String!
}

type Profile {
  id: ID!
  userName: String!
  age: Int!
  gender: GenderEnum!
  weight: Int!
  height: Int!
  targetWeight: Int!
  targetDailyExerciseTime: Int!
  targetDailyCarorie: Int!
  isCreated: Boolean!
  favorability: Int
}

type Food {
  id: ID!
  name: String!
  estimateCalorie: Int!
  lastUsedDate: String!
}

type JWTTokenResponse {
  userId: ID!
  token: String!
}
```

## API 利用例

### ユーザー登録 (signUp)

```graphql
mutation {
  signUp(input: { email: "test@example.com", password: "password123" })
}
```

### トークン認証 (tokenAuth)

`signUp`後に返されるユーザー ID と、登録時に生成されたトークン（通常はメールなどでユーザーに通知される）を使用します。

```graphql
mutation {
  tokenAuth(input: { userId: "[signUpで返されたID]", token: 123456 }) {
    userId
    token
  }
}
```

### ログイン (login)

```graphql
mutation {
  login(input: { email: "test@example.com", password: "password123" }) {
    userId
    token
  }
}
```

### ユーザー情報取得 (user) - 認証が必要な例

このクエリを実行する際は、HTTP ヘッダーに認証トークンを含める必要があります。

**HTTP Headers:**
`Authorization: Bearer [loginまたはtokenAuthで返されたJWTトークン]`

```graphql
query {
  user(id: "[ユーザーのID]") {
    id
    email
    level
    profile {
      userName
    }
  }
}
```
