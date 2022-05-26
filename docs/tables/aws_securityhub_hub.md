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


### List hubs that do not automatically enable new controls

```sql
select
  hub_arn,
  auto_enable_controls
from
  aws_securityhub_hub
where
  not auto_enable_controls;
```

### List administrator account details for the hub 

```sql
select
  hub_arn,
  auto_enable_controls,
  administrator_account ->> 'AccountId' as administrator_account_id,
  administrator_account ->> 'InvitationId' as administrator_invitation_id,
  administrator_account ->> 'InvitedAt' as administrator_invitation_time,
  administrator_account ->> 'MemberStatus' as administrator_status
from
  aws_securityhub_hub
where
  administrator_account is not null;
```