---
title: "Steampipe Table: aws_transfer_user - Query AWS Transfer for users in a server using SQL"
description: "Allows users to query AWS Transfer for users in a server, equivalent to list/describe user functions."
folder: "Transfer Family"
---

# Table: aws_transfer_user - Query AWS Transfer Family users using SQL

AWS Transfer Family is a secure transfer service that enables you to transfer files into and out of AWS storage services.

## Table Usage Guide

The `aws_transfer_user` table in Steampipe provides you with information about users inside defined servers in the AWS Transfer Family service. This table allows you, as a DevOps engineer, to query user-specific details, including home directories, ssh keys, usernames and IAM roles.

## Examples

### Basic info
Explore which AWS Transfer users are defined in a server

```sql+postgres
select
  arn,
  server_id,
  user_name
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx';
```

```sql+sqlite
select
  arn,
  server_id,
  user_name
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx';
```

### Sort users descending by SSH public key count

```sql+postgres
select
  arn,
  server_id,
  user_name,
  ssh_public_key_count
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx'
order by
  ssh_public_key_count desc;
```

```sql+sqlite
select
  arn,
  server_id,
  user_name,
  ssh_public_key_count
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx'
order by
  ssh_public_key_count desc;
```

### Get home directory mappings for users

```sql+potgres
select
  server_id,
  user_name,
  home_directory_mappings->0->>'Entry' as entry_home_directory,
  home_directory_mappings->0->>'Target' as target_home_directory
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx';
```

```sql+sqlite
select
  server_id,
  user_name,
  json_extract(home_directory_mappings, '$[0].Entry') as entry_home_directory,
  json_extract(home_directory_mappings, '$[0].Target') as target_home_directory
from
  aws_transfer_user
where
  server_id = 's-xxxxxxxxxxxxxxxxx';
```


### Find user_name across multiple servers

```sql+postgres
select
  server_id,
  user_name,
  arn
from
  aws_transfer_user
where
  server_id in (select server_id from aws_transfer_server)
and
  user_name = 'my_user_to_search';
```

```sql+sqlite
select
  server_id,
  user_name,
  arn
from
  aws_transfer_user
where
  server_id in (select server_id from aws_transfer_server)
and
  user_name = 'my_user_to_search';
```

### Count users by server_id descending

```sql+postgres
select
  count(*) as total_users,
  server_id
from
  aws_transfer_user
group by
  server_id
order by
  total_users desc;
```

```sql+sqlite
select
  count(*) as total_users,
  server_id
from
  aws_transfer_user
group by
  server_id
order by
  total_users desc;
```
