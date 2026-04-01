#!/bin/bash
docker run -d --name lifecycle-test nginx
docker stop lifecycle-test
docker start lifecycle-test
