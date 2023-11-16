---
title: "Table: aws_guardduty_member - Query AWS GuardDuty Member using SQL"
description: "Allows users to query AWS GuardDuty Member data, including member account details, detector ID, invitation status, and relationship status."
---

# Table: aws_guardduty_member - Query AWS GuardDuty Member using SQL

The `aws_guardduty_member` table in Steampipe provides information about member accounts within AWS GuardDuty. This table allows security analysts to query member-specific details, including account details, detector ID, invitation status, and relationship status. Users can utilize this table to gather insights on member accounts, such as the status of invitations sent to these accounts, the relationship status between the master and member accounts, and more. The schema outlines the various attributes of the GuardDuty member, including the account ID, email, detector ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_member` table, you can use the `.inspect aws_guardduty_member` command in Steampipe.

**Key columns**:

- `account_id`: This is the AWS Account ID of the GuardDuty member account. It is a key column that can be used to join this table with other AWS tables to gather more comprehensive data about the member account.
- `detector_id`: This is the detector ID of the GuardDuty member account. It can be used to join this table with other GuardDuty tables to gather detailed information about the specific detector associated with the member account.
- `relationship_status`: This column indicates the status of the relationship between the master and member accounts. It can be useful when querying the status of member accounts in relation to the master account.

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
