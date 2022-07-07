# paaticle
A simple PaaS service

## 動作シーケンス (予定)
- Gitリポジトリをローカルにクローン
- 指定されたブランチにチェックアウト
- `.paaticlegnore`で指定されていないファイルをtarにまとめてコンテナへ
- `paaticle.yml`で指定 or 各言語デフォルトのビルドコマンドを実行
- 指定されている場合のみビルド後フックを実行
- コードを実行

