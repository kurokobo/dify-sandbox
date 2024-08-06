# 使い方

🌏 [**日本語**](./usage.ja.md)
🌏 [**English**](./usage.md)

> [!CAUTION]
>
> - **このコンテナイメージを使用すると、あなたの環境は悪意のあるコードからの攻撃に非常に脆弱になります。**
> - **自分が行おうとしていることを自分で充分に理解できていない場合は、このコンテナイメージは使用しないでください。**
> - **本番環境では使用しないでください。**
> - **複数人で共有している Dify 環境では使用しないでください。**
> - **自分だけが使用できる環境で、実験とテストの目的でのみ使用してください。**

`docker-compose.yaml` ファイル内で、`sandbox` サービスの `image` を置き換えます。

```yaml
...
  sandbox:
    image: ghcr.io/kurokobo/dify-unsandboxed-sandbox:0.2.4-unsandboxed
...
```

その後、`sandbox` サービスを再起動します。

```bash
docker compose down sandbox
docker compose pull sandbox
docker compose up -d sandbox
```

## 追加機能

サンドボックス内に `apt` でインストールしたいパッケージがある場合、`./volumes/sandbox/dependencies/apt-requirements.txt` に列挙することで、コンテナの起動時にインストールされます。
