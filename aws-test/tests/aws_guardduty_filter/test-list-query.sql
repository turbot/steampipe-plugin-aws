select
  name,
  detector_id
from
  aws_guardduty_filter
where
  action = 'ARCHIVE';