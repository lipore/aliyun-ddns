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

Run in EdgeOS with ppp ip-up.d 
--------

### Build
```shell
GOOS=linux GOARCH=mipsle go build -o aliyun-ddns-edgeos edgeos/main.go
```

### RUN

1. Copy aliyun-ddns-edgeos to edgeRouter
2. create config.yml in edgeRouter
3. add run script to /config/scripts/ppp/ip-up.d/ddns.sh
```shell
/usr/bin/aliyun-ddns-edgeos $1 $4
```

Run in EdgeOS with task-scheduler
--------

### Build
```shell
GOOS=linux GOARCH=mipsle go build -o aliyun-ddns-edgeos main.go
```

### Setup task

1. Copy aliyun-ddns-edgeos to edgeRouter in /config/scripts/
2. create config.yml in edgeRouter
3. create task 
```
set system task-scheduler task ddns interval 20m
set system task-scheduler task ddns executable path /config/scripts/aliyun-ddns-edgeos
set system task-scheduler task ddns executable arguments /config/scripts/ddns.config.yml
commit
save
```
