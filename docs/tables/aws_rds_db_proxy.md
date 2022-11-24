# Table: aws_rds_db_proxy

RDS Proxy can allow applications to pool and share database connections to improve their ability to scale. RDS Proxy makes applications more resilient to database failures by automatically connecting to a standby DB instance while preserving application connections.

## Examples

### List of DB proxies and corresponding engine families

```sql
select
  db_proxy_name,
  status,
  engine_family
from
  aws_rds_db_proxy;
```

### Authentication info of each DB proxy.

```sql
select
  db_proxy_name,
  engine_family,
  a ->> 'AuthScheme' as auth_scheme,
  a ->> 'Description' as description,
  a ->> 'IAMAuth' as iam_auth,
  a ->> 'SecretArn' as secret_arn,
  a ->> 'UserName' as user_name
from
  aws_rds_db_proxy,
  jsonb_array_elements(auth) as a;
```
