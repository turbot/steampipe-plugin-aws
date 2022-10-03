select
  member_account_id,
  email,
  detector_id
from
  aws_guardduty_member
where
  detector_id = '{{ output.detector_id.value }}';