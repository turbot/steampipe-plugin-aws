select
	name,
	product_arn,
	company_name,
	description,
	categories,
	activation_url
from
	aws_securityhub_product
where
	product_arn = 'arn:aws:securityhub:{{ output.aws_region.value }}::product/aws/guardduty';