# Table: aws_servicecatalog_product

AWS ServiceCatalog Product contains information about the product, such as its name, description, version, associated resources, configurations, and metadata. It serves as a central repository for managing and provisioning standardized IT services and resources across an organization.

## Examples

### Basic info

```sql
select
  name,
  id,
  product_id,
  type,
  support_url,
  support_email
from
  aws_servicecatalog_product;
```

### List products that has a default path

```sql
select
  name,
  id,
  product_id,
  type,
  distributor,
  owner,
  has_default_path
from
  aws_servicecatalog_product
where
  has_default_path;
```

### List products that are owned by AWS

```sql
select
  name,
  id,
  product_id,
  type,
  support_url,
  support_description
from
  aws_servicecatalog_product
where
  type = 'MARKETPLACE';
```

### Get budget details for each product

```sql
select
  name,
  product_id,
  b ->> 'BudgetName' as budget_name
from
  aws_servicecatalog_product,
  jsonb_array_elements(budgets) as b;
```

### Get launch path details for each product

```sql
select
  name,
  id,
  l ->> 'Id' as launch_path_id,
  l ->> 'Name' as launch_path_name
from
  aws_servicecatalog_product,
  jsonb_array_elements(launch_paths) as l;
```

### Get provisioning artifact details for each product

```sql
select
  name,
  id,
  p ->> 'Id' as provisioning_artifact_id,
  p ->> 'Name' as provisioning_artifact_name,
  p ->> 'CreatedTime' as provisioning_artifact_created_time,
  p ->> 'Description' as provisioning_artifact_description,
  p ->> 'Guidance' as provisioning_artifact_guidance
from
  aws_servicecatalog_product,
  jsonb_array_elements(provisioning_artifacts) as p;
```