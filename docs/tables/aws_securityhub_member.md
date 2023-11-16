---
title: "Table: aws_securityhub_member - Query AWS Security Hub Members using SQL"
description: "Allows users to query AWS Security Hub Members for detailed information about each member's account, including its ID, email, status, and more."
---

# Table: aws_securityhub_member - Query AWS Security Hub Members using SQL

The `aws_securityhub_member` table in Steampipe provides information about each member account within AWS Security Hub. This table allows DevOps engineers to query member-specific details, including account ID, email, status, and the timestamp of the invitation. Users can utilize this table to gather insights on member accounts, such as their invitation and verification status, the email associated with each account, and more. The schema outlines the various attributes of the Security Hub member, including the member account ID, email, status, invited at timestamp, and updated at timestamp.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_member` table, you can use the `.inspect aws_securityhub_member` command in Steampipe.

### Key columns:

- `account_id`: This is the AWS account ID of the member account. It can be used to join this table with others that contain information about AWS accounts.
- `email`: This is the email of the member account. It can be used to join this table with others that contain information about the email addresses associated with AWS accounts.
- `status`: This is the status of the member account. It can be used to filter or sort the table based on the status of the member accounts.

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