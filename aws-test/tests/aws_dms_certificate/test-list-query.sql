select certificate_identifier, arn
from aws.aws_dms_certificate
where arn = '{{ output.resource_aka.value }}'
