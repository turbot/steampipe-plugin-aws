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

### List configurations shared with other accounts

```sql
select
  title as resource,
  case
    when authorized_account_id is not null then 'alarm'
    else 'ok'
  end as status,
  case
    when authorized_account_id is not null then title || ' is sharing configuration and compliance data with account ' || authorized_account_id || '.'
    else title || ' is not sharing configuration and compliance data with any other account.'
  end as reason,
  -- Additional columns
  region,
  account_id
from
  aws_config_aggregate_authorization;
```