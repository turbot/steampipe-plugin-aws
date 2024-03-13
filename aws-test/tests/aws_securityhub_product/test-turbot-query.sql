select name, title, akas, region, account_id
from aws_securityhub_product
where name = 'Aqua Security' and region = '{{ output.aws_region.value }}';
