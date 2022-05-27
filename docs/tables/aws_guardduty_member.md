# Table: aws_guardduty_member

AWS GuardDuty member resource can be used to add an AWS account as a GuardDuty member account to the current GuardDuty administrator account. If the value of the Status property is not provided or is set to Created, a member account is created but not invited. If the value of the Status property is set to Invited, a member account is created and invited.

## Examples

### Basic info

```sql
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member;
```

### List members that failed email verification

```sql
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  relationship_status = 'EmailVerificationFailed';
```

### List uninvited members

```sql
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at is null;
```

### List members which were invited within the last 10 days

```sql
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at >= (now() - interval '10' day);
```
