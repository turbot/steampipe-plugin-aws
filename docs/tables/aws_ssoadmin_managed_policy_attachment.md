---
title: "Table: aws_ssoadmin_managed_policy_attachment - Query AWS SSO Managed Policy Attachments using SQL"
description: "Allows users to query AWS SSO Managed Policy Attachments, providing information about the managed policy attachments of AWS SSO permission sets."
---

# Table: aws_ssoadmin_managed_policy_attachment - Query AWS SSO Managed Policy Attachments using SQL

The `aws_ssoadmin_managed_policy_attachment` table in Steampipe provides information about the managed policy attachments of AWS SSO permission sets. This table allows DevOps engineers to query policy-specific details, including the instance ARN, permission set ARN, and managed policy ARN. Users can utilize this table to gather insights on policy attachments, such as the attached policies for each permission set, and more. The schema outlines the various attributes of the managed policy attachment, including the instance ARN, permission set ARN, and managed policy ARN.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssoadmin_managed_policy_attachment` table, you can use the `.inspect aws_ssoadmin_managed_policy_attachment` command in Steampipe.

**Key columns**:

- `instance_arn`: This is the ARN of the SSO instance under which the operation will be executed. It is an important column as it uniquely identifies the SSO instance.
- `permission_set_arn`: This is the ARN of the permission set. It is useful for joining with other tables that contain permission set information.
- `managed_policy_arn`: This is the ARN of the managed policy that is being attached to a permission set. It is useful for joining with other tables that contain managed policy information.

## Examples

### Basic info

```sql
select
  mpa.managed_policy_arn,
  mpa.name
from
  aws_ssoadmin_managed_policy_attachment as mpa
join
  aws_ssoadmin_permission_set as ps on mpa.permission_set_arn = ps.arn;
```
