# Table: aws_glue_catalog_table

AWS Glue catalog table is the metadata definition that represents your data, including its schema. A table can be used as a source or target in a job definition.

## Examples

### Basic info

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  database_name
from
  aws_glue_catalog_table;
```

### Count the number of tables per catalog

```sql
select
  catalog_id,
  count(name) as table_count
from
  aws_glue_catalog_table
group by
  catalog_id;
```

### List tables with retention less than 30 days

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  retention
from
  aws_glue_catalog_table
where
  retention::bigint < 30;
```
