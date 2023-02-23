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
### List clusters with Block Public Access enabled

```sql
select
  created_by_arn,
  creation_date,
  configurations,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
where 
  block_public_security_group_rules;
```
