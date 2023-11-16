---
title: "Table: aws_wafv2_regex_pattern_set - Query AWS WAFv2 Regex Pattern Set using SQL"
description: "Allows users to query AWS WAFv2 Regex Pattern Set data, providing details about the regex pattern sets used in AWS WAFv2 to filter web requests."
---

# Table: aws_wafv2_regex_pattern_set - Query AWS WAFv2 Regex Pattern Set using SQL

The `aws_wafv2_regex_pattern_set` table in Steampipe provides information about Regex Pattern Sets within AWS WAFv2. This table allows DevOps engineers to query regex pattern set-specific details, including the ID, name, and the regular expressions included in the set. Users can utilize this table to gather insights on the regex patterns, such as the ARN, ID, lock token, and the regular expressions included in the pattern set. The schema outlines the various attributes of the regex pattern set, including the ARN, ID, lock token, name, regular expression list, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafv2_regex_pattern_set` table, you can use the `.inspect aws_wafv2_regex_pattern_set` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the regex pattern set. This can be used to join with other tables that contain AWS ARN data.
- `id`: The unique identifier for the regex pattern set. This is useful for joining with tables that require the regex pattern set ID.
- `name`: The name of the regex pattern set. This is important for identifying the specific regex pattern set when joining with other tables.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  regular_expressions,
  region
from
  aws_wafv2_regex_pattern_set,
  jsonb_array_elements_text(regular_expressions) as regex
where
  regex = '^steampipe';
```
