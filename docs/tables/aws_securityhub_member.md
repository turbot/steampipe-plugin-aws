# Table: aws_securityhub_member

An AWS Security Hub member

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
  member_status='Enabled';
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
  member_status='Created';
```

### List members which were enabled within the last 10 days

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
  member_status='Enabled'
and
  date_part('day', now() - invited_at) <= 10;
```