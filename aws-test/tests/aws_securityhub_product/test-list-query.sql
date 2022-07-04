select name, product_arn, company_name
from aws_securityhub_product
where name = 'Aqua Security' and region = '{{ output.aws_region.value }}';
