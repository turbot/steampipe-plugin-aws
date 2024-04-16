select title, akas, region, partition, account_id
from aws_dms_endpoint
where endpoint_identifier = '{{ resourceName }}'
