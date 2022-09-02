# Table: aws_glue_dev_endpoint

A development endpoint is an environment that you can use to develop and test your AWS Glue scripts. You can use AWS Glue to create, edit, and delete development endpoints.

## Examples

### Basic info

```sql
select
  endpoint_name,
  status,
  database_name,
  availability_zone,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint;
```

### List dev endpoints that are not in ready state

```sql
select
  endpoint_name,
  status,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint
where
  status <> 'READY'; 
```

### List dev endpoints updated in the last 30 days

```sql
select
  title,
  arn,
  status,
  glue_version,
  last_modified_timestamp
from
  aws_glue_dev_endpoint
where
   last_modified_timestamp >= now() - interval '30' day;
```

### List dev endpoints older than 30 days

```sql
select
  endpoint_name,
  arn,
  status,
  glue_version,
  created_timestamp
from
  aws_glue_dev_endpoint
where
   created_timestamp >= now() - interval '30' day;
```

### Get subnet details attached to a particular dev endpoint

```sql
select
  e.endpoint_name,
  s.availability_zone,
  s.available_ip_address_count,
  s.cidr_block,
  s.default_for_az,
  s.map_customer_owned_ip_on_launch,
  s.map_public_ip_on_launch,
  s.state
from
  aws_glue_dev_endpoint as e,
  aws_vpc_subnet as s
where
  e.endpoint_name = 'test5'
and
  e.subnet_id = s.subnet_id;
```

### Get extra jars s3 bucket details for a dev endpoint 

```sql
select
  e.endpoint_name,
  split_part(j, '/', '3') as extra_jars_s3_bucket,
  b.versioning_enabled,
  b.policy,
  b.object_lock_configuration,
  b.restrict_public_buckets,
  b.policy
from
  aws_glue_dev_endpoint as e,
  aws_s3_bucket as b,
  unnest (string_to_array(e.extra_jars_s3_path, ',')) as j
where
  b.name = split_part(j, '/', '3')
and
  e.endpoint_name = 'test34';
```