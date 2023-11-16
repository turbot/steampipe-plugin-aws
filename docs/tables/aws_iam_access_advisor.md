---
title: "Table: aws_iam_access_advisor - Query AWS IAM Access Advisor using SQL"
description: "Allows users to query AWS IAM Access Advisor to retrieve information about the service last accessed data for IAM entities (users, groups, and roles)."
---

# Table: aws_iam_access_advisor - Query AWS IAM Access Advisor using SQL

The `aws_iam_access_advisor` table in Steampipe provides information about the service last accessed data for IAM entities like users, groups, and roles within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query details about the services that IAM entities can access, the actions they can perform, and when they last accessed the services. Users can utilize this table to gather insights on access patterns, such as identifying unused permissions or verifying least privilege policies. The schema outlines the various attributes of the IAM entity, including the entity ARN, last accessed time, and the services accessible.

**Important Notes:**

- You ***must*** specify a single `principal_arn` in a `where` or `join` clause in order to use this table.  

- The service last accessed data includes all attempts to access an AWS API, not just the successful ones. This includes all attempts that were made using the AWS Management Console, the AWS API through any of the SDKs, or any of the command line tools. An unexpected entry in the service last accessed data does not mean that your account has been compromised, because the request might have been denied. Refer to your CloudTrail logs as the authoritative source for information about all API calls and whether they were successful or denied access. 

- Service last accessed data does not use other policy types when determining whether a resource could access a service. These other policy types include resource-based policies, access control lists, AWS Organizations policies, IAM permissions boundaries, and AWS STS assume role policies. It only applies permissions policy logic. For more about the evaluation of policy types, see Evaluating Policies in the IAM User Guide.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_access_advisor` table, you can use the `.inspect aws_iam_access_advisor` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the IAM entity. This can be used to join this table with other tables that contain IAM entity information.
- `service_name`: The name of the AWS service. This can be used to join this table with other tables that contain AWS service information.
- `last_authenticated`: The date and time, in ISO 8601 date-time format, when the IAM entity last attempted to access the AWS service. This can be used to identify unused permissions or verify least privilege policies.

## Examples

### Show the most recently used services for a user, role, group, or policy
```sql
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

### Show unused services granted a user, role, group, or policy (unused in the last year)
```sql
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
```sql
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

### Show unused services granted to a role, including the policy that grants access and the actions granted  

```sql
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

### Show action-level last accessed info (currently, only supported for S3)
```sql
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


### For all users in the account, find unused services
```sql
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
