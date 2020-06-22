# Markdown to PDFコアサービス

MarkdownをPDFに変換するサービス。

## 仕様

HTTPでMarkdownを受け取り、PDFに変換して返す。

### リクエスト

* メソッドはPOST
* ボディは以下のどちらか
  * Markdown
    * ヘッダのContent-Typeはtext/markdownにする
  * tar
    * ヘッダのContent-Typeはapplication/x-tarにする
* Multipartは受け付けない

### レスポンス

正常時、リクエストに含まれるMarkdownの数によって以下のどちらかを返す。

* Markdownが1つのとき、そのファイルから変換したPDFを返す
  * ヘッダのContent-Typeはapplication/pdf
* Markdownが複数のとき、それらから変換したPDFをまとめたtarを返す
  * ヘッダのContent-Typeはapplication/x-tar

### 挙動

#### 複数のMarkdown

リクエストのtarに含まれる複数のMarkdownは、PDFに変換された後もその相対パスを保つようにレスポンスのtarに入れられる。

#### 画像の埋め込み

変換したPDFには画像が埋め込まれる。
ローカルの画像を利用する場合は、リクエストのtarに画像も入れ、Markdownからその相対パスを参照するようにする。

## 使い方

以下はローカル実行時のURLを使用。
URLは適切に置き換えること。

```console
$ curl http://localhost:8080/ -H 'Content-Type: text/markdown' -d '@test.md' -o test.pdf
```

```console
$ tar cf req.tar test-dir/
$ curl http://localhost:8080/ -H 'Content-Type: application/x-tar' -d '@req.tar' -o res.tar
```

## 開発

### ローカル実行

```console
$ go run invoke.go
```

#### Dockerで

```console
$ docker build --tag markdown-to-pdf:latest .
$ docker run -p 8080:8080 markdown-to-pdf:latest
```

### デプロイ

```console
$ gcloud builds submit --tag gcr.io/<PROJECT-NAME>/markdown-to-pdf
$ gcloud run deploy --image gcr.io/<PROJECT-NAME>/markdown-to-pdf --platform managed
```
