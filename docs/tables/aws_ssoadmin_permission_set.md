---
title: "Steampipe Table: aws_ssoadmin_permission_set - Query AWS SSO Admin Permission Set using SQL"
description: "Allows users to query AWS SSO Admin Permission Set to retrieve data related to the permissions sets of AWS Single Sign-On (SSO) service."
folder: "SSO"
---

# Table: aws_ssoadmin_permission_set - Query AWS SSO Admin Permission Set using SQL

The AWS SSO Admin Permission Set is a component of AWS Single Sign-On (SSO) that defines the level of access that a user or group has to AWS resources. It holds a collection of permissions that can be used to restrict or allow actions on specific resources. AWS SSO Admin Permission Sets can be customized to suit specific needs, ensuring secure and efficient access management across your AWS environment.

## Table Usage Guide

The `aws_ssoadmin_permission_set` table in Steampipe provides you with information about the permission sets associated with AWS Single Sign-On (SSO) service. This table allows you, as a DevOps engineer, to query permission set-specific details, including the permission set name, description, created date, and related metadata. You can utilize this table to gather insights on permission sets, such as the instances of each permission set, associated policies, and more. The schema outlines the various attributes of the permission set for you, including the permission set ARN, created date, session duration, and associated tags.

## Examples

### Basic info
Explore the details of AWS SSO permission sets, including when they were created and their current state. This information can be useful for auditing purposes, understanding access controls, or reviewing the configuration of your AWS environment.

```sql+postgres
select
  name,
  arn,
  created_date,
  description,
  relay_state,
  session_duration,
  tags
from
  aws_ssoadmin_permission_set;
```

```sql+sqlite
select
  name,
  arn,
  created_date,
  description,
  relay_state,
  session_duration,
  tags
from
  aws_ssoadmin_permission_set;
```