select name, arn, id, scope, description, ip_address_version, addresses, partition, region, account_id
from aws.aws_wafv2_ip_set
where id = '{{ output.resource_id_regional.value }}';