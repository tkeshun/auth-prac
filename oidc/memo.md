https://openid.net/

## OpenIDについて

https://openid.net/foundation/

## OpenID Connect

OAuth 2.0 仕様フレームワーク (IETF RFC 6749 および 6750) に基づく相互運用可能な認証プロトコル
認可サーバーによって実行される認証に基づいてユーザーのIDを検証する

パスワードの設定、保存、管理の責任をアプリケーション開発者が持つ必要がなくなり、第三者に委託できる→※自分でIdP建てるときは対象外？

### OAuth2.0とOpenID Connectとの関係

OAuth2.0 → RFC6749および6750で定義された認証及び承認プロトコルの開発をサポートするように設計されたフレームワーク  

Auth屋さんの本より、「OAuthは認可のフレームワークであり、認証のフレームワークでないことが述べられてた」はず  

OAtuh2.0を拡張して、OIDCは認証を行う。  


### OpenID ConnectとSAMLの関係

SAMLは一部の企業や学術分野のユースケースで使用されるXMLベースの認証  
OpenID Connectは、よりシンプルなJSON/RESTベースのプロトコルを使用して、同じユースケースに対応可能  

# 本とか資料

- OAuth、OAuth認証、OpenID Connectの違いを整理して理解できる本
- 雰囲気でOAuth2.0を使っているエンジニアがOAuth2.0を整理して、手を動かしながら学べる本


