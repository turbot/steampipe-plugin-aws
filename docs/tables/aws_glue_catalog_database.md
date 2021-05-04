# Table: aws_glue_catalog_database

AWS Glue Catalog Database is a set of associated table definitions, organized into a logical group.

## Examples

### Basic info

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  location_uri,
  create_table_default_permissions
from
  aws_glue_catalog_database;
```


### Count databases per Catalog ID

```sql
select
  catalog_id,
  count(name)
from
  aws_glue_catalog_database
group by
  catalog_id;
```
