---
title: "Steampipe Table: aws_rolesanywhere_profile - Query AWS Roles Anywhere Trust Profiles using SQL"
description: "Allows users to query Roles Anywhere for detailed information about the profile configurations."
folder: "Roles Anywhere"
---

# Table: aws_rolesanywhere_profile - Query AWS Roles Anywhere Profiles using SQL

AWS Roles Anywhere enables trusted entities outside of an AWS account to carry out operations within that account by using certificates for IAM Role assumption. Profiles can be used to control the role sessions for a Trust Anchor by configuring features like boundary policies.

## Table Usage Guide

The `aws_rolesanywhere_profile` table in Steampipe provides you with information about Roles Anywhere Profiles. This table allows you, as a DevOps engineer, to query Profile-specific details, including boundary policy, certificate field mappings, and associated metadata. You can utilize this table to gather insights on Profiles, such as the IAM roles a Profile can access, which permission restrictions apply to role sessions, and more. The schema outlines the various attributes of the Profile for you, including the ARN, attribute mappings, IAM information, and create and update times.

## Examples

### List enabled Profiles.
Determine the Profiles that are currently enabled. 
This can be useful to determine which Profiles can create temporary session in the account.

```sql+postgres
select
  arn,
  role_arns
from
  aws_rolesanywhere_profile
where
  enabled;
```

```sql+sqlite
select
  arn,
  role_arns
from
  aws_rolesanywhere_profile
where
  enabled;
```

### Find IAM Roles that can be assumed by Profiles
Determine the IAM Roles that can be assumed by a Profile and their attached/inline IAM policies. 
This can be useful to determine the effective permissions for a Trust Anchor's Profile.

```sql+postgres
select 
  profile.arn as profile_arn, 
  role.arn as role_arn,
  role.attached_policy_arns as policy_arns,
  role.inline_policies as inline_policies
from 
  aws_rolesanywhere_profile as profile,
  jsonb_array_elements_text(profile.role_arns) as role_arn 
  join aws_iam_role as role on role_arn = role.arn
```

```sql+sqlite
select 
  profile.arn as profile_arn, 
  role.arn as role_arn,
  role.attached_policy_arns as policy_arns,
  role.inline_policies as inline_policies
from 
  aws_rolesanywhere_profile as profile,
  json_each(profile.role_arns) as role_arn 
  join aws_iam_role as role on role_arn = role.arn
```

### List Profiles that have a session policy.
Determine the Profiles that have a session policy configured. 
This can be useful to determine if temporary role sessions created using the Profile have restricted permissions relative to the role's identity policy(s).

```sql+postgres
select
  arn,
  session_policy
from
  aws_rolesanywhere_profile
where
  session_policy is not null;
```

```sql+sqlite
select
  arn,
  session_policy
from
  aws_rolesanywhere_profile
where
  session_policy is not null;
```
