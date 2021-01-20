select account_id, akas, region, title, tags, partition
from aws.aws_vpc_endpoint
where vpc_endpoint_id = '{{output.resource_id.value}}:asdf'