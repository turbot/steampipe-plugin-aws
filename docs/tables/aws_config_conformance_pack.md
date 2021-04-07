# Table: aws_config_conformance_pack

A conformance pack is a collection of AWS Config rules and remediation actions that can be easily deployed as a single entity in an account and a Region or across an organization in AWS Organizations.

## Examples

### Basic info

```sql
select
  name,
  conformance_pack_id,
  created_by,
  last_update_requested_time,
  title,
  akas
from
  aws_config_conformance_pack;
```


### Get S3 bucket info for each conformance pack

```sql
select
  name,
  conformance_pack_id,
  delivery_s3_bucket,
  delivery_s3_key_prefix
from
  aws_config_conformance_pack;
```


### Get input parameter details of each conformance pack

```sql
select
  name,
  inp ->> 'ParameterName' as parameter_name,
  inp ->> 'ParameterValue' as parameter_value,
  title,
  akas
from
  aws_config_conformance_pack,
  jsonb_array_elements(input_parameters) as inp;
```

