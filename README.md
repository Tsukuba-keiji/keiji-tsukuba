# keiji-tsukuba

情報の授業で作成した情報統括システムのソースコード。

## main.go

LINEのコールバック設定（https://github.com/line/line-bot-sdk-go のkitchensinkを利用）、サーバーの設定を行う。

## /static

LINEで使う画像ファイルを置く。

## /src

jsonの一時保存場所。

## /assets

ホームページに表示するHTMLの置き場。

### /assets/

ルート。検索フォームがある。

### /assets/search/

検索結果をGETメソッドで表示する。

### /assets/add/

index.htmlで入力し、verify.htmlで確認・POST送信する。

## GAS（外部API）

Google Driveにjsonをバックアップする。
