---
title: "Steampipe Table: aws_cloudformation_stack - Query AWS CloudFormation Stack using SQL"
description: "Allows users to query AWS CloudFormation Stack data, including stack name, status, creation time, and associated tags."
folder: "CloudFormation"
---

# Table: aws_cloudformation_stack - Query AWS CloudFormation Stack using SQL

The AWS CloudFormation Stack is a service that allows you to manage and provision AWS resources in an orderly and predictable fashion. You can use AWS CloudFormation to leverage AWS products such as Amazon EC2, Amazon Elastic Block Store, Amazon SNS, Elastic Load Balancing, and Auto Scaling to build highly reliable, highly scalable, cost-effective applications without creating or configuring the underlying AWS infrastructure. With CloudFormation, you describe your desired resources in a template, and AWS CloudFormation takes care of provisioning and configuring those resources for you.

## Table Usage Guide

The `aws_cloudformation_stack` table in Steampipe provides you with information about stacks within AWS CloudFormation. This table enables you as a DevOps engineer to query stack-specific details, including stack name, status, creation time, and associated tags. You can utilize this table to gather insights on stacks, such as stack status, stack resources, stack capabilities, and more. The schema outlines the various attributes of the CloudFormation stack for you, including stack ID, stack name, creation time, stack status, and associated tags.

## Examples

### Find the status of each cloudformation stack
Explore the current status of each AWS CloudFormation stack to monitor the health and progress of your infrastructure deployments. This can help in identifying any potential issues or failures in your stack deployments.

```sql+postgres
select
  name,
  id,
  status
from
  aws_cloudformation_stack;
```

```sql+sqlite
select
  name,
  id,
  status
from
  aws_cloudformation_stack;
```


### List of cloudformation stack where rollback is disabled
Discover the segments that have disabled rollback in their AWS CloudFormation stacks. This can be useful for identifying potential risk areas, as these stacks will not automatically revert to a previous state if an error occurs during stack operations.

```sql+postgres
select
  name,
  disable_rollback
from
  aws_cloudformation_stack
where
  disable_rollback;
```

```sql+sqlite
select
  name,
  disable_rollback
from
  aws_cloudformation_stack
where
  disable_rollback = 1;
```


### List of stacks where termination protection is not enabled
Discover the segments that have not enabled termination protection in their stacks. This is crucial to identify potential risk areas and ensure the safety of your resources.

```sql+postgres
select
  name,
  enable_termination_protection
from
  aws_cloudformation_stack
where
  not enable_termination_protection;
```

```sql+sqlite
select
  name,
  enable_termination_protection
from
  aws_cloudformation_stack
where
  enable_termination_protection = 0;
```


### Rollback configuration info for each cloudformation stack
Explore the settings of your AWS CloudFormation stacks to understand their rollback configurations, including how long they monitor for signs of trouble and what triggers a rollback. This can help optimize your stack management by adjusting these settings based on your operational needs.

```sql+postgres
select
  name,
  rollback_configuration ->> 'MonitoringTimeInMinutes' as monitoring_time_in_min,
  rollback_configuration ->> 'RollbackTriggers' as rollback_triggers
from
  aws_cloudformation_stack;
```

```sql+sqlite
select
  name,
  json_extract(rollback_configuration, '$.MonitoringTimeInMinutes') as monitoring_time_in_min,
  json_extract(rollback_configuration, '$.RollbackTriggers') as rollback_triggers
from
  aws_cloudformation_stack;
```


### Resource ARNs where notifications about stack actions will be sent
Determine the areas in which notifications related to stack actions will be sent. This is useful for managing and tracking changes in your AWS CloudFormation stacks.

```sql+postgres
select
  name,
  jsonb_array_elements_text(notification_arns) as resource_arns
from
  aws_cloudformation_stack;
```

```sql+sqlite
select
  name,
  json_extract(json_each.value, ') as resource_arns
from
  aws_cloudformation_stack,
  json_each(notification_arns);
```