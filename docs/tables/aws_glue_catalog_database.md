# Table: aws_glue_catalog_database

AWS Glue Database is a set of associated table definitions, organized into a logical group.

## Examples

### Basic info

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  location_uri
from
  aws_glue_catalog_database;
```


### List databases by Catalog ID

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  location_uri
from
  aws_glue_catalog_database
where
  catalog_id = '453319552164';
```
