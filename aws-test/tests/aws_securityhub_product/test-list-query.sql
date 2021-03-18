select name, product_arn, company_name
from aws_securityhub_product
order by name asc limit 5;
