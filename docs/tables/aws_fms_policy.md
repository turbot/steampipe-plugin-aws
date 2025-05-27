---
title: "Steampipe Table: aws_fms_policy - Query AWS Firewall Manager Policies using SQL"
description: "Allows users to query AWS Firewall Manager Policies using SQL. This table provides information about each AWS Firewall Manager (FMS) policy in an AWS account. It can be used to gain insights into policy details such as the policy name, ID, resource type, security service type, and more."
folder: "Firewall Manager (FMS)"
---

# Table: aws_fms_policy - Query AWS Firewall Manager Policies using SQL

The AWS Firewall Manager Policies is a feature of AWS Firewall Manager, a security management service designed to centrally configure and manage firewall rules across your accounts and applications in AWS. It enables you to easily apply AWS WAF, AWS Shield Advanced, and VPC security group rules across your AWS resources. This centralized control effectively helps in maintaining a consistent security posture across your entire AWS environment.

## Table Usage Guide

The `aws_fms_policy` table in Steampipe provides you with information about each AWS Firewall Manager (FMS) policy in your AWS account. This table allows you, as a DevOps engineer, security professional, or other user, to query policy-specific details, including policy ID, policy name, resource type, security service type, and more. You can utilize this table to gather insights on policies, such as identifying which resources are protected by which policies, understanding the type of security service provided by each policy, and knowing the status of each policy. The schema outlines the various attributes of the FMS policy for you, including the policy ID, policy name, resource tags, remediation enabled status, and more.

## Examples

### Basic info
Explore the details of your AWS Firewall Manager policies to understand their configuration and purpose. This can be beneficial for maintaining security standards and managing resources effectively.

```sql+postgres
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type
from
  aws_fms_policy;
```

```sql+sqlite
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type
from
  aws_fms_policy;
```

### List policies that has remediation enabled
Identify instances where specific policies have remediation enabled to better manage potential security threats and vulnerabilities within your AWS environment.

```sql+postgres
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type,
  remediation_enabled
from
  aws_fms_policy
where
  remediation_enabled;
```

```sql+sqlite
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type,
  remediation_enabled
from
  aws_fms_policy
where
  remediation_enabled = 1;
```

### Count policies by resource type
Discover which policies are being applied and how frequently, across various resource types within your AWS Firewall Manager, providing insights into your security configurations and their distribution.

```sql+postgres
select
  policy_name,
  resource_type,
  count(policy_id) as policy_applied
from
  aws_fms_policy
group by
  policy_name,
  resource_type;
```

```sql+sqlite
select
  policy_name,
  resource_type,
  count(policy_id) as policy_applied
from
  aws_fms_policy
group by
  policy_name,
  resource_type;
```

### List policies that are not active
Identify instances where certain policies are inactive to ensure all necessary measures are being enforced for optimal security management.

```sql+postgres
select
  policy_name,
  policy_id,
  policy_status
from
  aws_fms_policy
where
  policy_status <> 'ACTIVE';
```

```sql+sqlite
select
  policy_name,
  policy_id,
  policy_status
from
  aws_fms_policy
where
  policy_status <> 'ACTIVE';
```