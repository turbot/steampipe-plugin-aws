---
title: "Steampipe Table: aws_api_gateway_usage_plan - Query AWS API Gateway Usage Plans using SQL"
description: "Allows users to query AWS API Gateway Usage Plans in order to retrieve information about the usage plans configured in the AWS API Gateway service."
folder: "API Gateway"
---

# Table: aws_api_gateway_usage_plan - Query AWS API Gateway Usage Plans using SQL

The AWS API Gateway Usage Plans are a feature of Amazon API Gateway that allows developers to manage and restrict the usage of their APIs. These plans can be associated with API keys to enable cost recovery, as well as to control the usage of APIs by third-party developers. This ensures a smooth and controlled distribution of your APIs, protecting them from misuse and overuse.

## Table Usage Guide

The `aws_api_gateway_usage_plan` table in Steampipe provides you with information about usage plans within AWS API Gateway. This table allows you, as a DevOps engineer, to query usage plan specific details, including associated API stages, throttle and quota limits, and associated metadata. You can utilize this table to gather insights on usage plans, such as plans with specific rate limits, the number of requests your clients can make per a given period, and more. The schema outlines the various attributes of the usage plan, including the plan ID, name, description, associated API keys, and associated tags for you.

## Examples

### Basic info
Explore the various usage plans associated with your AWS API Gateway. This can help you better manage and monitor your API usage, ensuring optimal performance and cost-effectiveness.

```sql+postgres
select
  name,
  id,
  product_code,
  description,
  api_stages
from
  aws_api_gateway_usage_plan;
```

```sql+sqlite
select
  name,
  id,
  product_code,
  description,
  api_stages
from
  aws_api_gateway_usage_plan;
```


### List the API gateway usage plans where quota ( i.e the number of api call a user can make within a time period) is disabled
Identify instances where the API gateway usage plans do not have a set quota, which indicates that there is no limit to the number of API calls a user can make within a certain time period. This might be useful in understanding potential areas of vulnerability or overuse in your system.

```sql+postgres
select
  name,
  id,
  quota
from
  aws_api_gateway_usage_plan
where
  quota is null;
```

```sql+sqlite
select
  name,
  id,
  quota
from
  aws_api_gateway_usage_plan
where
  quota is null;
```


### List the API gateway usage plan where throttle ( i.e the rate at which user can make request ) is disabled
Determine the areas in which the API gateway usage plan lacks a throttle feature, indicating that there are no restrictions on user request rates.

```sql+postgres
select
  name,
  id,
  throttle
from
  aws_api_gateway_usage_plan
where
  throttle is null;
```

```sql+sqlite
select
  name,
  id,
  throttle
from
  aws_api_gateway_usage_plan
where
  throttle is null;
```