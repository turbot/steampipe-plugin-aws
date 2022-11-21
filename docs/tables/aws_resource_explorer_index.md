# Table: aws_resource_explorer_index

Retrieves the indexes in AWS regions that are currently collecting resource
information for AWS Resource Explorer.

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

### Get the details for the aggregator index

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
