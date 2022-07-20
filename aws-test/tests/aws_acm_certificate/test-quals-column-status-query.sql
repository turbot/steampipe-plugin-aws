select title, status
from aws.aws_acm_certificate
where status = 'ISSUED' and certificate_arn = '{{ output.resource_aka.value }}';
