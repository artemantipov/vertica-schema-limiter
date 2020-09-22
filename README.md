# vertica-schema-limiter
Tool for [Vertica](https://www.vertica.com/) to limit usage on schemas after specified threshold with prometheus metrics endpoint.

### Usage
* Create schema
* Create role with RW permissions to schema above
* Create service user with dbadmin privileges to check service tables for schemas size and revoke usage on schemas roles if threshold exceeded
* Configure vertica-schema-limiter
* Run service

### Configuration
Config file with additional env variable override option.
Default location: /config.yaml (path could be reassigned with env `CONFIG_PATH`)

#### Config example
```yaml
vertica:
  user: dbadmin
  pass: password
  host: vertica.host.com
  port: 5433
  db: main
  checkinterval: 5  #How often check schemas size in minutes
schemas:
  example_schema:  #Name of schema
    limit: 100   #Size of limited schema in Gb
    role: RW_example_schema_role  #Name of RW role for schema
```
For override any parameters above use env with prefix `LIMITER_` and `__` for next level in config

Example env for overriding vertica hostname is `LIMITER_VERTICA__HOST`

### Metrics
All schema measure results provided as prometheus metrics and available at `0.0.0.0:2112/metrics`

Metric for example above:
`vertica_limiter_schema_size{schema="example_schema"} 13.37`
 