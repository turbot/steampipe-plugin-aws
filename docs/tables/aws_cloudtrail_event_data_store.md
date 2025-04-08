---
title: "Steampipe Table: aws_cloudtrail_event_data_store - Query AWS CloudTrail Event Data using SQL"
description: "Allows users to query AWS CloudTrail Event Data, providing information about API activity in AWS accounts. This includes details about API calls, logins, and other events captured by AWS CloudTrail."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_event_data_store - Query AWS CloudTrail Event Data using SQL

The AWS CloudTrail Event Data is an AWS service that enables governance, compliance, operational auditing, and risk auditing of your AWS account. It allows you to log, continuously monitor, and retain account activity related to actions across your AWS infrastructure. The service provides event history of your AWS account activity, including actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services.

## Table Usage Guide

The `aws_cloudtrail_event_data_store` table in Steampipe provides you with information about API activity in your AWS accounts. This includes details about your API calls, logins, and other events captured by AWS CloudTrail. This table allows you, as a DevOps engineer, to query event-specific details, including event names, event sources, and related metadata. You can utilize this table to gather insights on API activity, such as identifying unusual API calls, tracking login activity, and monitoring changes to your AWS resources. The schema outlines the various attributes of the CloudTrail event for you, including the event ID, event time, event name, and user identity.

## Examples

### Basic info
Explore the status and configuration of your AWS CloudTrail event data stores, including when they were created and their current settings. This can help you maintain security and compliance by ensuring features like multi-region access, organization-wide access, and termination protection are enabled as needed.

```sql+postgres
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store;
```

```sql+sqlite
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store;
```

### List event data stores which are not enabled
Identify instances where event data stores in the AWS CloudTrail service are not enabled. This query is useful in pinpointing potential security vulnerabilities or areas in your system that may not be properly logging and storing event data.

```sql+postgres
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  status <> 'ENABLED';
```

```sql+sqlite
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  status != 'ENABLED';
```

### List event data stores with termination protection disabled
Determine the areas in which event data stores have termination protection disabled in your AWS CloudTrail. This is useful to identify potential vulnerabilities and ensure data safety.

```sql+postgres
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  not termination_protection_enabled;
```

```sql+sqlite
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  termination_protection_enabled = 0;
```