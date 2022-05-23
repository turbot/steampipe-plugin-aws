# Table: aws_config_aggregate_authorization

AWS Config provides a detailed view of the configuration of AWS resources in your AWS account. You can add authorization to grant permission to aggregator accounts and regions to collect AWS Config configuration and compliance data.

## Examples

### Basic info

```sql
select
  arn,
  authorized_account_id,
  authorized_aws_region,
  creation_time
from
  aws_config_aggregate_authorization;
```
