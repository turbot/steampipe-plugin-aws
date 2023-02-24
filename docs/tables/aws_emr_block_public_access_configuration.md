# Table: aws_emr_block_public_access_configuration

The `aws_emr_block_public_access_configuration` table provides Amazon EMR block public access configuration for your AWS account in the current Region.

## Examples

### Basic info

```sql
select
  created_by_arn,
  block_public_security_group_rules,
  creation_date,
  classification,
  configurations,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
order by
  created_by_arn,
  creation_date;
```


### List Block Public Access enabled settings for emr

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where 
  block_public_security_group_rules;
```


### List Block Public Access with Permitted Public Security Group Rule Max and Min Ranges 

```sql
select
  created_by_arn,
  creation_date,
  rules ->> 'MaxRange' as max_range,
  rules ->> 'MinRange' as min_range
from
  aws_emr_block_public_access_configuration
  cross join jsonb_array_elements(permitted_public_security_group_rule_ranges) as rules;
```


### List EMR Block Public Access Configuration created in last 90 days

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where
  date_part('day', now() - creation_date) < 90;
```


### Configuration details of EMR Block Public Access

```sql
select
  created_by_arn,
  creation_date,
  config ->> 'Classification' as classification,
  config ->> 'Configurations' as configurations,
  config ->> 'Properties' as properties
from
  aws_emr_block_public_access_configuration
  cross join jsonb_array_elements(configurations) as config;
```