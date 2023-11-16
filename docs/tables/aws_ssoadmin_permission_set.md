---
title: "Table: aws_ssoadmin_permission_set - Query AWS SSO Admin Permission Set using SQL"
description: "Allows users to query AWS SSO Admin Permission Set to retrieve data related to the permissions sets of AWS Single Sign-On (SSO) service."
---

# Table: aws_ssoadmin_permission_set - Query AWS SSO Admin Permission Set using SQL

The `aws_ssoadmin_permission_set` table in Steampipe provides information about the permission sets associated with AWS Single Sign-On (SSO) service. This table allows DevOps engineers to query permission set-specific details, including the permission set name, description, created date, and related metadata. Users can utilize this table to gather insights on permission sets, such as the instances of each permission set, associated policies, and more. The schema outlines the various attributes of the permission set, including the permission set ARN, created date, session duration, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssoadmin_permission_set` table, you can use the `.inspect aws_ssoadmin_permission_set` command in Steampipe.

**Key columns**:

- `permission_set_arn`: The ARN of the permission set. This can be used to join with other tables that reference permission sets.
- `created_date`: The date and time when the permission set was created. This can be used for tracking and auditing purposes.
- `session_duration`: The length of time that the SSO session is valid. This can be used for managing session lifetimes and security configurations.

## Examples

### Basic info

```sql
select
  name,
  arn,
  created_date,
  description,
  relay_state,
  session_duration,
  tags
from
  aws_ssoadmin_permission_set;
```
