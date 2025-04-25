---
title: "Steampipe Table: aws_ssoadmin_managed_policy_attachment - Query AWS SSO Managed Policy Attachments using SQL"
description: "Allows users to query AWS SSO Managed Policy Attachments, providing information about the managed policy attachments of AWS SSO permission sets."
folder: "SSO"
---

# Table: aws_ssoadmin_managed_policy_attachment - Query AWS SSO Managed Policy Attachments using SQL

The AWS SSO Managed Policy Attachment is a feature of AWS Single Sign-On (SSO) service. It allows you to attach and manage access permissions for AWS SSO users and groups through managed policies. This helps in streamlining the process of assigning permissions, ensuring secure access to AWS resources.

## Table Usage Guide

The `aws_ssoadmin_managed_policy_attachment` table in Steampipe provides you with information about the managed policy attachments of AWS SSO permission sets. This table allows you, as a DevOps engineer, to query policy-specific details, including the instance ARN, permission set ARN, and managed policy ARN. You can utilize this table to gather insights on policy attachments, such as the attached policies for each permission set, and more. The schema outlines the various attributes of the managed policy attachment for you, including the instance ARN, permission set ARN, and managed policy ARN.

## Examples

### Basic info
Analyze the connection between AWS SSO managed policy attachments and permission sets to understand the allocation of permissions within your AWS environment. This can help you maintain security and compliance by ensuring correct policy attachments.

```sql+postgres
select
  mpa.managed_policy_arn,
  mpa.name
from
  aws_ssoadmin_managed_policy_attachment as mpa
join
  aws_ssoadmin_permission_set as ps on mpa.permission_set_arn = ps.arn;
```

```sql+sqlite
select
  mpa.managed_policy_arn,
  mpa.name
from
  aws_ssoadmin_managed_policy_attachment as mpa
join
  aws_ssoadmin_permission_set as ps on mpa.permission_set_arn = ps.arn;
```