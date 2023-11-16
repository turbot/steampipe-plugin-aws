---
title: "Table: aws_region - Query AWS Region using SQL"
description: "Allows users to query AWS Region to retrieve details about AWS regions including their names, descriptions, and statuses."
---

# Table: aws_region - Query AWS Region using SQL

The `aws_region` table in Steampipe provides information about regions within AWS. This table allows DevOps engineers to query region-specific details, including the region name, description, and status. Users can utilize this table to gather insights on regions, such as their geographical distribution, operational status, and more. The schema outlines the various attributes of the AWS region, including the region name, endpoint, and whether the region is opt-in status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_region` table, you can use the `.inspect aws_region` command in Steampipe.

**Key columns**:

- `name`: The name of the region. This can be used to join with other tables where region name is a foreign key.
- `description`: The description of the region. This provides context about the region which can be useful when joining with other tables.
- `opt_in_status`: The opt-in status of the region. This can be used to filter regions based on their opt-in status.

## Examples

### AWS region info

```sql
select
  name,
  opt_in_status
from
  aws_region;
```


### List of AWS regions which are enable

```sql
select
  name,
  opt_in_status
from
  aws_region
where
  opt_in_status = 'not-opted-in';
```
