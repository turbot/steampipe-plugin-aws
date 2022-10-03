# Table: aws_securityhub_member

An AWS Security Hub member account. An administrator account can view data from and manage configuration for its member accounts.

## Examples

### Basic info

```sql
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member;
```

### List members which are enabled

```sql
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled';
```

### List members which are invited but did not accept

```sql
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member
where
  member_status = 'Created';
```

### List members which were invited within the last 10 days

```sql
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled'
and
  invited_at <= (now() - interval '10' day);
```