---
title: "Steampipe Table: aws_iam_instance_profile - Query AWS Identity and Access Management (IAM) Instance Profiles using SQL"
description: "Allows users to query IAM Instance Profiles to gain insights into their configurations, associated roles, and metadata."
folder: "IAM"
---

# Table: aws_iam_instance_profile - Query AWS Identity and Access Management (IAM) Instance Profiles using SQL

An IAM instance profile is a container for an IAM role that you can use to pass role information to an EC2 instance when the instance starts. An instance profile can contain only one role, and that limit cannot be increased. If you create a role by using the IAM console, the console creates an instance profile automatically and gives it the same name as the role to which it corresponds.

## Table Usage Guide

The `aws_iam_instance_profile` table in Steampipe provides you with information about IAM instance profiles within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query instance profile-specific details, including associated roles, creation date, and associated metadata. You can utilize this table to gather insights on instance profiles, such as which roles are attached to each instance profile, identify instance profiles without roles, verification of instance profile configurations, and more. The schema outlines the various attributes of the IAM instance profile for you, including the instance profile ARN, creation date, associated roles, and associated tags.

## Examples

### List all IAM instance profiles
Discover all IAM instance profiles in your AWS account to understand what roles are available for EC2 instances.

```sql+postgres
select
  instance_profile_name,
  arn,
  create_date,
  path
from
  aws_iam_instance_profile
order by
  instance_profile_name;
```

```sql+sqlite
select
  instance_profile_name,
  arn,
  create_date,
  path
from
  aws_iam_instance_profile
order by
  instance_profile_name;
```

### List instance profiles with their associated roles
Identify the relationships between instance profiles and their associated IAM roles to understand permission assignments.

```sql+postgres
select
  ip.instance_profile_name as instance_profile_name,
  ip.arn as instance_profile_arn,
  role ->> 'RoleName' as role_name,
  role ->> 'Arn' as role_arn
from
  aws_iam_instance_profile as ip,
  jsonb_array_elements(roles) as role;
```

```sql+sqlite
select
  ip.instance_profile_name as instance_profile_name,
  ip.arn as instance_profile_arn,
  json_extract(role.value, '$.RoleName') as role_name,
  json_extract(role.value, '$.Arn') as role_arn
from
  aws_iam_instance_profile as ip,
  json_each(roles) as role;
```

### Find instance profiles without any associated roles
Identify instance profiles that don't have any roles attached, which may indicate incomplete configurations.

```sql+postgres
select
  instance_profile_name,
  arn,
  create_date
from
  aws_iam_instance_profile
where
  roles is null
  or jsonb_array_length(roles) = 0;
```

```sql+sqlite
select
  instance_profile_name,
  arn,
  create_date
from
  aws_iam_instance_profile
where
  roles is null
  or json_array_length(roles) = 0;
```

### List instance profiles created in the last 30 days
Identify recently created instance profiles to track new configurations and potential security implications.

```sql+postgres
select
  instance_profile_name,
  arn,
  create_date,
  path
from
  aws_iam_instance_profile
where
  create_date >= now() - interval '30 days'
order by
  create_date desc;
```

```sql+sqlite
select
  instance_profile_name,
  arn,
  create_date,
  path
from
  aws_iam_instance_profile
where
  create_date >= datetime('now', '-30 days')
order by
  create_date desc;
```

### Get instance profiles with specific tags
Find instance profiles that have specific tags to understand organization and categorization.

```sql+postgres
select
  instance_profile_name,
  arn,
  tags
from
  aws_iam_instance_profile
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  instance_profile_name,
  arn,
  tags
from
  aws_iam_instance_profile
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List instance profiles with roles that have specific permissions
Identify instance profiles whose associated roles have specific policies attached.

```sql+postgres
select
  ip.instance_profile_name as instance_profile_name,
  role ->> 'RoleName' as role_name,
  p.policy_name
from
  aws_iam_instance_profile as ip,
  jsonb_array_elements(roles) as role,
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as policy_arn,
  aws_iam_policy as p
where
  r.arn = role ->> 'Arn'
  and p.arn = policy_arn
  and p.policy_name like '%EC2%';
```

```sql+sqlite
select
  ip.instance_profile_name as instance_profile_name,
  json_extract(role.value, '$.RoleName') as role_name,
  p.policy_name
from
  aws_iam_instance_profile as ip,
  json_each(roles) as role,
  aws_iam_role as r,
  json_each(r.attached_policy_arns) as policy_arn,
  aws_iam_policy as p
where
  r.arn = json_extract(role.value, '$.Arn')
  and p.arn = policy_arn.value
  and p.policy_name like '%EC2%';
```

### Count instance profiles by path prefix
Analyze the distribution of instance profiles across different organizational paths.

```sql+postgres
select
  substring(path from 1 for position('/' in substring(path from 2)) + 1) as path_prefix,
  count(*) as instance_profile_count
from
  aws_iam_instance_profile
where
  path != '/'
group by
  path_prefix
order by
  instance_profile_count desc;
```

```sql+sqlite
select
  substr(path, 1, instr(substr(path, 2), '/') + 1) as path_prefix,
  count(*) as instance_profile_count
from
  aws_iam_instance_profile
where
  path != '/'
group by
  path_prefix
order by
  instance_profile_count desc;
```
