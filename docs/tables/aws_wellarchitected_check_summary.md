---
title: "Steampipe Table: aws_wellarchitected_check_summary - Query AWS Well-Architected Tool Check Summary using SQL"
description: "Allows users to query AWS Well-Architected Tool Check Summary for detailed information about the checks for all workloads. This table provides insights into the state of your workloads, highlighting potential risks and areas for improvement."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_check_summary - Query AWS Well-Architected Tool Check Summary using SQL

The AWS Well-Architected Tool is a service that helps you review the state of your workloads and compares them to the latest AWS architectural best practices. The service provides a summary of the check, highlighting high-risk issues where your architecture deviates from AWS best practices. This tool supports you in improving your workloads based on the five pillars of the Well-Architected Framework: Operational Excellence, Security, Reliability, Performance Efficiency, and Cost Optimization.

## Table Usage Guide

The `aws_wellarchitected_check_summary` table in Steampipe provides you with information about the check summaries of all workloads within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query check-specific details, including the workload ID, lens alias, pillar ID, and risk counts. You can utilize this table to gather insights on checks, such as the number of high-risk items, medium-risk items, and checks that are not applicable. The schema outlines the various attributes of the check summary for you, including the workload ID, lens alias, pillar ID, and risk counts.

## Examples

### Basic info
This query is used to gain insights into the summary of checks performed within AWS Well-Architected Tool. It's practical application lies in helping users understand the status of their AWS workloads, making it easier to manage and improve their cloud architectures.

```sql+postgres
select
  id,
  name,
  description,
  jsonb_pretty(account_summary) as account_summary,
  choice_id,
  lens_arn,
  pillar_id,
  question_id,
  status,
  region,
  workload_id
from
  aws_wellarchitected_check_summary;
```

```sql+sqlite
select
  id,
  name,
  description,
  json_pretty(account_summary) as account_summary,
  choice_id,
  lens_arn,
  pillar_id,
  question_id,
  status,
  region,
  workload_id
from
  aws_wellarchitected_check_summary;
```

### Get summarized trusted advisor check report for a workload
This query is useful for gaining insights into the overall health and status of a specific workload in your AWS environment. It can help identify areas of concern or improvement, making it an essential tool for effective workload management and optimization.

```sql+postgres
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  status,
  choice_id,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

```sql+sqlite
select
  workload_id,
  id,
  name,
  account_summary,
  status,
  choice_id,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

### List trusted advisor checks with errors
Explore which Trusted Advisor checks have encountered errors. This is useful to quickly pinpoint areas in your AWS infrastructure that may need immediate attention or remediation.

```sql+postgres
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  status = 'ERROR';
```

```sql+sqlite
select
  workload_id,
  id,
  name,
  account_summary,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  status = 'ERROR';
```

### Get account summary for trusted advisor checks
Determine the areas in which trusted advisor checks may have encountered issues or warnings, aiding in the assessment of the overall health and security of your AWS account.

```sql+postgres
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary;
```

```sql+sqlite
select
  workload_id,
  id,
  name,
  json_extract(account_summary, '$.ERROR') as errors,
  json_extract(account_summary, '$.FETCH_FAILED') as fetch_failed,
  json_extract(account_summary, '$.NOT_AVAILABLE') as not_available,
  json_extract(account_summary, '$.OKAY') as okay,
  json_extract(account_summary, '$.WARNING') as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary;
```

### Get account summary for trusted advisor checks for well-architected lens in a particular workload
This query is designed to pinpoint the specific areas of a given workload that may need attention or improvement as per the AWS Trusted Advisor checks. It helps in assessing the health of the workload under the well-architected lens, providing insights into any errors, warnings, or failed fetches that have occurred.

```sql+postgres
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

```sql+sqlite
select
  workload_id,
  id,
  name,
  json_extract(account_summary, '$.ERROR') as errors,
  json_extract(account_summary, '$.FETCH_FAILED') as fetch_failed,
  json_extract(account_summary, '$.NOT_AVAILABLE') as not_available,
  json_extract(account_summary, '$.OKAY') as okay,
  json_extract(account_summary, '$.WARNING') as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```