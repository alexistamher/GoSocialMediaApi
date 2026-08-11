FROM postgres:18.1-alpine

COPY social-media-db.sql /docker-entrypoint-initdb.d/

