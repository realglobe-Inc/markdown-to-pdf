# Cloud Runの使い方

```console
$ gcloud builds submit --tag gcr.io/<PROJECT-ID>/helloworld
$ gcloud run deploy --image gcr.io/<PROJECT-ID>/helloworld --platform managed
```

## 構想

### Markdown to PDF

* これはできそう

Multipartで.mdファイルや画像を受け取り、.mdファイルが1つならPDFファイルを返す。
.mdファイルが複数なら複数のPDFをtarで固めたファイルを返す。

### Google Driveコネクタ

Google Driveアプリ。

* これはできるか未検証

Driveの選択されたファイルをMultipartに入れてMarkdown to PDFに投げて、返ってきたPDFをDriveの元の.mdファイルの場所に置く。
