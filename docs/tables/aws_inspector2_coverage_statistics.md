# Table: aws_inspector2_coverage_statistics

AWS Inspector coverage statistic refers to the measurement or analysis of the extent of coverage provided by AWS Inspector for assessing the security and compliance of AWS resources. It provides insights into the assessment scope and the level of evaluation performed by AWS Inspector.


## Examples

### Basic info

```sql
select
  total_counts,
  counts_by_group
from
  aws_inspector2_coverage_statistics;
```

### Get the count of resources within a group

```sql
select
  g ->> 'Count' as count,
  g ->> 'GroupKey' as group_key
from
  aws_inspector2_coverage_statistics,
  jsonb_array_elements(counts_by_group) as g;
```