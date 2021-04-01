Aliyun-DDNS
============

This project use ipify.org to determine the external ip. And then call aliyun api to update the dns A record or AAAA record

requirement
-----------

1. setup aliyun account
2. create an rebot RAM account, with app key and app secret. The RAM must have dns privilege.
3. change you domain's ns record to ns1.alidns.com and ns2.alidns.com
4. add you domain to aliyun dns

Run in Docker
------

### Build image
```
docker build -t aliyun-ddns:0.1 .
```

### Configure
1. make a folder to store the config file
2. copy config.yml to your config folder
3. update the config file

### Run
```
docker run -d -v /your/config/folder:/opt/aliyun-ddns/config aliyun-ddns:0.1
```

