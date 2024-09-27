---
title: "Steampipe Table: aws_shield_attack - Query attacks detected by AWS Shield Advanced in the past using SQL"
description: "Allows users to query information about attacks AWS Shield Advanced was able to detect in the past and provide detailed information about the attack."
---

# Table: aws_shield_attack - Query information about AWS Shield Advanced detect attacks using SQL

AWS Shield is a DDoS protection service from AWS. AWS Shield Advanced provide you detailed information about attacks that it was able to detect in the past. This information contains details, such as the start and end time of the attack, the type of attack, the resources that were targeted, the most requested URLs and the mitigation actions that were taken.

## Table Usage Guide

The `aws_shield_attack` table in Steampipe allows you to query AWS Shield Advanced for more details about a DDoS event it was able to detect. This table requires you to specify the `attack_id` of the attack in the WHERE statement of your SQL query for which you want to retrieve information about. For more information about the different columns and their values of this table, please refer to the [official AWS documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeAttack.html#API_DescribeAttack_ResponseSyntax).

## Examples

### Basic info

```sql+postgres
select
  resource_arn,
  start_time,
  end_time,
  jsonb_pretty(sub_resources) as sub_resources,
  jsonb_pretty(attack_properties) as attack_properties,
  mitigations
from
  aws_shield_attack
where
  attack_id = '5dc752eb-96b1-4a9a-8bca-cd3a1096dd56'
```

```sql+sqlite
select
  resource_arn,
  start_time,
  end_time,
  json_pretty(sub_resources) as sub_resources,
  json_pretty(attack_properties) as attack_properties,
  mitigations
from
  aws_shield_attack
where
  attack_id = '5dc752eb-96b1-4a9a-8bca-cd3a1096dd56'
```
