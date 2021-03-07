select
  name,
  role_arn,
  status,
  title
from
  aws.aws_config_configuration_recorder;