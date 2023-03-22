select display_name, title, akas, region, account_id
from aws_servicecatalog_portfolio
where display_name = '{{ resourceName }}' and region = '{{ output.aws_region.value }}';
