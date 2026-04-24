---
title: "Steampipe Table: aws_synthetics_canary_run - Query AWS CloudWatch Synthetics Canary Runs using SQL"
description: "Allows users to query AWS CloudWatch Synthetics Canary Runs for information about their run configuration and status."
folder: "Synthetics"
---

# Table: aws_synthetics_canary_run - Query AWS CloudWatch Synthetics Canary Runs using SQL

The AWS CloudWatch Synthetics Canary Run is a synthetic monitoring service that allows you to execute canary runs that
simulates user actions to monitor your endpoints and APIs.

## Table Usage Guide

The `aws_synthetics_canary_run` table in Steampipe provides you with information about Synthetics Canary Runs within AWS
CloudWatch. This table allows you, as a DevOps engineer, to query canary run details, including execution timestamps,
browser types, and status information. You can utilize this table to gather insights on the health of your endpoints
and APIs that have been setup to be monitored by canaries. The schema outlines the various attributes of the Synthetics
Canary Run for you, including the canary name, timeline, and status.

**Important Notes**
- You must specify a `name` in a where or join clause in order to use this table.
- `dry_run_id` must be specified together with `run_type = 'DRY_RUN'` or an exception will be raised. See
  https://docs.aws.amazon.com/AmazonSynthetics/latest/APIReference/API_GetCanaryRuns.html for more information.

## Examples

### List Runs of a Canary
Query run information for a specific canary.

```sql+postgres
select
  name,
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
```

```sql+sqlite
select
  name,
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
```

### List Failed Canary Runs
Query run information for failed canary runs. This could be helpful to identify potential health issues in your systems.

```sql+postgres
select
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
  and status ->> 'State' = 'FAILED'
```

```sql+sqlite
select
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
  and status ->> 'State' = 'FAILED'
```

### List Canary Dry Runs
Query information for canary dry runs. All rows returned are dry runs when `run_type = 'DRY_RUN''.

```sql+postgres
select
  name,
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
  and run_type = 'DRY_RUN'
```

```sql+sqlite
select
  name,
  id,
  status
from
  aws_synthetics_canary_run
where
  name = 'TargetCanary'
  and run_type = 'DRY_RUN'
```
