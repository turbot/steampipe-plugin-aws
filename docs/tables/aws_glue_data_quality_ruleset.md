# Table: aws_glue_data_quality_ruleset

AWS Glue Data Quality automatically computes statistics for your datasets. It uses these statistics to recommend a set of quality rules that checks for freshness, accuracy, and integrity. You can adjust recommended rules, discard rules, or add new rules as needed.

## Examples

### Basic info

```sql
select
  name,
  database_name,
  table_name,
  created_on,
  description,
  rule_set,
  recommendation_run_id
from
  aws_glue_data_quality_ruleset;
```

### List rulesets created in the last 30 days

```sql
select
  name,
  database_name,
  table_name,
  created_on,
  description,
  rule_set,
  recommendation_run_id
from
  aws_glue_data_quality_ruleset
where
  created_on >= now() - interval '30' day;
```

### Count ruleset by database

```sql
select
  database_name,
  count("name") as rulset_count
from
  aws_glue_data_quality_ruleset
group by
  database_name;
```

### Get Glue database details for a ruleset

```sql
select
  r.name,
  r.database_name,
  d.catalog_id,
  d.create_time as databse_create_time,
  d.location_uri
from
  aws_glue_data_quality_ruleset as r,
  aws_glue_catalog_database as d
where
  r.database_name = d.name
and
  r.name = 'ruleset1';
```

### Get Glue table details for a ruleset

```sql
select
  r.name,
  r.database_name,
  t.catalog_id,
  t.create_time as table_create_time,
  t.table_type
from
  aws_glue_data_quality_ruleset as r,
  aws_glue_catalog_table as t
where
  r.table_name = t.name
and
  r.name = 'ruleset1';
```

### Count rules per ruleset

```sql
select
  name,
  rule_count
from
  aws_glue_data_quality_ruleset;
```