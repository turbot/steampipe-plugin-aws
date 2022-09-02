# Table: aws_securityhub_finding_aggregator

AWS Security Hub now allows you to designate an aggregation Region and link some or all Regions to that aggregation Region. This gives you a centralized view of all your findings across all of your accounts and all of your linked Regions.

## Examples

### Basic info

```sql
select
  arn,
  finding_aggregation_region,
  region_linking_mode
from
  aws_securityhub_finding_aggregator;
```

### List finding aggregators linked to all regions

```sql
select
  arn,
  finding_aggregation_region,
  region_linking_mode
from
  aws_securityhub_finding_aggregator
where
  region_linking_mode = 'ALL_REGIONS';
```

### List regions for finding aggregators that include specific regions

```sql
select
  arn,
  region_linking_mode,
  r as linked_region
from
  aws_securityhub_finding_aggregator,
  jsonb_array_elements_text(regions) as r
where
  region_linking_mode = 'SPECIFIED_REGIONS';
```

### List regions for finding aggregators that exclude specific regions

```sql
select
  arn,
  a.name as linked_region
from
  aws_redhood.aws_securityhub_finding_aggregator as f,
  aws_redhood.aws_region as a,
  jsonb_array_elements_text(f.regions) as r
where
  region_linking_mode = 'ALL_REGIONS_EXCEPT_SPECIFIED'
and
  a.name <> r;
```