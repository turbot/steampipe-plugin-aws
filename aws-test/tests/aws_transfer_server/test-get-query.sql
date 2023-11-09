select akas
from aws.aws_transfer_server
where tags ->> 'Name' = '{{resourceName}}'
