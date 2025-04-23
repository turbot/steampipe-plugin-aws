---
title: "Steampipe Table: aws_inspector2_member - Query AWS Inspector Members using SQL"
description: "Allows users to query AWS Inspector Members to retrieve detailed information about the member accounts within an AWS Inspector assessment target."
folder: "Inspector2"
---

# Table: aws_inspector2_member - Query AWS Inspector Members using SQL

The AWS Inspector is a security assessment service that helps improve the security and compliance of applications deployed on AWS. It automatically assesses applications for exposure, vulnerabilities, and deviations from best practice. After performing an assessment, AWS Inspector produces a detailed list of security findings prioritized by level of severity.

## Table Usage Guide

The `aws_inspector2_member` table in Steampipe provides you with information about AWS Inspector Members. This table allows you, as a DevOps engineer, to query member-specific details, including account IDs, emails, and associated metadata. You can utilize this table to gather insights on member accounts, such as the account status, the account's relationship with the AWS Inspector assessment target, and more. The schema outlines the various attributes of the AWS Inspector Member for you, including the account ID, email, and the ARN of the AWS Inspector assessment target.

## Examples

### Basic info
Identify instances where the status of the relationship between member and admin accounts in AWS Inspector has changed, which can be useful for auditing or tracking changes over time.

```sql+postgres
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member;
```

```sql+sqlite
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member;
```

### Retrieve a list of members whose status hasn't changed in the past 30 days
Identify the members who have maintained a consistent status over the past month. This can be useful for tracking stability within your organization or for identifying members who may need attention or updates.

```sql+postgres
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

```sql+sqlite
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member
where
  updated_at >= datetime('now', '-30 day');
```

### List invited members
Explore which members have been invited to join your AWS Inspector service. This is useful for tracking pending invitations and managing your AWS Inspector member relationships.

```sql+postgres
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status
from
  aws_inspector2_member
where
  relationship_status = 'INVITED';
```

```sql+sqlite
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status
from
  aws_inspector2_member
where
  relationship_status = 'INVITED';
```