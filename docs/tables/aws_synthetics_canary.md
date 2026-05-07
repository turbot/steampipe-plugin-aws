---
title: "Steampipe Table: aws_synthetics_canary - Query AWS CloudWatch Synthetics Canaries using SQL"
description: "Allows users to query AWS CloudWatch Synthetics Canaries for information about their synthetic monitoring
configuration and status."
folder: "Synthetics"
---

# Table: aws_synthetics_canary - Query AWS CloudWatch Synthetics Canaries using SQL

The AWS CloudWatch Synthetics Canary is a synthetic monitoring service that allows you to create and configure canaries,
which are configurable scripts that run on schedule to monitor endpoints and APIs by simulating actions that your users
may perform.

## Table Usage Guide

The `aws_synthetics_canary` table in Steampipe provides you with information about Synthetics Canaries within AWS
CloudWatch. This table allows you, as a DevOps engineer, to query canary-specific details, including runtime, code
configuration, dry runs, schedule, and status information. You can utilize this table to gather insights on your
canaries, such as the last time it was started or stopped, browser type coverage, and even the engines used to execute
canary runs. The schema outlines the various attributes of the Synthetics Canary for you, including the name, runtime,
engine, schedule, and status.

## Examples

### List of Canaries using the syn-nodejs-puppeteer-15.0 Runtime
Identify canaries using a specific runtime. This could help when upgrading system dependencies.

```sql+postgres
select
    name,
    id
from
    aws_synthetics_canary
where
    runtime_version = 'syn-nodejs-puppeteer-15.0'
```

```sql+sqlite
select
    name,
    id
from
    aws_synthetics_canary
where
    runtime_version = 'syn-nodejs-puppeteer-15.0'
```
