
## 起動コマンド

`command: 'start-from-init --masterkey "MasterkeyNeedsToHave32Characters" --tlsMode disabled'`

start-from-init:
ZITADEL を初回起動モードでセットアップする。
--masterkey "MasterkeyNeedsToHave32Characters"
マスターキー (暗号化のため) は 32文字 必要。
これが適切に設定されていないと、起動時にエラーになる。
--tlsMode disabled
TLS（SSL）を無効化（開発環境向け）。
本番環境では TLS を有効化するべき。

## 環境変数

- ZITADEL_DATABASE_POSTGRES_HOST: db
データベースのホスト名。db サービスを参照。
- ZITADEL_DATABASE_POSTGRES_PORT: 5432
PostgreSQL のポート番号。
- ZITADEL_DATABASE_POSTGRES_DATABASE: zitadel
ZITADEL のデータベース名。
- ZITADEL_DATABASE_POSTGRES_USER_USERNAME: zitadel
ZITADEL が使用するデータベースユーザー。
- ZITADEL_DATABASE_POSTGRES_USER_PASSWORD: zitadel
上記ユーザーのパスワード。
- ZITADEL_DATABASE_POSTGRES_USER_SSL_MODE: disable
SSL を無効化（開発環境向け）。
- ZITADEL_DATABASE_POSTGRES_ADMIN_USERNAME: postgres
管理者（postgres）のユーザー名。
- ZITADEL_DATABASE_POSTGRES_ADMIN_PASSWORD: postgres
上記ユーザーのパスワード。
- ZITADEL_DATABASE_POSTGRES_ADMIN_SSL_MODE: disable
管理者接続でも SSL を無効化（開発環境向け）。
- ZITADEL_EXTERNALSECURE: false
外部アクセスのセキュリティをオフ（開発環境向け）。


## メモ

初回起動だとマイグレーションが走ってるっぽい
localhost:8080にアクセスするとユーザー登録画面にいくが、公式docの設定では、SMTP関連の設定がないので、認証コードをおくれない。
メールを送れるようにするには、ZITADEL側にSMTPの設定をするのと、外部のメールサーバーが必要

