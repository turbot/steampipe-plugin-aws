select
  name,
  arn,
  id,
  status,
  assessment_report_destination,
  assessment_report_destination_type,
  description,
  aws_account
from
  aws.aws_auditmanager_assessment
where
  id = '{{ output.assessment_id.value }}';