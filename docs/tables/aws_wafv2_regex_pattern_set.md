# Table: aws_wafv2_regex_pattern_set

An AWS WAFv2 regex pattern set contains a set of regex patterns to allow or block web requests if the regex patterns appears in the request.

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
