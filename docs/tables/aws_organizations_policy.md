---
title: "Steampipe Table: aws_organizations_policy - Query AWS Organizations Policy using SQL"
description: "Allows users to query AWS Organizations Policy to retrieve detailed information on policies within AWS Organizations. This table can be utilized to gain insights on policy-specific details, such as policy type, content, and associated metadata."
folder: "Organizations"
---

# Table: aws_organizations_policy - Query AWS Organizations Policy using SQL

The AWS Organizations Policy is a service that allows you to centrally manage and enforce policies for multiple AWS accounts. With this service, you can create policies that apply across your organization, or to specific organizational units (OUs) or accounts. It provides control over AWS service use, ensures consistent tags, and helps maintain your accounts as per your company's compliance requirements.

## Table Usage Guide

The `aws_organizations_policy` table in Steampipe provides you with information about policies within AWS Organizations. This table allows you, as a DevOps engineer, to query policy-specific details, including policy type, content, and associated metadata. You can utilize this table to gather insights on policies, such as policy names, policy types, and the contents of the policies. The schema outlines the various attributes of the policy for you, including the policy ARN, policy type, policy content, policy name, and associated tags.

**Important Notes**
- You must specify `type` in the `where` clause to query this table.

## Examples

### Basic info
Analyze the settings to understand the policies managed by AWS within your organization. This is particularly useful for gaining insights into service control policies, allowing you to manage access across your AWS environment effectively.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```

### List tag policies that are not managed by AWS
Identify tag policies within your AWS organization that are not directly managed by AWS. This can be useful for maintaining oversight of custom policies and ensuring they align with your organization's specific requirements and standards.

```sql+postgres
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  not aws_managed
  and type = 'TAG_POLICY';
```

```sql+sqlite
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  not aws_managed
  and type = 'TAG_POLICY';
```

### List backup policies
Explore the list of backup policies in your AWS organization to understand which ones are managed by AWS and which ones you've implemented. This can help you maintain compliance and ensure data recovery in case of accidental deletion or system failure.

```sql+postgres
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'BACKUP_POLICY';
```

```sql+sqlite
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'BACKUP_POLICY';
```

### Get policy details of the service control policies
Determine the details of service control policies within your AWS organization. This query can help you understand the version and statement of each policy, providing valuable insights for policy management and compliance.

```sql+postgres
select
  name,
  id,
  content ->> 'Version' as policy_version,
  content ->> 'Statement' as policy_statement
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```

```sql+sqlite
select
  name,
  id,
  json_extract(content, '$.Version') as policy_version,
  json_extract(content, '$.Statement') as policy_statement
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```
### List the targets each service control policy is attached to
Find where each service control policy is attached in your organization. Each policy's targets are listed with one call per policy, so this is quicker than looking up the policies of every root, organizational unit and account.

```sql+postgres
select
  p.name,
  t ->> 'Name' as target_name,
  t ->> 'Type' as target_type,
  t ->> 'TargetId' as target_id
from
  aws_organizations_policy as p,
  jsonb_array_elements(p.targets) as t
where
  p.type = 'SERVICE_CONTROL_POLICY';
```

```sql+sqlite
select
  p.name,
  json_extract(t.value, '$.Name') as target_name,
  json_extract(t.value, '$.Type') as target_type,
  json_extract(t.value, '$.TargetId') as target_id
from
  aws_organizations_policy as p,
  json_each(p.targets) as t
where
  p.type = 'SERVICE_CONTROL_POLICY';
```

### List service control policies that are not attached to any target
Identify customer managed service control policies that aren't attached to a root, organizational unit or account, and so have no effect.

```sql+postgres
select
  name,
  id,
  arn
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY'
  and not aws_managed
  and jsonb_array_length(targets) = 0;
```

```sql+sqlite
select
  name,
  id,
  arn
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY'
  and aws_managed = 0
  and json_array_length(targets) = 0;
```
