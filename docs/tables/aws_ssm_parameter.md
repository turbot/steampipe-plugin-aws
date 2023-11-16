---
title: "Table: aws_ssm_parameter - Query AWS Systems Manager Parameter Store using SQL"
description: "Allows users to query AWS Systems Manager Parameter Store to retrieve information about parameters, their types, values, and associated metadata."
---

# Table: aws_ssm_parameter - Query AWS Systems Manager Parameter Store using SQL

The `aws_ssm_parameter` table in Steampipe provides information about parameters within AWS Systems Manager Parameter Store. This table allows DevOps engineers to query parameter-specific details, such as parameter names, types, values, and associated metadata. Users can utilize this table to gather insights on parameters, such as parameter descriptions, last modification dates, and the user who last modified the parameter. The schema outlines the various attributes of the parameter, including the parameter ARN, type, value, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_parameter` table, you can use the `.inspect aws_ssm_parameter` command in Steampipe.

### Key columns:

- `name`: The name of the parameter. It can be used to join this table with others that require a parameter name.
- `type`: The type of the parameter. It can be used to filter parameters based on their types.
- `last_modified_date`: The date the parameter was last changed. It can be used to track parameter modifications over time.

## Examples

### SSM parameter basic info

```sql
select
  name,
  type,
  data_type,
  tier,
  region
from
  aws_ssm_parameter;
```


### Policy details of advanced tier ssm parameter

```sql
select
  name,
  tier,
  p ->> 'PolicyType' as policy_type,
  p ->> 'PolicyStatus' as Policy_status,
  p ->> 'PolicyText' as policy_text
from
  aws_ssm_parameter,
  jsonb_array_elements(policies) as p;
```


### List of SSM parameters which do not have owner or app_id tag key

```sql
select
  name
from
  aws_ssm_parameter
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```