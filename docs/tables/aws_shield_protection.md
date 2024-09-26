---
title: "Steampipe Table: aws_shield_protection - Query AWS Shield Advanced Protections using SQL"
description: "Allows users to query AWS Shield Advanced Protections and retrieve detailed information about each protection's settings."
---

# Table: aws_shield_protection - Query AWS Shield Advanced Protections using SQL

AWS Shield Advanced Protections are safeguards provided by AWS to protect AWS resources against Distributed Denial of Service (DDoS) attacks.

## Table Usage Guide

The `aws_shield_protection` table in Steampipe allows you to query AWS Shield Advanced Protections and retrieve detailed information about each protection's settings. This table provides you with insights into the protections that are currently active in your AWS environment, including the ARN of the resource that is protected and the automatic application layer DDoS mitigation setting. You can use this table to monitor the status of your AWS Shield Advanced Protections and ensure that your resources are protected against DDoS attacks.

## Examples

### Basic info

Discover the protections that are currently active in your account and what kind of resources they are protecting.

```sql+postgres
select
  name,
  resource_arn
from
  aws_shield_protection;
```

```sql+sqlite
select
  name,
  resource_arn
from
  aws_shield_protection;
```

### Identify Protections without Route 53 Health Checks

Identify the protections that are missing Route 53 Health Checks. This information can be useful to see which protections still need Health Checks in order to be covered by the proactive engagement of the Shield Response Team.

```sql+postgres
select
  name,
  resource_arn
from
  aws_shield_protection
where
  health_check_ids is null;
```

```sql+sqlite
select
  name,
  resource_arn
from
  aws_shield_protection
where
  health_check_ids is null;
```
