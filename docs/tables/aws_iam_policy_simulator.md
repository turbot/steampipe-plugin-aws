---
title: "Table: aws_iam_policy_simulator - Query AWS IAM Policy Simulator using SQL"
description: "Allows users to query IAM Policy Simulator for evaluating the effects of IAM access control policies. It provides information such as evaluation results, matching resources, and involved actions."
---

# Table: aws_iam_policy_simulator - Query AWS IAM Policy Simulator using SQL

The `aws_iam_policy_simulator` table in Steampipe provides information about IAM Policy Simulator within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query evaluation results, matching resources, and involved actions. It can be used to understand the effects of IAM access control policies. Users can utilize this table to gather insights on policy simulation, such as the resources involved in the policy, the actions that can be performed, and the evaluation results. The schema outlines the various attributes of the IAM Policy Simulator, including the policy source, policy action, policy resource, and evaluation result.

**Note** that you ***must*** specify a single `action`, `resource_arn`, and `principal_arn` in a where or join clause in order to use this table. 

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_policy_simulator` table, you can use the `.inspect aws_iam_policy_simulator` command in Steampipe.

Key columns:

- `policy_source`: This is the source of the policy. It can be used to identify the origin of the policy.

- `policy_action`: This column provides the action in the policy. It can be used to understand what actions are allowed or denied by the policy.

- `evaluation_result`: This column provides the result of the policy simulation. It can be used to understand the effects of the policy on the resources.

## Examples

### Check if user has s3:DeleteBucket on any resource
```sql
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

```sql
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

### For all users in the account, check whether they have `sts:AssumeRole` on all roles.
```sql
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
