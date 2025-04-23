---
title: "Steampipe Table: aws_iam_access_key - Query AWS IAM Access Keys using SQL"
description: "Allows users to query IAM Access Keys in AWS to obtain details about the access keys associated with an IAM user. This includes the access key ID, status, creation date, and more."
folder: "IAM"
---

# Table: aws_iam_access_key - Query AWS IAM Access Keys using SQL

The AWS Identity and Access Management (IAM) Access Keys are long-term credentials for an IAM user or the AWS account root user. These keys are used in conjunction with the access key ID to cryptographically sign programmatic AWS requests for authentication. Managing access keys appropriately enables you to protect your AWS resources from unauthorized access.

## Table Usage Guide

The `aws_iam_access_key` table in Steampipe provides you with information about IAM Access Keys within AWS Identity and Access Management (IAM). This table lets you, as a DevOps engineer, query access key-specific details, including the access key ID, status, creation date, and more. You can utilize this table to gather insights on access keys, such as their current status (active/inactive), the IAM user they are associated with, and their creation date. The schema outlines the various attributes of the IAM Access Key for you, including the access key ID, status, creation date, and the IAM user to which it belongs.

## Examples

### List of access keys with their corresponding user name and date of creation
Discover the segments that hold information about user access keys, including who created them and when, to help manage and monitor AWS IAM security credentials effectively.

```sql+postgres
select
  access_key_id,
  user_name,
  create_date
from
  aws_iam_access_key;
```

```sql+sqlite
select
  access_key_id,
  user_name,
  create_date
from
  aws_iam_access_key;
```


### List of access keys which are inactive
Determine the areas in which AWS IAM access keys are inactive. This can be useful for identifying unused keys, potentially improving security by reducing the number of active keys in your system.

```sql+postgres
select
  access_key_id,
  user_name,
  status
from
  aws_iam_access_key
where
  status = 'Inactive';
```

```sql+sqlite
select
  access_key_id,
  user_name,
  status
from
  aws_iam_access_key
where
  status = 'Inactive';
```


### Access key count by user name
Determine the number of access keys associated with each user in your AWS IAM service. This can be useful for understanding how access is distributed across your users, potentially highlighting areas where access can be consolidated or better managed.

```sql+postgres
select
  user_name,
  count (access_key_id) as access_key_count
from
  aws_iam_access_key
group by
  user_name;
```

```sql+sqlite
select
  user_name,
  count(access_key_id) as access_key_count
from
  aws_iam_access_key
group by
  user_name;
```