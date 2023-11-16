---
title: "Table: aws_appconfig_application - Query AWS AppConfig Applications using SQL"
description: "Allows users to query AWS AppConfig Applications to gather detailed information about each application, including its name, description, associated environments, and more."
---

# Table: aws_appconfig_application - Query AWS AppConfig Applications using SQL

The `aws_appconfig_application` table in Steampipe provides information about AWS AppConfig Applications. This table allows DevOps engineers and other technical professionals to query application-specific details, including its ID, name, description, and associated environments. Users can utilize this table to gather insights on applications, such as their deployment strategies, associated configurations, and more. The schema outlines the various attributes of the AppConfig application, including the application ID, name, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_appconfig_application` table, you can use the `.inspect aws_appconfig_application` command in Steampipe.

Key columns:

- `id`: The unique ID of the application. It can be used to join this table with other tables to fetch more detailed information.
- `name`: The name of the application. It provides a human-readable reference to the application and can be used to filter results based on application names.
- `description`: The description of the application. It provides additional context about the application's purpose and can be used to filter results based on specific keywords or phrases.

## Examples

### Basic info

```sql
select
  arn,
  id,
  name,
  description,
  tags
from
  aws_appconfig_application;
```