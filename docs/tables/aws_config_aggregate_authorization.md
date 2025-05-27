---
title: "Steampipe Table: aws_config_aggregate_authorization - Query AWS Config Aggregate Authorizations using SQL"
description: "Allows users to query AWS Config Aggregate Authorizations, providing vital information about AWS Config rules and their respective authorizations in an aggregated form."
folder: "Config"
---

# Table: aws_config_aggregate_authorization - Query AWS Config Aggregate Authorizations using SQL

The AWS Config Aggregate Authorization is a feature of AWS Config that allows you to authorize the aggregator account to collect AWS Config data from source accounts. It simplifies compliance auditing by enabling you to collect configuration and compliance data across multiple accounts and regions, and aggregate it into a central account. This centralized data can then be accessed using SQL queries for analysis and reporting.

## Table Usage Guide

The `aws_config_aggregate_authorization` table in Steampipe provides you with information about AWS Config Aggregate Authorizations. This table allows you, as a DevOps engineer, to query authorization-specific details, including the account ID and region that are allowed to aggregate AWS Config rules. You can utilize this table to gather insights on AWS Config Aggregate Authorizations, such as the permissions and trust policies associated with each authorization, the AWS account that has been granted the authorization, and more. The schema outlines the various attributes of the AWS Config Aggregate Authorization for you, including the account ID, region, and associated ARN.

## Examples

### Basic info
Discover the segments that are authorized to access your AWS configuration data, including the region and account details. This can help you manage access control and understand when these authorizations were created.

```sql+postgres
select
  arn,
  authorized_account_id,
  authorized_aws_region,
  creation_time
from
  aws_config_aggregate_authorization;
```

```sql+sqlite
select
  arn,
  authorized_account_id,
  authorized_aws_region,
  creation_time
from
  aws_config_aggregate_authorization;
```