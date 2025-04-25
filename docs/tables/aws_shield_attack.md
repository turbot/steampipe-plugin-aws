---
title: "Steampipe Table: aws_shield_attack - Query attacks detected by AWS Shield Advanced in the past using SQL"
description: "Allows users to query information about attacks AWS Shield Advanced was able to detect in the past and provide detailed information about the attack."
folder: "Shield"
---

# Table: aws_shield_attack - Query information about AWS Shield Advanced detect attacks using SQL

AWS Shield is a DDoS protection service from AWS. AWS Shield Advanced provide you detailed information about attacks that it was able to detect in the past. This information contains details, such as the start and end time of the attack, the type of attack, the resources that were targeted, the most requested URLs and the mitigation actions that were taken.

## Table Usage Guide

The `aws_shield_attack` table in Steampipe allows you to query AWS Shield Advanced for more details about a DDoS event it was able to detect. For more information about the different columns and their values of this table, please refer to the AWS Shield Advanced documentation of the [ListAttacks](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_ListAttacks.html#API_ListAttacks_ResponseSyntax) and [DescribeAttack](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeAttack.html#API_DescribeAttack_ResponseSyntax) API.

## Examples

### List all attacks detected by AWS Shield Advanced in the last 30 days

```sql+postgres
select
  resource_arn,
  start_time,
  end_time
from
  aws_shield_attack
where
  start_time between current_date - interval '30 day' and current_date;
```

```sql+sqlite
select
  resource_arn,
  start_time,
  end_time
from
  aws_shield_attack
where
  start_time between date('now', '-30 day') and date('now');
```

### List the most attacked resources of the last 30 days

```sql+postgres
select
  resource_arn,
  count(*) as attacks
from
  aws_shield_attack
where
  start_time between current_date - interval '30 day' and current_date
group by
  resource_arn
order by
  attacks
desc;
```

```sql+sqlite
select
  resource_arn,
  count(*) as attacks
from
  aws_shield_attack
where
  start_time between date('now', '-30 day') and date('now')
group by
  resource_arn
order by
  attacks
desc;
```

### List countries from which the most requests of the attacks of the last 30 days originated

```sql+postgres
select
  top_contributor ->> 'Name' as country,
  sum(cast(top_contributor ->> 'Value' as integer)) as requests
from
  aws_shield_attack,
  jsonb_array_elements(attack_properties) as attack_property,
  jsonb_array_elements(attack_property -> 'TopContributors') as top_contributor
where
  start_time between current_date - interval '30 day' and current_date
  and attack_property ->> 'AttackPropertyIdentifier' = 'SOURCE_COUNTRY'
group by
  country
order by
  requests
desc;
```

```sql+sqlite
select
  top_contributor -> 'Name' as country,
  sum(cast(top_contributor -> 'Value' as integer)) as requests
from
  aws_shield_attack,
  json_each(attack_properties) as attack_property,
  json_each(attack_property -> 'TopContributors') as top_contributor
where
  start_time between date('now', '-30 day') and date('now')
  and attack_property_value_value.key = 'AttackPropertyIdentifier'
  and attack_property_value_value.value = 'SOURCE_COUNTRY'
group by
  country
order by
  requests
desc;
```
