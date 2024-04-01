select
  certificate_identifier,
  region
from
  aws.aws_dms_certificate
where
  certificate_identifier = '{{ resourceName }}-dummy';
