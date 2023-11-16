---
title: "Table: aws_config_aggregate_authorization - Query AWS Config Aggregate Authorizations using SQL"
description: "Allows users to query AWS Config Aggregate Authorizations, providing vital information about AWS Config rules and their respective authorizations in an aggregated form."
---

# Table: aws_config_aggregate_authorization - Query AWS Config Aggregate Authorizations using SQL

The `aws_config_aggregate_authorization` table in Steampipe provides information about AWS Config Aggregate Authorizations. This table allows DevOps engineers to query authorization-specific details, including the account ID and region that are allowed to aggregate AWS Config rules. Users can utilize this table to gather insights on AWS Config Aggregate Authorizations, such as the permissions and trust policies associated with each authorization, the AWS account that has been granted the authorization, and more. The schema outlines the various attributes of the AWS Config Aggregate Authorization, including the account ID, region, and associated ARN.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_config_aggregate_authorization` table, you can use the `.inspect aws_config_aggregate_authorization` command in Steampipe.

### Key columns:

- `account_id`: This column contains the ID of the AWS account that is allowed to aggregate data. This key column is useful for joining this table with others that contain account-related information.
- `region`: This column holds the region that is allowed to aggregate data. It can be used to join this table with other tables that contain region-specific information.
- `arn`: This column contains the Amazon Resource Name (ARN) of the aggregate authorization. This key column is useful for joining this table with others that contain ARN-related information.

## Examples

### Basic info

```sql
select
  arn,
  authorized_account_id,
  authorized_aws_region,
  creation_time
from
  aws_config_aggregate_authorization;
```
