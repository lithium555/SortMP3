#Description of migration commands

* install [migrate] command: (https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```console
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

* creating migrations files:

```console
make create-migrations n=your_name_of_migration
```

### Before we use next commands, we need set `database` variable in Makefile.

#### Set your params for connects to postgres DB.

* rolling migrations:

```console
make migrate-up
```

* rolling back migrations:
  (If you specify flag `s=i` this will rollback `i` migrations.)
```console
make migrate-down s=i
```

* dropping migrations:

```console
make migrate-drop
```