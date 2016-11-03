#!/bin/sh
docker run -dit --name roach1  -p 26257:26257 -p 8080:8080 cockroachdb/cockroach:beta-20161027  start --insecure
