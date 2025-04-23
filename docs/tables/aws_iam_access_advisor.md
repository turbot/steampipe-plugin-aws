---
title: "Steampipe Table: aws_iam_access_advisor - Query AWS IAM Access Advisor using SQL"
description: "Allows users to query AWS IAM Access Advisor to retrieve information about the service last accessed data for IAM entities (users, groups, and roles)."
folder: "IAM"
---

# Table: aws_iam_access_advisor - Query AWS IAM Access Advisor using SQL

The AWS IAM Access Advisor is a tool that helps you view and understand the permissions granted to your IAM entities (users, groups of users, or roles). It provides you with service last accessed information to help you refine your IAM policies. This service can be queried using SQL to fetch detailed insights.

## Table Usage Guide

The `aws_iam_access_advisor` table in Steampipe provides you with information about the service last accessed data for IAM entities like users, groups, and roles within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query details about the services that IAM entities can access, the actions they can perform, and when they last accessed the services. You can utilize this table to gather insights on access patterns, such as identifying unused permissions or verifying least privilege policies. The schema outlines the various attributes of the IAM entity, including the entity ARN, last accessed time, and the services accessible.

**Important Notes**

- You ***must*** specify a single `principal_arn` in a `where` or `join` clause in order to use this table.  

- The service last accessed data includes all your attempts to access an AWS API, not just the successful ones. This includes all attempts that you have made using the AWS Management Console, the AWS API through any of the SDKs, or any of the command line tools. An unexpected entry in the service last accessed data does not mean that your account has been compromised, because the request might have been denied. Refer to your CloudTrail logs as the authoritative source for information about all API calls and whether they were successful or denied access. 

- Service last accessed data does not use other policy types when determining whether a resource could access a service. These other policy types include resource-based policies, access control lists, AWS Organizations policies, IAM permissions boundaries, and AWS STS assume role policies. It only applies permissions policy logic. For more about the evaluation of policy types, see Evaluating Policies in the IAM User Guide.

## Examples

### Show the most recently used services for a user, role, group, or policy
Discover the segments that have recently interacted with a specific user, role, group, or policy. This helps in understanding the usage pattern and can aid in optimizing resource allocation.

```sql+postgres
select 
  principal_arn,
  service_name,
  last_authenticated,
  age(last_authenticated::date) 
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::123456789123:user/john'
  and last_authenticated is not null
order by 
  age asc;
```

```sql+sqlite
select 
  principal_arn,
  service_name,
  last_authenticated,
  julianday('now') - julianday(last_authenticated) as age
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::123456789123:user/john'
  and last_authenticated is not null
order by 
  age asc;
```

### Show unused services granted a user, role, group, or policy (unused in the last year)
This example helps to identify unused services that have been granted to a user, role, group, or policy, specifically those that have not been used in the past year. This is useful for maintaining security and efficiency by ensuring that unnecessary access permissions are revoked.

```sql+postgres
select 
  principal_arn,
  service_name
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::123456789123:role/turbot/admin'
  and last_authenticated is null
order by 
  service_name
```

```sql+sqlite
select 
  principal_arn,
  service_name
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::123456789123:role/turbot/admin'
  and last_authenticated is null
order by 
  service_name
```

### Show the last time a policy was used to access services, and the user, role, or group that used it
Determine the last time a specific policy was used to access services in AWS, along with the user, role, or group that used it. This can help in assessing the activity and access patterns related to your AWS resources, aiding in their management and security.

```sql+postgres
select 
  principal_arn,
  service_name,
  last_authenticated,
  age(last_authenticated::date),
  last_authenticated_entity,
  last_authenticated_region
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::aws:policy/AdministratorAccess'
  and last_authenticated is not null
order by 
  age asc;
```

```sql+sqlite
select 
  principal_arn,
  service_name,
  last_authenticated,
  julianday('now') - julianday(last_authenticated),
  last_authenticated_entity,
  last_authenticated_region
from 
  aws_iam_access_advisor
where
  principal_arn = 'arn:aws:iam::aws:policy/AdministratorAccess'
  and last_authenticated is not null
order by 
  julianday('now') - julianday(last_authenticated) asc;
```

### Show unused services granted to a role, including the policy that grants access and the actions granted  
Determine the areas in which services have been granted to a role but remain unused, including the specific policy that provides access and the actions permitted. This is particularly useful for identifying potential security risks and optimizing resource usage by revoking unnecessary permissions.

```sql+postgres
select 
  adv.service_name,
  action as action_granted,
  attached as granted_in,
  adv.service_namespace
from 
  aws_iam_access_advisor as adv,
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as attached,
  aws_iam_policy as p,  
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action
where
  principal_arn = 'arn:aws:iam::123456789123:role/turbot/admin'
  and r.arn = adv.principal_arn
  and last_authenticated is null
  and attached  = p.arn
  and stmt ->> 'Effect' = 'Allow'
  and action like adv.service_namespace || ':%'
order by 
  adv.service_name;
```

```sql+sqlite
select 
  adv.service_name,
  action.value as action_granted,
  attached.value as granted_in,
  adv.service_namespace
from 
  aws_iam_access_advisor as adv,
  aws_iam_role as r,
  json_each(r.attached_policy_arns) as attached,
  aws_iam_policy as p,  
  json_each(p.policy_std, '$.Statement') as stmt,
  json_each(stmt.value, '$.Action') as action
where
  principal_arn = 'arn:aws:iam::123456789123:role/turbot/admin'
  and r.arn = adv.principal_arn
  and last_authenticated is null
  and attached.value  = p.arn
  and json_extract(stmt.value, '$.Effect') = 'Allow'
  and action.value like adv.service_namespace || ':%'
order by 
  adv.service_name;
```

### Show action-level last accessed info (currently, only supported for S3)
This query allows you to gain insights into the most recent actions taken within the S3 service for a specific user. It helps in identifying key details such as the action name, the entity accessed, the region of access, and the time of access, which can be useful for auditing and security purposes.
```sql+postgres
select 
  principal_arn,
  service_name,
  last_authenticated,
  age(last_authenticated::date),
  a ->> 'ActionName' as action_name,
  a ->> 'LastAccessedEntity' as action_last_accessed_entity,
  a ->> 'LastAccessedRegion' as action_last_accessed_region,
  a ->> 'LastAccessedTime' as action_last_accessed_time
from 
  aws_iam_access_advisor,
  jsonb_array_elements(tracked_actions_last_accessed) as a
where
  principal_arn = 'arn:aws:iam::123456789123:user/jane'
  and last_authenticated is not null
  and service_namespace = 's3'
order by 
  age asc;
```

```sql+sqlite
select 
  principal_arn,
  service_name,
  last_authenticated,
  julianday('now') - julianday(last_authenticated) as age,
  json_extract(a.value, '$.ActionName') as action_name,
  json_extract(a.value, '$.LastAccessedEntity') as action_last_accessed_entity,
  json_extract(a.value, '$.LastAccessedRegion') as action_last_accessed_region,
  json_extract(a.value, '$.LastAccessedTime') as action_last_accessed_time
from 
  aws_iam_access_advisor,
  json_each(tracked_actions_last_accessed) as a
where
  principal_arn = 'arn:aws:iam::123456789123:user/jane'
  and last_authenticated is not null
  and service_namespace = 's3'
order by 
  age asc;
```


### For all users in the account, find unused services
This query helps you identify the services associated with each user in your account that have not been used. This is particularly useful for optimizing resource allocation and ensuring security by minimizing unnecessary access permissions.

```sql+postgres
select 
  principal_arn,
  service_name
from
  aws_iam_user as u,
  aws_iam_access_advisor as adv
where
  adv.principal_arn = u.arn
  and last_authenticated is null;
```

```sql+sqlite
select 
  principal_arn,
  service_name
from
  aws_iam_user as u,
  aws_iam_access_advisor as adv
where
  adv.principal_arn = u.arn
  and last_authenticated is null;
```