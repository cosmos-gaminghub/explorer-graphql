# GAME Explorer GraphQL API

**Prerequisites**
* go1.16.0+
* mongoDB 4.4.0+

## Get Started
```
git clone https://github.com/cosmos-gaminghub/explorer-graphql
cd explorer-graphql
cp .env.example .env
make all
./build/graphql
```

### Use Service
```sh
### service
tee /etc/systemd/system/explorer-graphql.service > /dev/null <<EOF
[Unit]
Description=exploder-graphql

[Service]
Type=simple
User=root
Group=root

WorkingDirectory=/root/explorer-graphql
ExecStart=/root/explorer-graphql/build/graphql

Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable explorer-graphql
systemctl start explorer-graphql
```

## Setting environment variables

Env Variables| Description | Default
------------------ | ---------------------------- | --------
DB_ADDR            | mongo server's address       | 127.0.0.1
DB_DATABASE        | database's name              | test
DB_USER            | database's username          |
DB_PASSWORD        | database's password          |
DB_POOL_LOMIT      | database max connection num  | 4096
PORT               | explorer server's port       | 8080
LCD                | node light client daemon URI | http://198.13.33.206:1317

Access to `localhost:8080`, then you can refer to the graphql schema docs in graphql playground. `localhost:8080/query` will be the graphql endpoint which should be set in a frontend env file.
