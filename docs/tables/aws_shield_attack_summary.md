---
title: "Steampipe Table: aws_shield_attack_summary - Query AWS Shield Advanced Attack Summaries using SQL"
description: "Allows users to query AWS Shield Advanced Attack Summaries and get an overview of the attacks that have been detected by AWS Shield Advanced."
---

# Table: aws_shield_attack_summary - Query AWS Shield Advanced Attack Summaries using SQL

AWS Shield is a DDoS protection service from AWS. AWS Shield Advanced provide you detailed information about attacks that it was able to detect in the past. This information contains details, such as the start and end time of the attack, the type of attack, the resources that were targeted, the most requested URLs and the mitigation actions that were taken.

## Table Usage Guide

The `aws_shield_attack_summary` table in Steampipe allows you to get an overview of attacks that Shield Advanced has detected in the past. This table provides you with insights into the number of the attacks, the attacked resources, the attack vector, as well as the start and end time of the attack. For more details about the individual fields, please refer to the [AWS Shield Advanced API documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_ListAttacks.html#API_ListAttacks_ResponseSyntax).

## Examples

### Basic info

```sql+postgres
select
  attack_id,
  resource_arn,
  start_time,
  end_time,
  attack_vectors
from
  aws_shield_attack_summary;
```

```sql+sqlite
select
  attack_id,
  resource_arn,
  start_time,
  end_time,
  attack_vectors
from
  aws_shield_attack_summary;
```
