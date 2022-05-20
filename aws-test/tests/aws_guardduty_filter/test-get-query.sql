select
  name,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter
where
  name = '{{ output.filter_name.value }}';