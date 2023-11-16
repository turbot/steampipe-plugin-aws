---
title: "Table: aws_iam_service_specific_credential - Query AWS IAM Service Specific Credentials using SQL"
description: "Allows users to query AWS IAM Service Specific Credentials, retrieving detailed information about each credential, such as the username, status, creation date, and service name."
---

# Table: aws_iam_service_specific_credential - Query AWS IAM Service Specific Credentials using SQL

The `aws_iam_service_specific_credential` table in Steampipe provides information about service-specific credentials within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query credential-specific details, including the associated user, status, creation date, and service name. Users can utilize this table to gather insights on credentials, such as those associated with a specific user, the status of each credential, and the services for which they are used. The schema outlines the various attributes of service-specific credentials, including the username, status, creation date, and service name.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_service_specific_credential` table, you can use the `.inspect aws_iam_service_specific_credential` command in Steampipe.

**Key columns**:

- `user_name`: The name of the IAM user associated with the service-specific credential. This column is useful for joining with other tables that contain user-specific information.
- `status`: The status of the service-specific credential (Active or Inactive). This column can be used to filter or sort the results based on the status of the credentials.
- `service_name`: The name of the service that the credentials are associated with. This column can be used to join with other tables that contain service-specific information.

## Examples

### Basic info

```sql
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential;
```

### IAM user details for service specific credentials

```sql
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

### Service specific credentials older than 30 days

```sql
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
