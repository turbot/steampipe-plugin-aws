select name, product_arn, company_name
from aws_securityhub_product
where name = 'TestNotFound';
