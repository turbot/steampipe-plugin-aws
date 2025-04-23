---
title: "Steampipe Table: aws_wafv2_regex_pattern_set - Query AWS WAFv2 Regex Pattern Set using SQL"
description: "Allows users to query AWS WAFv2 Regex Pattern Set data, providing details about the regex pattern sets used in AWS WAFv2 to filter web requests."
folder: "WAFv2"
---

# Table: aws_wafv2_regex_pattern_set - Query AWS WAFv2 Regex Pattern Set using SQL

The AWS WAFv2 Regex Pattern Set is a feature within the AWS Web Application Firewall (WAF) service. It enables users to define a list of regular expressions (regex) that AWS WAF will use to inspect web requests. This tool is essential for identifying and blocking malicious requests based on pattern matching, thereby providing an additional layer of security for your web applications.

## Table Usage Guide

The `aws_wafv2_regex_pattern_set` table in Steampipe provides you with information about Regex Pattern Sets within AWS WAFv2. This table allows you, as a DevOps engineer, to query regex pattern set-specific details, including the ID, name, and the regular expressions included in the set. You can utilize this table to gather insights on the regex patterns, such as the ARN, ID, lock token, and the regular expressions included in the pattern set. The schema outlines the various attributes of the regex pattern set for you, including the ARN, ID, lock token, name, regular expression list, and associated tags.

## Examples

### Basic info
Determine the areas in which specific AWS WAFv2 regex pattern sets are implemented. This helps in understanding the distribution and application of these pattern sets across different regions.

```sql+postgres
select
  name,
  description,
  arn,
  id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set;
```

```sql+sqlite
select
  name,
  description,
  arn,
  id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set;
```


### List global (CloudFront) regex pattern sets
Explore the global pattern sets used in AWS CloudFront to gain insights into the regular expressions being utilized. This is useful for understanding the scope and configuration of your AWS WAFv2, aiding in the optimization and security of your cloud resources.

```sql+postgres
select
  name,
  description,
  arn,
  id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set
where
  scope = 'CLOUDFRONT';
```

```sql+sqlite
select
  name,
  description,
  arn,
  id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set
where
  scope = 'CLOUDFRONT';
```


### List regex pattern sets with a specific regex pattern
Identify instances where specific regex patterns are used in your AWS WAFv2 regex pattern sets. This can help in managing and monitoring your security configurations.

```sql+postgres
select
  name,
  description,
  arn,
  wrps.id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set as wrps,
  jsonb_array_elements_text(regular_expressions) as regex
where
  regex = '^steampipe';
```

```sql+sqlite
select
  name,
  description,
  arn,
  wrps.id,
  scope,
  json_extract(regex.value, '$') as regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set as wrps,
  json_each(regular_expressions) as regex
where
  json_extract(regex.value, '$') = '^steampipe';
```