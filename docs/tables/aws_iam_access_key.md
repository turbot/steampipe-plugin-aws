---
title: "Table: aws_iam_access_key - Query AWS IAM Access Keys using SQL"
description: "Allows users to query IAM Access Keys in AWS to obtain details about the access keys associated with an IAM user. This includes the access key ID, status, creation date, and more."
---

# Table: aws_iam_access_key - Query AWS IAM Access Keys using SQL

The `aws_iam_access_key` table in Steampipe provides information about IAM Access Keys within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query access key-specific details, including the access key ID, status, creation date, and more. Users can utilize this table to gather insights on access keys, such as their current status (active/inactive), the IAM user they are associated with, and their creation date. The schema outlines the various attributes of the IAM Access Key, including the access key ID, status, creation date, and the IAM user to which it belongs.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_access_key` table, you can use the `.inspect aws_iam_access_key` command in Steampipe.

### Key columns:

- `access_key_id`: The ID of the access key. This can be used to join with other tables that contain information about access keys.
- `user_name`: The name of the IAM user that the access key is associated with. This can be used to join with other tables that contain user information.
- `status`: The status of the access key (Active/Inactive). This can be used to filter or sort the access keys based on their status.

## Examples

### List of access keys with their corresponding user name and date of creation

```sql
select
  access_key_id,
  user_name,
  create_date
from
  aws_iam_access_key;
```


### List of access keys which are inactive

```sql
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

```sql
select
  user_name,
  count (access_key_id) as access_key_count
from
  aws_iam_access_key
group by
  user_name;
```