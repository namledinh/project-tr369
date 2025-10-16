#!/bin/sh
export APP_NAME=usp-management-backend-dev
export APP_PORT=2105
export APP_LOG_LEVEL=debug
export SQL_HOST=1.52.246.136
export SQL_PORT=15421
export SQL_USP_USER=
export SQL_USP_PASS=
export SQL_USP_DB=usp_system_v2
export MINIO_ENDPOINT=s3-api.dc11.inf.fpt.net
export MINIO_ACCESS_KEY=1b8xAQl4JRXixJoXgmQn
export MINIO_SECRET_KEY=RqiQUNFOKUtNMiuvQkD9F3j7Sc6uG57wFmtduBAt
cd /app/ && ./main
