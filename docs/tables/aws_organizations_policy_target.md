---
title: "Steampipe Table: aws_organizations_policy_target - Query AWS Organizations Policy Targets using SQL"
description: "Allows users to query AWS Organizations Policy Targets to retrieve detailed information about the application of policies to roots, organizational units (OUs), and accounts."
folder: "Organizations"
---

# Table: aws_organizations_policy_target - Query AWS Organizations Policy Targets using SQL

The AWS Organizations Policy Target is a part of AWS Organizations service. It allows you to attach policies to roots, organizational units (OUs), or accounts in your organization, thereby enabling centralized control over your AWS workloads. It simplifies the process of managing permissions and enhances the security of your AWS resources.

## Table Usage Guide

The `aws_organizations_policy_target` table in Steampipe provides you with information about policy targets within AWS Organizations. This table allows you, as a DevOps engineer, to query policy target-specific details, including the policy ID, target ID, and the type of target (root, OU, or account). You can utilize this table to gather insights on policy applications, such as which policies are applied to which roots, OUs, or accounts, and the status of these applications. The schema outlines the various attributes of the policy target for you, including the ARN, policy ID, target ID, and target type.

**Important Notes**
- You must specify `type` and `target_id` in the `where` clause to query this table.

## Examples

### Basic info
Explore which AWS managed services are controlled by a specific policy. This is particularly useful for assessing security measures and ensuring only authorized services are being used within an organization.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  aws_managed 
from
  aws_organizations_policy_target 
where
  type = 'SERVICE_CONTROL_POLICY' 
  and target_id = '123456789098';
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  aws_managed 
from
  aws_organizations_policy_target 
where
  type = 'SERVICE_CONTROL_POLICY' 
  and target_id = '123456789098';
```

### List tag policies of a targeted organization that are not managed by AWS
Determine the areas in which your organization has implemented tag policies that are not directly managed by AWS. This is useful for maintaining oversight of policy management and ensuring compliance with your organization's specific requirements.

```sql+postgres
select
  id,
  name,
  arn,
  type,
  aws_managed 
from
  aws_organizations_policy_target 
where
  not aws_managed 
  and type = 'TAG_POLICY' 
  and target_id = 'ou-jsdhkek';
```

```sql+sqlite
select
  id,
  name,
  arn,
  type,
  aws_managed 
from
  aws_organizations_policy_target 
where
  not aws_managed 
  and type = 'TAG_POLICY' 
  and target_id = 'ou-jsdhkek';
```

### List backup organization policies of an account
Explore which backup policies are managed by AWS for a specific account. This is useful for assessing the security measures in place and identifying any potential vulnerabilities.

```sql+postgres
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy_target
where
  type = 'BACKUP_POLICY'
  and target_id = '123456789098';
```

```sql+sqlite
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy_target
where
  type = 'BACKUP_POLICY'
  and target_id = '123456789098';
```

### Get policy details of the service control policies of a root account
Determine the specifics of service control policies linked to a root account to understand the policy version and statements. This can be useful in managing and auditing policy details for security and compliance purposes.

```sql+postgres
select
  name,
  id,
  content ->> 'Version' as policy_version,
  content ->> 'Statement' as policy_statement
from
  aws_organizations_policy_target
where
  type = 'SERVICE_CONTROL_POLICY'
  and target_id = 'r-9ijkl7';
```

```sql+sqlite
select
  name,
  id,
  json_extract(content, '$.Version') as policy_version,
  json_extract(content, '$.Statement') as policy_statement
from
  aws_organizations_policy_target
where
  type = 'SERVICE_CONTROL_POLICY'
  and target_id = 'r-9ijkl7';
```