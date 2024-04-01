select title, akas, tags
from aws.aws_dms_certificate
where certificate_identifier = '{{ resourceName }}';
