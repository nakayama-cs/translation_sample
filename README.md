# 翻訳サンプルアプリ

スプレッドシートから日本語文字列を取得し英字の変数名に変換後、
スプレッドシートに変換後の変数名を保存します。

## 実行方法
1. 秘密鍵の.jsonを環境変数に格納する
```
$ export GOOGLE_APPLICATION_CREDENTIALS=keys/translationsample-402703-1d82b934413a.json
```
2. アプリを実行する
```
$ cd translate_sample/
$ go run .
```
