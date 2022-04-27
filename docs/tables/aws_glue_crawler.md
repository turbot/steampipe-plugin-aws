# Table: aws_glue_crawler

Crawler is used to populate the AWS Glue Data Catalog with tables. This is the primary method used by most AWS Glue users. A crawler can crawl multiple data stores in a single run. Upon completion, the crawler creates or updates one or more tables in your Data Catalog. Extract, transform, and load (ETL) jobs that you define in AWS Glue use these Data Catalog tables as sources and targets. The ETL job reads from and writes to the data stores that are specified in the source and target Data Catalog tables.

## Examples

### Basic info

```sql
select
  name,
  state,
  database_name,
  creation_time,
  description,
  recrawl_behavior
from
  aws_glue_crawler;
```

### List running crawlers

```sql
select
  name,
  state,
  database_name,
  creation_time,
  description,
  recrawl_behavior
from
  aws_glue_crawler
where
  state = 'RUNNING'; 
```