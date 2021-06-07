select
  id,
  resource_id,
  status
from
  aws.aws_ssm_managed_instance_compliance
where
  id = '{{ output.compliance_id.value }}';