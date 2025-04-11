---
title: "Steampipe Table: aws_iam_service_specific_credential - Query AWS IAM Service Specific Credentials using SQL"
description: "Allows users to query AWS IAM Service Specific Credentials, retrieving detailed information about each credential, such as the username, status, creation date, and service name."
folder: "IAM"
---

# Table: aws_iam_service_specific_credential - Query AWS IAM Service Specific Credentials using SQL

The AWS IAM Service Specific Credentials are a type of AWS credentials that allow you to programmatically access AWS services. These credentials are used when the access to services is from an application that is running on an EC2 instance. They provide a secure way to distribute and rotate credentials on multiple EC2 instances without having to manage the underlying AWS keys.

## Table Usage Guide

The `aws_iam_service_specific_credential` table in Steampipe provides you with information about service-specific credentials within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query credential-specific details, including the associated user, status, creation date, and service name. You can utilize this table to gather insights on credentials, such as those associated with a specific user, the status of each credential, and the services for which they are used. The schema outlines the various attributes of service-specific credentials for you, including the username, status, creation date, and service name.

## Examples

### Basic info
Explore which specific AWS IAM services have associated credentials, along with their creation dates and linked user names. This can help in auditing and managing access controls in your AWS environment.

```sql+postgres
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential;
```

```sql+sqlite
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential;
```

### IAM user details for service specific credentials
Discover the segments that are using service-specific credentials in AWS IAM, including details like user names and whether multi-factor authentication is enabled. This query is beneficial for auditing security practices and ensuring adherence to best practices.

```sql+postgres
select
  s.service_name as service_name,
  s.service_specific_credential_id as service_specific_credential_id,
  u.name as user_name,
  u.user_id as user_id,
  u.password_last_used as password_last_used,
  u.mfa_enabled as mfa_enabled
from
  aws_iam_service_specific_credential as s,
  aws_iam_user as u
where
  s.user_name = u.name;
```

```sql+sqlite
select
  s.service_name as service_name,
  s.service_specific_credential_id as service_specific_credential_id,
  u.name as user_name,
  u.user_id as user_id,
  u.password_last_used as password_last_used,
  u.mfa_enabled as mfa_enabled
from
  aws_iam_service_specific_credential as s
join
  aws_iam_user as u
on
  s.user_name = u.name;
```

### Service specific credentials older than 30 days
Determine the areas in which service-specific credentials in AWS IAM are older than 30 days. This can be useful for identifying potential security risks associated with outdated credentials.

```sql+postgres
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential
where
  create_date <= current_date - interval '30' day;
```

```sql+sqlite
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential
where
  create_date <= date('now', '-30 day');
```