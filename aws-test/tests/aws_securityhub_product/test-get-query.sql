select name, product_arn, company_name, description, categories, activation_url
from aws_securityhub_product
where product_arn = 'arn:aws:securityhub:us-east-1:324264561773:product/guardicore/aws-infection-monkey';
