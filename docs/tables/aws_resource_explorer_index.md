# Table: aws_resource_explorer_index

Retrieves the indexes in AWS Regions that are currently collecting resource information for AWS Resource Explorer.

## Examples

### Basic info

```sql
select
  arn,
  region,
  type
from
  aws_resource_explorer_index;
```

### Get the details of `AGGREGATOR` index in the account

```sql
select
  arn,
  region,
  type
from
  aws_resource_explorer_index
where
  type = 'AGGREGATOR';
```
