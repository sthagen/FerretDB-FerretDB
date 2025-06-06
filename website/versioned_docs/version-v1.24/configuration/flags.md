---
sidebar_position: 1
---

# Configuration flags

FerretDB provides numerous configuration flags you can customize to suit your needs and environment.
You can always see the complete list by using `--help` flag.
To make user experience cloud native, every flag has its environment variable equivalent.
There is no configuration file.

:::info
Some default values are overridden in [our Docker image](../quickstart-guide/docker.md).
:::

<!-- Keep order in sync with the `--help` output -->

## General

| Flag              | Description                                                       | Environment Variable     | Default Value                  |
| ----------------- | ----------------------------------------------------------------- | ------------------------ | ------------------------------ |
| `-h`, `--help`    | Show context-sensitive help                                       |                          | false                          |
| `--version`       | Print version to stdout and exit                                  |                          | false                          |
| `--handler`       | Backend handler                                                   | `FERRETDB_HANDLER`       | `pg` (PostgreSQL)              |
| `--mode`          | [Operation mode](operation-modes.md)                              | `FERRETDB_MODE`          | `normal`                       |
| `--state-dir`     | Path to the FerretDB state directory<br />(set to `-` to disable) | `FERRETDB_STATE_DIR`     | `.`<br />(`/state` for Docker) |
| `--repl-set-name` | Replica set name<br />(should be set for OpLog to work correctly) | `FERRETDB_REPL_SET_NAME` | empty                          |

## Interfaces

| Flag                     | Description                                                                               | Environment Variable            | Default Value                                |
| ------------------------ | ----------------------------------------------------------------------------------------- | ------------------------------- | -------------------------------------------- |
| `--listen-addr`          | Listen TCP address                                                                        | `FERRETDB_LISTEN_ADDR`          | `127.0.0.1:27017`<br />(`:27017` for Docker) |
| `--listen-unix`          | Listen Unix domain socket path                                                            | `FERRETDB_LISTEN_UNIX`          |                                              |
| `--listen-tls`           | Listen TLS address (see [here](../security/tls-connections.md))                           | `FERRETDB_LISTEN_TLS`           |                                              |
| `--listen-tls-cert-file` | TLS cert file path                                                                        | `FERRETDB_LISTEN_TLS_CERT_FILE` |                                              |
| `--listen-tls-key-file`  | TLS key file path                                                                         | `FERRETDB_LISTEN_TLS_KEY_FILE`  |                                              |
| `--listen-tls-ca-file`   | TLS CA file path                                                                          | `FERRETDB_LISTEN_TLS_CA_FILE`   |                                              |
| `--proxy-addr`           | Proxy address                                                                             | `FERRETDB_PROXY_ADDR`           |                                              |
| `--proxy-tls-cert-file`  | Proxy TLS cert file path                                                                  | `FERRETDB_PROXY_TLS_CERT_FILE`  |                                              |
| `--proxy-tls-key-file`   | Proxy TLS key file path                                                                   | `FERRETDB_PROXY_TLS_KEY_FILE`   |                                              |
| `--proxy-tls-ca-file`    | Proxy TLS CA file path                                                                    | `FERRETDB_PROXY_TLS_CA_FILE`    |                                              |
| `--debug-addr`           | Listen address for HTTP handlers for metrics, profiling, etc<br />(set to `-` to disable) | `FERRETDB_DEBUG_ADDR`           | `127.0.0.1:8088`<br />(`:8088` for Docker)   |

## Backend handlers

<!-- Do not document alpha backends -->

### PostgreSQL

[PostgreSQL backend](../understanding-ferretdb.md#postgresql) can be enabled by
`--handler=pg` flag or `FERRETDB_HANDLER=pg` environment variable.

| Flag               | Description                     | Environment Variable      | Default Value                        |
| ------------------ | ------------------------------- | ------------------------- | ------------------------------------ |
| `--postgresql-url` | PostgreSQL URL for 'pg' handler | `FERRETDB_POSTGRESQL_URL` | `postgres://127.0.0.1:5432/ferretdb` |

FerretDB uses [pgx v5](https://github.com/jackc/pgx) library for connecting to PostgreSQL.
Supported URL parameters are documented there:

- https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#ParseConfig
- https://pkg.go.dev/github.com/jackc/pgx/v5#ParseConfig
- https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#ParseConfig

Additionally:

- `pool_max_conns` parameter is set to 50 if it is unset in the URL;
- `application_name` is always set to "FerretDB";
- `timezone` is always set to "UTC".

### SQLite

[SQLite backend](../understanding-ferretdb.md#sqlite) can be enabled by
`--handler=sqlite` flag or `FERRETDB_HANDLER=sqlite` environment variable.

| Flag           | Description                                 | Environment Variable  | Default Value                                     |
| -------------- | ------------------------------------------- | --------------------- | ------------------------------------------------- |
| `--sqlite-url` | SQLite URI (directory) for 'sqlite' handler | `FERRETDB_SQLITE_URL` | `file:data/` `.`<br />(`file:/state/` for Docker) |

FerretDB uses [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) library for accessing SQLite database files.
Supported URL parameters are documented there:

- https://www.sqlite.org/uri.html
- https://pkg.go.dev/modernc.org/sqlite#Driver.Open
- https://www.sqlite.org/pragma.html

Additionally:

- `_pragma=auto_vacuum(none)` parameter is set if that PRAGMA is not present;
- `_pragma=busy_timeout(10000)` parameter is set if that PRAGMA is not present;
- `_pragma=journal_mode(wal)` parameter is set if that PRAGMA is not present.

One difference is that URI should point to the existing directory (with absolute or relative path), not to a single database file.
That allows FerretDB to work with multiple databases.

In-memory SQLite databases are fully supported.
In that case, the URI should still point to the existing directory (that will be unused).
For example: `file:./?mode=memory`.

## Miscellaneous

| Flag                     | Description                                                                     | Environment Variable            | Default Value    |
| ------------------------ | ------------------------------------------------------------------------------- | ------------------------------- | ---------------- |
| `--log-level`            | Log level: 'debug', 'info', 'warn', 'error'                                     | `FERRETDB_LOG_LEVEL`            | `info`           |
| `--[no-]log-uuid`        | Add instance UUID to all log messages                                           | `FERRETDB_LOG_UUID`             |                  |
| `--[no-]metrics-uuid`    | Add instance UUID to all metrics                                                | `FERRETDB_METRICS_UUID`         |                  |
| `--otel-traces-url`      | OpenTelemetry OTLP/HTTP traces endpoint URL (e.g. `http://host:4318/v1/traces`) | `FERRETDB_OTEL_TRACES_URL`      | empty (disabled) |
| `--test-enable-new-auth` | Enable new authentication mode                                                  | `FERRETDB_TEST_ENABLE_NEW_AUTH` | false            |
| `--setup-database`       | Setup database during backend initialization                                    | `FERRETDB_SETUP_DATABASE`       |                  |
| `--setup-username`       | Setup user during backend initialization                                        | `FERRETDB_SETUP_USERNAME`       |                  |
| `--setup-password`       | Setup user's password                                                           | `FERRETDB_SETUP_PASSWORD`       |                  |
| `--setup-timeout`        | Setup timeout                                                                   | `FERRETDB_SETUP_TIMEOUT`        | `30s`            |
| `--telemetry`            | Enable or disable [basic telemetry](telemetry.md)                               | `FERRETDB_TELEMETRY`            | `undecided`      |

<!-- Do not document `--test-XXX` flags here -->
