# go-sqlboiler-graphql-boilerplate
Go SQLboilerのGraphQLのボイラープレート

## 技術構成
- go
- sqlboiler
- sql-migrate
- gqlgen
- ozzo-validation
- godotenv
- go-txdb
- stretchr/testify
- go-randomdata
- factory-go/factory

## 機能
### 認証
| | URI | 権限 |
| ------------- | ------------- | ------------- |
| 会員登録 | signUp(Mutation) | 権限なし |
| 会員登録 | signIn(Mutation) | 権限なし |

### TODOリスト
| | URI | 権限 |
| ------------- | ------------- | ------------- |
| 作成 | createTodo(Mutation) | 認証済み |
| 一覧取得 | fetchTodoLists(Query) | 認証済み |
| 詳細 | fetchTodo(Query) | 認証済み |
| 更新 | updateTodo(Mutation) | 認証済み |
| 削除 | updateTodo(Mutation) | 認証済み |

## 環境構築
### 1. 環境変数のファイルの作成
- rootディレクトリ配下の.env.sampleをコピーし、.envとする
- .env.sampleは開発環境をサンプルとしているため、設定値の調整は不要

```
cp .env.sample .env
```

### 2. Dockerのbuild・立ち上げ
- ビルド
```
docker-compose build
```

- 起動
```
docker-compose up -d
```

- コンテナに入る
```
docker-compose exec api_server sh
```

### 3. 開発環境DBのマイグレーション
- 2で起動・コンテナに入った上で、以下のコマンドを実行

マイグレーションファイルの作成
```
cd db

godotenv -f /app/.env sql-migrate new -env="mysql" <filename>
```

マイグレーション実行
```
cd db

godotenv -f /app/.env sql-migrate up -env="mysql"
```

### 4. Webサーバの起動
- 2の起動・コンテナに入った上で、以下のコマンドを実行
```
cd /app

godotenv -f /app.env go run server.go
```

4により、Webサーバの起動ができたら、graphiqlでアクセスができていることを確認

http://localhost:8080

## sqlboilerによるコードの自動生成
コンテナに入った上で、以下のコマンドを実行
```
make prepare-sqlboiler

make gen-models
```

## GraphQLの開発
- app/graphql配下の*.graphqlsファイルにスキーマを書く

- コンテナに入った上で、以下のコマンドを実行
```
make gen-gql
```

- 自動生成されたResolverに処理を書いていく

## 設計方針
- Resolver - Serviceのレイヤードアーキテクチャ
	- ロジックはServiceに寄せる
	- Controllerはリクエスト・レスポンスのハンドリング
		- リクエスト・レスポンスそれぞれのデータの加工含む
	- 将来的に複数モデルの保存などを行うケースが出てきたら、Service層からTransactionsクラスに切り出したりすると良さそう

## テスト方針
- それぞれの層でテストを書く
	- DB接続があるところはモックを使わずに行う(実環境に近い形でテストする方が不具合検知に役立つため)
		- サービスが大きくなってくると、モックを使用せず結合テストした方がリグレッションテストにも寄与するため

- Resolverで処理が意図通り行われ、意図したレスポンスが返ってくるかのテストを書く
	- Resolver単体は最初はすごく薄いと思うが(Service層の内容をそのまま返しているだけなので)、Service層のテストができていれば、Resolverのテストはメリットが薄いという考えもあるが、開発が進んでくると、同じFieldを複数のResolverから取得するようになってくると思うので、リクエストレベルできちんと意図した型のデータが返ってくるか確かめたいケースが出てくると考える
	- CIでの課金の心配がない状態であれば(Playwrightとか使えば)、バックエンドの修正だけでもE2E走らせてあげるのもアリかな？と思う一方、そうなるとCIの実行時間が長くなりすぎて、リリース速度に影響しそう
	- バックエンドの不具合はバックエンドのテストで拾えるのが理想かと考える

- 正常系/異常系ともに書く
	- 可能な限りC1カバレッジで書きたいところ
	- 事故があるとまずい機能については、C2カバレッジで書いても良さそう

## テスト実行
### テスト用DBの作成・マイグレーション
- テスト用のDBをdbコンテナのホストにログインし、DB名`go_sqlboiler_graphql_boilerplate_test`で作成する
- DBを作成した上で、api_serverコンテナに入った上で、以下のコマンドを実行
```
cd db

# テスト用のDBのマイグレーション
godotenv -f /app/.env.test.local sql-migrate up -env="mysql"
```

### テスト実行
api_serverコンテナに入った上で、以下のコマンドを実行
```
godotenv -f /app/.env.test.local go test -v ./...
```
