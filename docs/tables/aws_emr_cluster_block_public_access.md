# Table: aws_emr_cluster_block_public_access

The `aws_emr_cluster_block_public_access` table provides Amazon EMR block public access configuration for your AWS account in the current Region..

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
  aws_emr_cluster_block_public_access
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
  aws_emr_cluster_block_public_access
where 
  block_public_security_group_rules;
```
