---
title: "Steampipe Table: aws_ecr_registry_scanning_configuration - Query AWS ECR Registry Scanning Configuration using SQL"
description: "Allows users to query AWS ECR Registry Scanning Configuration at the private registry level on a per-region basis."
folder: "Config"
---

# Table: aws_ecr_registry_scanning_configuration - Query AWS ECR Registry Scanning Configuration using SQL

The AWS ECR Registry Scanning Configurations are defined at the private registry level on a per-region basis. These refer to the settings and policies that govern how Amazon ECR scans your container images for vulnerabilities. Amazon ECR integrates with the Amazon ECR image scanning feature, which automatically scans your Docker and OCI images for software vulnerabilities.

## Table Usage Guide

The `aws_ecr_registry_scanning_configuration` table in Steampipe provides you with information about the scanning configurations of Amazon Elastic Container Registry (ECR). This table allows you, as a cloud administrator, security team member, or developer, to query the scanning rules associated with the registry. You can utilize this table to gather insights on scanning configurations, such as the rules, the repository filters, and the region name. The schema outlines the various attributes of the scanning configurations for you, including the region, rules, repository filters, scan type and scan frequency.

## Examples

### Basic configuration info
Analyze the configuration to understand that Amazon ECR scans your container images for vulnerabilities. This is essential for several reasons, primarily centered around security, compliance, and operational efficiency in managing container images.

```sql+postgres
select
  registry_id,
  jsonb_pretty(scanning_configuration),
  region
from
  aws_ecr_registry_scanning_configuration;
```

```sql+sqlite
select
  registry_id,
  scanning_configuration,
  region
from
  aws_ecr_registry_scanning_configuration;
```

### Configuration info for a particular region
Determine the scanning configuration of container images for a specific region. This query is beneficial for understanding the scanning configuration of your container images in that particular region.

```sql+postgres
select
  registry_id,
  jsonb_pretty(scanning_configuration),
  region
from
  aws_ecr_registry_scanning_configuration
where
  region = 'ap-south-1';
```

```sql+sqlite
select
  registry_id,
  scanning_configuration,
  region
from
  aws_ecr_registry_scanning_configuration
where
  region = 'ap-south-1';
```


### List the regions where enhanced scanning is enabled
Identify regions where the enhanced scanning is enabled for container images. This helps determine whether enhanced vulnerability scanning features are available through integrations with AWS services or third-party tools.

```sql+postgres
select
  registry_id,
  region
from
  aws_ecr_registry_scanning_configuration
where
  scanning_configuration ->> 'ScanType' = 'ENHANCED'
```

```sql+sqlite
select
  registry_id,
  region
from
  aws_ecr_registry_scanning_configuration
where
  json_extract(scanning_configuration, '$.ScanType') = 'ENHANCED';
```
