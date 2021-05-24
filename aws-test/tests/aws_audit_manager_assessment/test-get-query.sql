select
  name,
  arn,
  id,
  status,
  assessment_report_destination,
  assessment_report_destination_type,
  description,
  aws_account,
  scope
from
  aws.aws_audit_manager_assessment
where
  id = '{{ output.assessment_id.value }}';