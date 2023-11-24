# GoRestify Backend Template


#### Redis and MySQL with docker
```bash 
docker run --rm --name db-mysql -d -v mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=88888888 -e TZ='Asia/Baghdad' -p 3306:3306 mysql --innodb_lock_wait_timeout=1000 --innodb_buffer_pool_size=2147483648  --max_allowed_packet=1073741824  --max_connections=2000 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run --rm --name db-redis -d -v redis-data:/usr/local/etc/redis -p 6379:6379 redis:latest
```

#### Create Travis User
```SQL 
create user 'travis'@'%'; Grant all privileges on *.* To 'travis'@'%' with grant option;
```

### Run Project
``` bash
# run admin
source cmd/admin/sample.env && reflex -r '\.go' -s -- sh -c 'go run cmd/admin/main.admin.go'

# run user
source cmd/user/sample.env && reflex -r '\.go' -s -- sh -c 'go run cmd/user/main.user.go'
```