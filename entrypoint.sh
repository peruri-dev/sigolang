#!/bin/sh

/app db:init
/app db:status
/app db:migrate

exec /app "$@"