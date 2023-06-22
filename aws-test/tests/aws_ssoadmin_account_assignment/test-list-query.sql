select 
  principal_id,
  principal_type,
  instance_arn,
  permission_set_arn,
  target_account_id,
  partition,
  region,
  account_id
from aws.aws_ssoadmin_account_assignment
where permission_set_arn = '{{output.permission_set_arn.value}}' and target_account_id = '{{output.target_account_id.value}}';
