select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings
where
  region = 'us-east-2';