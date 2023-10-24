# 翻訳サンプルアプリ

スプレッドシートから日本語文字列を取得し英字の変数名に変換後、
スプレッドシートに変換後の変数名を保存します。

## 実行方法
1. [Cloud Console](https://console.cloud.google.com/)にログインする
2. プロジェクトを作成する
3. 以下のAPIを許可する
    - Google Sheets API
    - Google Drive API
    - Cloud Translation API
4. ターミナルを開き、デフォルトの認証情報を作成する。`<project ID>`には#2で作成したプロジェクトのIDを指定する。
```
$ gcloud auth application-default set-quota-project <project ID>
```
5. アプリを実行する
```
$ cd translate_sample/
$ go run .
```
