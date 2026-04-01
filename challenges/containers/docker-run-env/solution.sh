#!/bin/bash
docker run -d --name my-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_DATABASE=testdb mysql
