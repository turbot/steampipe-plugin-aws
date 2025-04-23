---
title: "Steampipe Table: aws_cloudwatch_log_resource_policy - Query AWS CloudWatch Log Resource Policies using SQL"
description: "Allows users to query AWS CloudWatch Log Resource Policies, providing details such as the policy name, policy document, and last updated timestamp."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_resource_policy - Query AWS CloudWatch Log Resource Policies using SQL

The AWS CloudWatch Log Resource Policy is a feature of Amazon CloudWatch that allows you to manage resource policies. These policies enable AWS services to perform tasks on your behalf without sharing your security credentials. They are crucial in controlling who can access your logs and what actions they can perform.

## Table Usage Guide

The `aws_cloudwatch_log_resource_policy` table in Steampipe provides you with information about log resource policies within Amazon CloudWatch Logs. This table allows you, as a DevOps engineer, to query policy-specific details, including the policy name, policy document, and last updated timestamp. You can utilize this table to gather insights on policies, such as what actions are allowed or denied, the resources to which the policy applies, and the conditions under which the policy takes effect. The schema outlines for you the various attributes of the CloudWatch Logs resource policy, including the policy name, policy document, and last updated timestamp.

## Examples

### Basic Info
Explore the updates made to your AWS CloudWatch log resource policies. This query can be used to track policy changes over time, ensuring your settings align with your security and operational requirements.

```sql+postgres
select
  policy_name,
  last_updated_time,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_cloudwatch_log_resource_policy;
```

```sql+sqlite
select
  policy_name,
  last_updated_time,
  policy,
  policy_std
from
  aws_cloudwatch_log_resource_policy;
```