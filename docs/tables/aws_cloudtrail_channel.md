# Table: aws_cloudtrail_channel

AWS services can create a service-linked channel to receive CloudTrail events on your behalf. The AWS service creating the service-linked channel configures advanced event selectors for the channel and specifies whether the channel applies to all regions, or a single region.

## Examples

### Basic info

```sql
select
  name,
  arn,
  source,
  apply_to_all_regions
from
  aws_cloudtrail_channel;
```

### List channels that are not applyed to all region

```sql
select
  name,
  arn,
  source,
  apply_to_all_regions,
  advanced_event_selectors
from
  aws_cloudtrail_channel
where
  not apply_to_all_regions;
```

### Get advance event selector details for each channel

```sql
select
  name,
  a ->> 'Name' as advanced_event_selector_name,
  a ->> 'FieldSelectors' as field_selectors
from
  aws_cloudtrail_channel
  jsonb_array_elements(advanced_event_selectors) as a;
```

