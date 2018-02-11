# HarekazeCTF2018-server

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
  - [mysql, redisの簡単構築](https://github.com/TeamHarekaze/HarekazeCTF2018-server/tree/master/src#mysql-redis%E3%81%AE%E7%B0%A1%E5%8D%98%E6%A7%8B%E7%AF%89)
1. クローン
```shell
git clone https://github.com/TeamHarekaze/HarekazeCTF2018-server.git
```

2. 作業ディレクトリに移動
```shell
cd HarekazeCTF2018-server/src
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

## Docker Composeを利用する場合
Docker Composeを利用すると手元の環境を汚さずにDockerで構成された開発環境を利用することが出来ます。

```console
% docker-compose build
% cat src/.env
APP_PORT=5000
APP_ADMIN_HASH="bda130d51a9b2e2f7f7929831af92e43"

# Competition Time
COMPETITION_START_TIME="2018-02-10 15:00:00" #UTC +09:00
COMPETITION_END_TIME="2018-02-11 15:00:00" #UTC +09:00

## MySQL
DB_HOST="db"
DB_NAME="HarekazeCTF"
DB_PORT=3306
DB_USER="root"
DB_PASSWORD="password"

## Redis
REDIS_HOST="redis"
% docker-compose up -d
% docker-compose run --rm app mysql -uroot -ppassword -hdb -e "CREATE DATABASE HarekazeCTF;"
% docker-compose run --rm app mysql -uroot -ppassword -hdb HarekazeCTF < src/migrate.sql
% docker-compose restart app
```
