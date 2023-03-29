select display_name, id
from aws_servicecatalog_portfolio
where display_name = '{{ resourceName }}';
