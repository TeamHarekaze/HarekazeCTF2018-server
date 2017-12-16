# HarekazeCTF-Competition

## Directory structure
```
└─web
    ├─controllers
    ├─models
    ├─public (static files)
    └─views
```

## 開発環境構築方法
- `golangのインストール済みであること`
- `mysql, redisにアクセスできる環境であること`
  - [mysql, redisの簡単構築](https://github.com/HayatoDoi/HarekazeCTF-Competition/tree/master/src#mysql-redis%E3%81%AE%E7%B0%A1%E5%8D%98%E6%A7%8B%E7%AF%89)
1. クローン
```shell
git clone https://github.com/HayatoDoi/HarekazeCTF-Competition.git
```

2. 作業ディレクトリに移動
```shell
cd HarekazeCTF-Competition/src
```

3. ライブラリのインストール
```
sh lib_install.sh 
```

4. 設定ファイルの作成
- `各自の環境に合わせて変更してください`
```shell
cp .env.example .env
vim .env
```

5. サーバの起動
```shell
go run harekazectf.go
```

# mysql, redisの簡単構築
`Docker`でやるよ！

1. 適当な作業ディレクトリを作成 & 移動
```shell
mkdir -p ~/docker/HarekazeCTF_DB
cd ~/docker/HarekazeCTF_DB
```

2. コンテナ設定ファイルをダウンロード
```shell
git clone https://github.com/HayatoDoi/docker-mysql.git
```

3. 起動スクリプトを作成
```shell
vim start.sh
```

  - ファイルの中身は次の通り
```
cd docker-mysql/
docker-compose up -d
docker run --rm --name redis -d -p 6379:6379 redis redis-server --appendonly yes
```

4. 起動
```shell
sh start.sh
```