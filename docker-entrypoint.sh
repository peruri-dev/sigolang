#!/bin/sh

/api db:init
/api db:status
/api db:migrate

exec /api "$@" 2>&1 | tee /var/log/api.log
