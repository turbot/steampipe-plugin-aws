---
title: "Table: aws_inspector2_member - Query AWS Inspector Members using SQL"
description: "Allows users to query AWS Inspector Members to retrieve detailed information about the member accounts within an AWS Inspector assessment target."
---

# Table: aws_inspector2_member - Query AWS Inspector Members using SQL

The `aws_inspector2_member` table in Steampipe provides information about AWS Inspector Members. This table allows DevOps engineers to query member-specific details, including account IDs, emails, and associated metadata. Users can utilize this table to gather insights on member accounts, such as the account status, the account's relationship with the AWS Inspector assessment target, and more. The schema outlines the various attributes of the AWS Inspector Member, including the account ID, email, and the ARN of the AWS Inspector assessment target.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector2_member` table, you can use the `.inspect aws_inspector2_member` command in Steampipe.

### Key columns:

- `account_id`: This is the AWS account ID of the member account. This column is useful for joining with other tables that contain AWS account data.
- `email`: This column contains the email of the member account. It can be used to join with other tables that contain email-related data.
- `arn`: This column contains the ARN of the AWS Inspector assessment target. It can be used to join with other tables that contain AWS Inspector assessment target data.

## Examples

### Basic info

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member;
```

### Retrieve a list of members whose status hasn't changed in the past 30 days

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member
where
  updated_at >= now() - interval '30' day;
```

### List invited members

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status
from
  aws_inspector2_member
where
  relationship_status = 'INVITED';
```