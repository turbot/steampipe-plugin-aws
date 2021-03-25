# Table: aws_securityhub_hub

Security Hub provides a comprehensive view of the security state of AWS environment and resources. It also provides the readiness status of environment based on controls from supported security standards.

## Examples

### Basic info

```sql
select
  hub_arn,
  auto_enable_controls,
  subscribed_at,
  region
from
  aws_securityhub_hub;
```


### List of hubs where Auto enable controls is disable

```sql
select
  hub_arn,
  auto_enable_controls
from
  aws_securityhub_hub
where
  not auto_enable_controls;
```


### Check if the hub missing required tags

```sql
select
  hub_arn,
  tags
from
  aws_securityhub_hub
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```