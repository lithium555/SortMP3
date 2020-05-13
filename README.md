#   SortMP3

#   Up postgres in container
```
docker run --name googleAPI_new -e POSTGRES_PASSWORD=rem -p 5432:5432 -d postgres

docker run --name sort_music -e POSTGRES_PASSWORD=master -e POSTGRES_DB=musicDB -e POSTGRES_USER=sorter -p 5432:5432 -d postgres

53aa124bb2af  postgres    "docker-entrypoint.s…"  About a minute ago   Up About a minute   0.0.0.0:5432->5432/tcp   googleAPI_new

```

#   Remove postgres-container
```docker rm -f sort_music```
#   Draft
PostgreSQL:
  restart: always
  image: sameersbn/postgresql:10-2
  ports:
    - "5432:5432"
  environment:
    - DEBUG=false

    - DB_USER=sorter
    - DB_PASS=master
    - DB_NAME=musicDB
    - DB_TEMPLATE=

    - DB_EXTENSION=

    - REPLICATION_MODE=
    - REPLICATION_USER=
    - REPLICATION_PASS=
    - REPLICATION_SSLMODE=
  volumes:
    - /srv/docker/postgresql:/var/lib/postgresql