---
title: "Steampipe Table: aws_appsync_graphql_api - Query AWS AppSync GraphQL API using SQL"
description: "Allows users to query AppSync GraphQL APIs to retrieve detailed information about each individual GraphQL API."
folder: "AppSync"
---

# Table: aws_appsync_graphql_api - Query AWS AppSync GraphQL APIs using SQL

AWS AppSync is a fully managed service provided by Amazon Web Services (AWS) that simplifies the development of scalable and secure GraphQL APIs. GraphQL is a query language for APIs that allows clients to request only the data they need, making it more efficient and flexible compared to traditional REST APIs.

## Table Usage Guide

The `aws_appsync_graphql_api` table in Steampipe provides you with information about GraphQL API within AWS Athena. This table allows you, as a data analyst or developer, to GraphQL API specific details, including authentication type, owner of the API, and log configuration details of the API.

## Examples

### List all merged APIs
A merged GraphQL API typically refers to a GraphQL API that aggregates or combines data from multiple sources into a single, unified GraphQL schema. This approach is often used to create a single, cohesive interface for clients, even when the underlying data comes from different services, databases, or microservices.

```sql+postgres
select
  name,
  api_id,
  arn,
  api_type,
  authentication_type,
  owner,
  owner_contact
from
  aws_appsync_graphql_api
where
  api_type = 'MERGED';
```

```sql+sqlite
select
  name,
  api_id,
  arn,
  api_type,
  authentication_type,
  owner,
  owner_contact
from
  aws_appsync_graphql_api
where
  api_type = 'MERGED';
```

### List public APIs of the current account
A public AppSync GraphQL API is accessible over the internet, and clients outside of your AWS account can make requests to it. Public APIs are typically configured with an authentication mechanism to control and secure access. Common authentication methods include API keys and OpenID Connect (OIDC) integration with an identity provider.

```sql+postgres
select
  name,
  api_id,
  api_type,
  visibility
from
  aws_appsync_graphql_api
where
  visibility = 'GLOBAL'
  and owner = account_id;
```

```sql+sqlite
select
  name,
  api_id,
  api_type,
  visibility
from
  aws_appsync_graphql_api
where
  visibility = 'GLOBAL'
  and owner = account_id;
```

### Get the log configuration details of APIs
Discover the queries that have the longest execution times to identify potential areas for performance optimization and enhance the efficiency of your AWS Athena operations.

```sql+postgres
select
  name,
  api_id,
  owner,
  log_config ->> 'CloudWatchLogsRoleArn' as cloud_watch_logs_role_arn,
  log_config ->> 'FieldLogLevel' as field_log_level,
  log_config ->> 'ExcludeVerboseContent' as exclude_verbose_content
from
  aws_appsync_graphql_api;
```

```sql+sqlite
select
  name,
  api_id,
  owner,
  json_extract(log_config, '$.CloudWatchLogsRoleArn') as cloud_watch_logs_role_arn,
  json_extract(log_config, '$.FieldLogLevel') as field_log_level,
  json_extract(log_config, '$.ExcludeVerboseContent') as exclude_verbose_content
from
  aws_appsync_graphql_api;
```