---
title: "Steampipe Table: aws_shield_attack_statistic - Query AWS Shield Attack Statistics using SQL"
description: "Allows users to query AWS Shield Attack Statistics and retrieve an overview and statistics over past attacks Shield was able detected."
folder: "Shield"
---

# Table: aws_shield_attack_statistic - Query AWS Shield Attack Statistics using SQL

AWS Shield is a DDoS protection service from AWS. AWS Shield provide statistics about the number and type of attacks Shield has detected in the last year for all resources that belong to your account. These statistics are available to you regardless of whether you have subscribed to Shield Advanced or not.

## Table Usage Guide

The `aws_shield_attack_statistic` table in Steampipe allows you to query AWS Shield Attack Statistics and information about layer 3, 4 and 7 attacks that Shield has detected in the last year. It gives you an overview which kind of attacks have the most impact on your resources and how many attacks Shield was able to detect. For more information about the individual columns and their values, please refer to the [official AWS documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeAttackStatistics.html).

## Examples

### Basic info

```sql+postgres
select
  max,
  unit,
  attack_count
from
  aws_shield_attack_statistic
order by
  attack_count desc;
```

```sql+sqlite
select
  max,
  unit,
  attack_count
from
  aws_shield_attack_statistic
order by
  attack_count desc;
```
