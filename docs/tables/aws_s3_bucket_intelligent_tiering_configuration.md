# Table: aws_s3_bucket_intelligent_tiering_configuration

AWS S3 Bucket Intelligent Tiering Configuration is a feature of Amazon S3 that enables automatic optimization of storage costs for objects in an S3 bucket. With Intelligent Tiering, Amazon S3 monitors the access patterns of your objects and moves them between two storage tiers: frequent access tier and infrequent access tier.

The Intelligent Tiering feature analyzes the access patterns and automatically moves objects that have not been accessed for a certain period of time to the infrequent access tier, which has a lower storage cost. If an object in the infrequent access tier is accessed again, it is automatically moved back to the frequent access tier.


## Examples

### Basic info

```sql
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration;
```

### Get buckets intelligent tiering configure status

```sql
with intelligent_tiering_configuration as MATERIALIZED (
select
  bucket_name, id, status
from
  aws_s3_bucket_intelligent_tiering_configuration ),
  bucket as MATERIALIZED (
  select
    name, region
  from
    aws_s3_bucket )
    select distinct
      b.name,
      b.region,
      case
        when
          i.id is null
        then
          'Bucket do not have intelligent tiering configured'
        else
          'Bucket have intelligent tiering configured'
      end
      as intelligent_tiering_configuration_status
    from
      bucket as b
      left join
        intelligent_tiering_configuration as i
        on b.name = i.bucket_name;
```

### List enabled intelligent tiering configurations of each bucket

```sql
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration
where
  status = 'Enabled';
```

### Get tiering details of each intelligent tiering configurations

```sql
select
  bucket_name,
  id,
  status,
  t ->> 'AccessTier' as access_tier,
  t ->> 'Days' as days
from
  aws_s3_bucket_intelligent_tiering_configuration,
  jsonb_array_elements(tierings) as t;
```

### Get filter details of intelligent tiering configurations

```sql
select
  bucket_name,
  id,
  filter -> 'And' as filter_and,
  filter -> 'Prefix' as filter_prefix,
  filter -> 'Tag' as filter_tag
from
  aws_s3_bucket_intelligent_tiering_configuration;
```
