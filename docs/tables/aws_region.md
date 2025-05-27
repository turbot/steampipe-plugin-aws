---
title: "Steampipe Table: aws_region - Query AWS Region using SQL"
description: "Allows users to query AWS Region to retrieve details about AWS regions including their names, descriptions, and statuses."
folder: "Region"
---

# Table: aws_region - Query AWS Region using SQL

The AWS Region is a geographical area that represents a physical location around the world where AWS clusters data centers. Each AWS Region is designed to be completely isolated from the other AWS Regions, which aids in achieving the greatest possible fault tolerance and stability. This isolation ensures that user data is not replicated between AWS Regions unless explicitly done so by the user.

## Table Usage Guide

The `aws_region` table in Steampipe provides you with information about regions within AWS. This table allows you, as a DevOps engineer, to query region-specific details, including the region name, description, and status. You can utilize this table to gather insights on regions, such as their geographical distribution, operational status, and more. The schema outlines the various attributes of the AWS region for you, including the region name, endpoint, and whether the region is opt-in status.

## Examples

### AWS region info
Determine the areas in which your AWS services are deployed and their opt-in statuses. This can help you manage your resources more effectively, particularly for services that require manual opt-in.

```sql+postgres
select
  name,
  opt_in_status
from
  aws_region;
```

```sql+sqlite
select
  name,
  opt_in_status
from
  aws_region;
```


### List of AWS regions that are enabled
Discover the segments that are not currently active in your AWS regions. This can help you understand which regions are not utilized, potentially highlighting areas for infrastructure optimization or cost savings.

```sql+postgres
select
  name,
  opt_in_status
from
  aws_region
where
  opt_in_status = 'not-opted-in';
```

```sql+sqlite
select
  name,
  opt_in_status
from
  aws_region
where
  opt_in_status = 'not-opted-in';
```
