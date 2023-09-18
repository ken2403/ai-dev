# LLMを用いたDeveloper Tool開発

## 実験の目標

**より少ないコードサンプルと、より少ない指示で、より良いコードを生成できるようにする**

- CLIツールで使用することを考えると、対話しながらコードを生成することは難しい

## 実験の流れ

  1. 簡単なCRUD操作を行うAPIの作成
    - 簡単なUserテーブルに対するCRUD操作を行うAPIの作成を目指す
    - ChatGPTをチューニングして、アーキテクチャやフレームワーク、ORMマッパーに適合できるかを実験
    - ChatGPTのチューニング方針を決める

  2. チューニングしたChatGPTを用いて、複数のテーブルが絡むようなAPIの作成

  3. チューニングしたChatGPTを用いて、Issueに対応するAPIの作成

## 1. 簡単なCRUD操作を行うAPIの作成

### 1.1. Input

- DB Schema (`schema_1.sql`)

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `mail` varchar(100) NOT NULL,
  `birthday` DATE NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
```

- API Spec (`api_1.yaml`)

```yaml
OpenAPI: 3.0.3
info:
  title: CRUD API
  version: 1.0.0
paths:
  /api/users/:id:
    get:
      summary: Get user by id.
      description: This API is used to get user info by ID.
      operationId: getUserById
      responses:
        '200':
          description: Successfully get user info.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/users'
        '404':
            description: Not Found.
            content:
                application/json:
                schema:
                    $ref: '#/components/schemas/error'
        '500':
          description: Internal Server Error.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/error'
    delete:
      ...
    put:
      ...
  /api/users:
    post:
      ...
```

### 1.2. チューニングせずにChatGPTに雑に投げる

```chatgpt
あなたはバックエンドエンジニアです。
今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

## 制約条件
- APIは３層アーキテクチャに従って作成すること
- APIはGo言語で作成すること

## 入力ファイル
- DBのスキーマファイル:
    ```
    ここにスキーマファイルを貼る
    ```

- APIの定義書:
    ```
    ここにAPIの定義書を貼る
    ```

## 出力ファイル
- アーキテクチャに従って分割されたAPIのコード
```

- [結果](https://chat.openai.com/share/09469c0c-7d1a-4a51-a2be-272881eb3e28)
  - DBの初期化やRoutingの設定などもすべて生成して出力してしまう
  - フレームワークやDBエンジンも指定していないので、よくあるものを選択して出力している（Gin, mysql）
  - サービス層に関しては「ここにロジックを実装」のようなコメントが生成され、実装はされていない
  - 既存のアーキテクチャに組み込むためには、どのAPIを生成し、どのようなアーキテクチャかを詳細に指定する必要がある気がする

### 1.3. 指示を細かくしてどこまでいけるかを確認

```chatgpt
あなたはバックエンドエンジニアです。
今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

## 制約条件
- APIは「Controller層」、「Service層」、「Repository層」の３層アーキテクチャに従って作成すること
- APIはGo言語で作成すること
- フレームワークはEchoを使用すること
- DBエンジンはMySQLを使用すること
- ORMマッパーはGORMを使用すること
- ユーザーを生成するAPIのみを作成すること

## 入力ファイル
- DBのスキーマファイル:
    ```
    ここにスキーマファイルを貼る
    ```

- APIの定義書: api_1.yaml
    ```
    ここにAPIの定義書を貼る
    ```

## 出力ファイル
- アーキテクチャに従って分割されたAPIのコード
```

- [結果](https://chat.openai.com/share/b9aeeff6-8d9e-45d5-a18a-6f796930ace0)
  - 他のAPIやRoutingの設定も、作成されてしまったが、コピペだけですむコードができた
  - 既存のプロジェクトのコードに組み込むとなると、既存のコードをある程度食わせないと（この程度のAPIでも）難しそう

### 1.4. 既存プロジェクトのコードを食わせて、どこまでいけるかを確認

`api/users/:id`の`GET`に対応するコードをこちらで実装し、それをChatGPTに食わせて`POST`に対応するコードを生成してもらう。

#### 1.4.1. Custom Instructionを使って、既存プロジェクトのコードを食わせる

既存プロジェクトのコードをCustom Instructionとして、改めて指示を出してみる。

`How would you like ChatGPT to respond?`の部分に指定してみる。

- Custom Instruction(文字数、Token数に制限あり)

    ```chatgpt
    APIを実装する際、アーキテクチャは以下のように層を分離してファイルを出力する。

        ```Controller
        controller層のAPIの実装を例として貼り付ける。
        ```
        ```service
        service層のAPIの実装を例として貼り付ける。
        ```
    repository層とDBのmodelも同様に分割する。
    レスポンスやリクエストのスキーマが必要な場合にも適宜ファイルを分割する。
    ```

- Prompt

    ```chatgpt
    あなたはバックエンドエンジニアです。
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - ユーザーを生成するAPI（operationId: createUser）のみを作成すること

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - アーキテクチャに従って分割されたAPIのコード
    ```

- [結果](https://chat.openai.com/share/c4a40867-c001-4914-a826-ad0fe04c4639)
  - 何度かやりとりしたらよいコードができた。
  - やはりmain関数とDBの初期化コードを作りたがる。-> 事前の制約条件に明示的に指定した方が良い
  - errorのレスポンスがスキーマ通りにならない
  - Custom instructionsでコードを渡すよりも、Promptでコードを渡した方が良いかもしれない

#### 1.4.2. Promptを使って、既存プロジェクトのコードを食わせる

Custom insturctionsだと文字数の制限がよりタイトなので、Promptでコードを渡すようにしてみる。

- Prompt1

    ```chatgpt
    あなたはバックエンドエンジニアで、私はプロジェクトマネージャーです。
    今後私からAPIを実装するように頼まれた際には、以下に送るようなスタイルで実装してください。
    ポイントは
    - アーキテクチャの各層に応じてファイルを分離して出力すること。
    - DBの初期化やルーティングは私が実装するので、あなたが実装する必要はなく、そのAPIによって新たに必要となるコードのみを実装すること。
    - 必要に応じて、以下のスタイルに含まれないが必要な要素（RequestやResponseのSchemaの定義など）は適宜ファイルを分割して実装すること。

      ```Controller
      controller層のAPIの実装を例として貼り付ける。
      ```
      ```service
      service層のAPIの実装を例として貼り付ける。
      ```
      ```repository
      repository層のAPIの実装を例として貼り付ける。
      ```
      ```model
      modelの実装を例として貼り付ける。
      ```

    以上の内容が理解出来たら"Yes"と返事をし、不明点がある場合にはその内容を具体的に明らかにして質問してください。
    ```

- Prompt2

    ```chatgpt
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - ユーザーを生成するAPI（operationId: createUser）のみを作成すること

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - アーキテクチャに従って分割されたAPIのコード
    ```

- [結果](https://chat.openai.com/share/dba93de0-1f1b-4f79-a999-d0cde4c62aac)
  - **対話なしで一番よいコード**
  - 相変わらずDBの初期化とルーティングのコードを生成するので、都度作る部分を明確にしないとダメそう。
  - errorのレスポンスがスキーマ通りにならない
  - Prompt1がCLIにした際の`COMMAND learn ...`的なイメージ
  - Prompt2がCLIにした際の`COMMAND gen ...`的なイメージ

これを踏まえて、Prompt2を改良してコードを生成する。

- Prompt2 改

    ```chatgpt
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - errorハンドリングは`echo.NewHTTPError`を使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - ユーザーを生成するAPI（operationId: createUser）のみを作成すること
    - DBの接続、ルーティングのコードを生成しないこと

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - アーキテクチャに従って分割されたAPIのコード
    ```

- [結果](https://chat.openai.com/share/dbe7c9b0-f6fe-47ae-b0bf-3036a7d7dcb5)
  - errorハンドリングはクリアしたが、DBの接続コードを生成しないようにすると、今度はDB周りのコードが生成されなくなってしまった。

さらに改良

- Prompt2 改改

    ```chatgpt
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを作成してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - errorハンドリングは`echo.NewHTTPError`を使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - ユーザーを生成するAPI（operationId: createUser）のみを作成すること
    - APIのコードをを生成する際には、事前に与えたスタイル通りにファイルを分割し、必要なファイルはすべて生成すること
    - DBの接続、ルーティングのコードを生成しないこと、ただし、DBには接続できていて、ルーティングも正しく設定されているものとしてコードを生成すること

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - アーキテクチャに従って分割されたAPIのコード
    ```

- [結果](https://chat.openai.com/share/d46880f5-2d20-4115-801e-ea000e9f856f)
  - ようやく理想通りのコードが生成された
  - すこし制約条件が多すぎるかもしれない

## 2. チューニングしたChatGPTを用いて、複数のテーブルが絡むようなAPIの作成

上の結果で一番よかった方針を、より複雑なAPIに適用してみる。

`users`, `notifications`, `users_notifications`の３つのテーブルを用意。ユーザーの既読がついていない通知のみを取得するAPIを作成する。

### 2.1. Input

- DB Schema (`schema_2.sql`)

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `mail` varchar(100) NOT NULL,
  `birthday` DATE NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL,
  `message` varchar(100) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
);

CREATE TABLE `users_notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `is_read` boolean NOT NULL DEFAULT false,
  `notification_id` bigint unsigned NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `users_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `users_notifications_notification_id` FOREIGN KEY (`notification_id`) REFERENCES `notifications` (`id`)
);
```

- API Spec (`api_2.yaml`)

```yaml
OpenAPI: 3.0.3
info:
  title: Notification API
  version: 1.0.0
paths:
  /api/users/notifications:
    get:
      summary: Get not read notifications
      description: This endpoint returns notifications that have not been read by the user.
      operationId: getNotReadNotifications
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/user_id'
      parameters:
        - in: query
          name: offset
          schema:
            type: integer
            minimum: 0
            default: 0
          required: false
          description: page offset
        - in: query
          name: limit
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 10
          required: false
          description: page size
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/notifications'
        '400':
          description: Bad Request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/error'
        '500':
          description: Internal Server Error.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/error'
```

*懸念点*

- 入力のValidation
- repository層の実装
  - Joinして一気にとるのか、user_id -> notification_id, notification_id -> notificationの順に取るのか

### 2.2. 同様のチューニング方針でコードを生成

先ほどと同様に、Prompt1として、UserをId指定して`GET`するコードを実装し、それをChatGPTに食わせて今回のAPIを生成してもらう。

- Prompt1

    ```chatgpt
    あなたはバックエンドエンジニアで、私はプロジェクトマネージャーです。
    今後私からAPIを実装するように頼まれた際には、以下に送るようなスタイルで実装してください。
    ポイントは
    - アーキテクチャの各層に応じてファイルを分離して出力すること。
    - DBの初期化やルーティングは私が実装するので、あなたが実装する必要はなく、実装を依頼したAPIによって新たに必要となるコードのみを実装すること。
    - APIのロジックもすべて実装した状態でファイルを出力すること。
    - 必要に応じて、以下のスタイルに含まれないが必要な要素（RequestやResponseのSchemaの定義など）は適宜ファイルを分割して実装すること。

      ```Controller
      controller層のAPIの実装を例として貼り付ける。
      ```
      ```service
      service層のAPIの実装を例として貼り付ける。
      ```
      ```repository
      repository層のAPIの実装を例として貼り付ける。
      ```
      ```model
      modelの実装を例として貼り付ける。
      ```

    以上の内容が理解出来たら"Yes"と返事をし、不明点がある場合にはその内容を具体的に明らかにして質問してください。
    ```

- Prompt2

    ```chatgpt
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書からAPIを実装してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - errorハンドリングは`echo.NewHTTPError`を使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - APIのコードをを生成する際には、事前に与えたスタイル通りにファイルを分割し、必要なファイルはすべて生成し、中身のロジックもすべて実装すること。
    - DBの接続、ルーティングのコードを生成しないこと、ただし、DBには接続できていて、ルーティングも正しく設定されているものとしてコードを生成すること

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - アーキテクチャに従って分割されたAPIのコード
    ```

- [結果](https://chat.openai.com/share/5ffb6874-66ae-40c6-8d46-858348d9ba84)
  - がわだけ実装して中身のロジックは実装してくれない（`中身のロジックもすべて実装すること`という制約条件を足してもダメだった。）
  - 「ロジックも実装してください」と続けると、書いてくれるが、ロジックが間違っている。レスポンスのスキーマも定義通りではない。

### 2.3. 各層を別々に出力させる

一回ですべてを出力させるのではなく、各層を別々に出力させることで、安定してコードを生成できるかを確認する。

- Prompt2

    ```chatgpt
    今から、以下の制約条件を満たすように、DBのスキーマファイルとAPIの定義書から、APIの「Controller層」を実装してください。

    ## 制約条件
    - APIはGo言語で作成すること
    - フレームワークはEchoを使用すること
    - errorハンドリングは`echo.NewHTTPError`を使用すること
    - DBエンジンはMySQLを使用すること
    - ORMマッパーはGORMを使用すること
    - APIのコードをを生成する際には、事前に与えたスタイル通りにファイルを分割し、中身のロジックもすべて実装すること。
    - DBの接続、ルーティングのコードを生成しないこと、ただし、DBには接続できていて、ルーティングも正しく設定されているものとしてコードを生成すること
    - 依存関係のある他のファイルはすでに実装されているものとしてコードを生成すること。

    ## 入力ファイル
    - DBのスキーマファイル:
        ```
        ここにスキーマファイルを貼る
        ```

    - APIの定義書:
        ```
        ここにAPIの定義書を貼る
        ```

    ## 出力ファイル
    - 「Controller層」のコード
    ```

  これを３回繰り返して、各層を。実装させる。

- [結果](https://chat.openai.com/share/cca33cb0-93f2-40c6-9d4b-33dff26f6973)
  - それぞれのロジックはそれっぽいものが生成されたが、、、
    - Schemaの定義が違う
    - テーブルの構造（中間テーブル）を理解していないため、Repository層のロジックがめちゃくちゃ

テーブルをJoinするよう指示してみる。

- Prompt2

  ```chatgpt
    ## 制約条件
    - users_notifications tableとnotifications tableをjoinしてリソースを取得してください。
  ```

- [結果](https://chat.openai.com/share/63477ea3-3d4d-4d29-94ae-1b41f2a22ec3)
  - 正しくRepository層が実装できた
  - 学習したコードにないResponseのSchema定義等は生成されず。

## まとめ

- 簡単なCRUD操作であれば、指示を細かくしていけば、よいコードが生成できる
- 複雑なAPIは、各層を別々に出力させることで、よいコードが生成できるかもしれない
  - ただし、現状だと制約条件が多すぎる
  - Projectにフィットさせるためにはもっと学習させるコードを増やさないとダメかもしれない
  - tableのJoinなどの複雑なロジックやリレーションを把握させるには、もっと別な方法が良いかもしれない（あまりアイデアはない）
