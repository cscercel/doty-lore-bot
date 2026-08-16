#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/migrations
goose turso $DATABASE_URL up
