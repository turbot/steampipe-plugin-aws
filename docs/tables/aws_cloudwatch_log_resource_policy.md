---
title: "Table: aws_cloudwatch_log_resource_policy - Query AWS CloudWatch Log Resource Policies using SQL"
description: "Allows users to query AWS CloudWatch Log Resource Policies, providing details such as the policy name, policy document, and last updated timestamp."
---

# Table: aws_cloudwatch_log_resource_policy - Query AWS CloudWatch Log Resource Policies using SQL

The `aws_cloudwatch_log_resource_policy` table in Steampipe provides information about log resource policies within Amazon CloudWatch Logs. This table allows DevOps engineers to query policy-specific details, including the policy name, policy document, and last updated timestamp. Users can utilize this table to gather insights on policies, such as what actions are allowed or denied, the resources to which the policy applies, and the conditions under which the policy takes effect. The schema outlines the various attributes of the CloudWatch Logs resource policy, including the policy name, policy document, and last updated timestamp.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_resource_policy` table, you can use the `.inspect aws_cloudwatch_log_resource_policy` command in Steampipe.

**Key columns**:

- `policy_name`: The name of the policy. This can be used to join with other tables that reference CloudWatch Logs resource policies by name.
- `policy_document`: The policy document, which details what actions are allowed or denied, the resources to which the policy applies, and the conditions under which the policy takes effect. This can be useful for analyzing the permissions granted by the policy.
- `last_updated_timestamp`: The timestamp indicating when the policy was last updated. This can be useful for tracking changes to the policy over time.

## Examples

### Basic Info

```sql
select
  policy_name,
  last_updated_time,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_cloudwatch_log_resource_policy;
```
