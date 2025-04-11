---
title: "Steampipe Table: aws_iam_policy_simulator - Query AWS IAM Policy Simulator using SQL"
description: "Allows users to query IAM Policy Simulator for evaluating the effects of IAM access control policies. It provides information such as evaluation results, matching resources, and involved actions."
folder: "IAM"
---

# Table: aws_iam_policy_simulator - Query AWS IAM Policy Simulator using SQL

The AWS IAM Policy Simulator is a tool that enables you to understand, test, and validate the effects of access control policies. It allows you to simulate how IAM policies and resource-based policies work together to grant or deny access to AWS resources. This helps you to ensure that your policies provide the appropriate permissions before you commit them into production.

## Table Usage Guide

The `aws_iam_policy_simulator` table in Steampipe provides you with information about IAM Policy Simulator within AWS Identity and Access Management (IAM). This table enables you to query evaluation results, matching resources, and involved actions as a DevOps engineer. You can use it to understand the effects of IAM access control policies. You can utilize this table to gather insights on policy simulation, such as the resources involved in the policy, the actions that can be performed, and the evaluation results. The schema outlines the various attributes of the IAM Policy Simulator for you, including the policy source, policy action, policy resource, and evaluation result.

**Important Notes**
- You must specify a single `action`, `resource_arn`, and `principal_arn` in a where or join clause in order to use this table.

## Examples

### Check if user has s3:DeleteBucket on any resource
Determine if a specific user has the ability to delete any bucket in the S3 service. This is useful for auditing user permissions and ensuring that sensitive operations are restricted to authorized individuals.

```sql+postgres
select
  decision
from
  aws_iam_policy_simulator
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```

```sql+sqlite
select
  decision
from
  aws_iam_policy_simulator
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```


### Check if user has 'ec2:terminateinstances' on any resource including details of any policy granting or denying access
Determine if a specific user has the ability to terminate any instances on your AWS EC2 service. This query is useful for identifying potential security risks and ensuring appropriate permissions are in place.

```sql+postgres
select
  decision,
  jsonb_pretty(matched_statements)
from
  aws_iam_policy_simulator
where
  action = 'ec2:terminateinstances'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```

```sql+sqlite
select
  decision,
  json_pretty(matched_statements)
from
  aws_iam_policy_simulator
where
  action = 'ec2:terminateinstances'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```

### For all users in the account, check whether they have `sts:AssumeRole` on all roles.
Determine the areas in which users have the ability to assume all roles within an account. This is useful for identifying potential security risks and ensuring appropriate access controls are in place.

```sql+postgres
select
  u.name,
  decision
from
  aws_iam_policy_simulator p,
  aws_iam_user u
where
  action = 'sts:AssumeRole'
  and resource_arn = '*'
  and p.principal_arn = u.arn;
```

```sql+sqlite
select
  u.name,
  decision
from
  aws_iam_policy_simulator p,
  aws_iam_user u
where
  action = 'sts:AssumeRole'
  and resource_arn = '*'
  and p.principal_arn = u.arn;
```