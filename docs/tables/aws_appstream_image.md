---
title: "Steampipe Table: aws_appstream_image - Query AWS AppStream Images using SQL"
description: "Allows users to query AWS AppStream Images to gain insights into their properties, states, and associated metadata."
folder: "AppStream"
---

# Table: aws_appstream_image - Query AWS AppStream Images using SQL

AWS AppStream Images are part of Amazon AppStream 2.0, a fully managed, secure application streaming service that allows you to stream desktop applications from AWS to any device running a web browser. These images act as templates for the creation of streaming instances, containing all the necessary applications, drivers, and settings. Administrators can create, maintain, and use these images to provide a consistent user experience, regardless of the device being used.

## Table Usage Guide

The `aws_appstream_image` table in Steampipe provides you with information about images within AWS AppStream. This table allows you as a DevOps engineer to query image-specific details, including the image's name, ARN, state, platform, and associated metadata. You can utilize this table to gather insights on images, such as their visibility, status, and the applications they are associated with. The schema outlines the various attributes of the AppStream Image for you, including the image ARN, creation time, visibility status, and associated tags.

## Examples

### Basic info
Explore the details of your AWS AppStream images to understand their configuration and attributes. This can be beneficial in managing your resources and ensuring they are optimally configured.

```sql+postgres
select
  name,
  arn,
  base_image_arn,
  description,
  created_time,
  display_name,
  image_builder_name,
  tags
from
  aws_appstream_image;
```

```sql+sqlite
select
  name,
  arn,
  base_image_arn,
  description,
  created_time,
  display_name,
  image_builder_name,
  tags
from
  aws_appstream_image;
```

### List available images
Determine the areas in which AWS AppStream images are available for use. This is useful for understanding what resources are currently usable in your environment.

```sql+postgres
select
  name,
  arn,
  display_name,
  platform,
  state
from
  aws_appstream_image
where
  state = 'AVAILABLE';
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  platform,
  state
from
  aws_appstream_image
where
  state = 'AVAILABLE';
```

### List Windows based images
Identify instances where Windows based images are used within the AWS Appstream service. This is beneficial for auditing purposes, ensuring the correct platform is being utilized.

```sql+postgres
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  platform = 'WINDOWS';
```

```sql+sqlite
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  platform = 'WINDOWS';
```

### List images that support image builder
Identify the AWS AppStream images that are compatible with the image builder feature. This is useful to ensure your applications are using images that support this functionality for streamlined image creation and management.

```sql+postgres
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  image_builder_supported;
```

```sql+sqlite
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  image_builder_supported = 1;
```

### List private images
Explore which AppStream images are set to private to manage access and ensure security. This can help identify instances where images may need to be shared or restricted further.

```sql+postgres
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_name,
  visibility
from
  aws_appstream_image
where
  visibility = 'PRIVATE';
```

```sql+sqlite
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_name,
  visibility
from
  aws_appstream_image
where
  visibility = 'PRIVATE';
```

### Get application details of images
Explore the various attributes of applications within images, such as creation time, display name, and platform compatibility. This can be useful to understand the application's configuration and behavior for effective management and troubleshooting.

```sql+postgres
select
  name,
  arn,
  a ->> 'AppBlockArn' as app_block_arn,
  a ->> 'Arn' as app_arn,
  a ->> 'CreatedTime' as app_created_time,
  a ->> 'Description' as app_description,
  a ->> 'DisplayName' as app_display_name,
  a ->> 'Enabled' as app_enabled,
  a ->> 'IconS3Location' as app_icon_s3_location,
  a ->> 'IconURL' as app_icon_url,
  a ->> 'InstanceFamilies' as app_instance_families,
  a ->> 'LaunchParameters' as app_launch_parameters,
  a ->> 'LaunchPath' as app_launch_path,
  a ->> 'Name' as app_name,
  a ->> 'Platforms' as app_platforms,
  a ->> 'WorkingDirectory' as app_WorkingDirectory
from
  aws_appstream_image,
  jsonb_array_elements(applications) as a;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(a.value, '$.AppBlockArn') as app_block_arn,
  json_extract(a.value, '$.Arn') as app_arn,
  json_extract(a.value, '$.CreatedTime') as app_created_time,
  json_extract(a.value, '$.Description') as app_description,
  json_extract(a.value, '$.DisplayName') as app_display_name,
  json_extract(a.value, '$.Enabled') as app_enabled,
  json_extract(a.value, '$.IconS3Location') as app_icon_s3_location,
  json_extract(a.value, '$.IconURL') as app_icon_url,
  json_extract(a.value, '$.InstanceFamilies') as app_instance_families,
  json_extract(a.value, '$.LaunchParameters') as app_launch_parameters,
  json_extract(a.value, '$.LaunchPath') as app_launch_path,
  json_extract(a.value, '$.Name') as app_name
from
  aws_appstream_image,
  json_each(applications) as a;
```

### Get the permission model of the images
Determine the access permissions of specific images within your AWS AppStream service. This query is useful if you want to understand which images are accessible by your fleet and image builder, providing insights into your resource utilization and access control.

```sql+postgres
select
  name,
  arn,
  image_permissions ->> 'AllowFleet' as allow_fleet,
  image_permissions ->> 'AllowImageBuilder' as allow_image_builder
from
  aws_appstream_image;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(image_permissions, '$.AllowFleet') as allow_fleet,
  json_extract(image_permissions, '$.AllowImageBuilder') as allow_image_builder
from
  aws_appstream_image;
```

### Get error details of failed images
Discover the segments that contain failed images within your AWS AppStream environment. This query can be used to identify and analyze the issues causing image failures, helping to improve the efficiency and reliability of your AppStream services.

```sql+postgres
select
  name,
  arn,
  e ->> 'ErrorCode' as error_code,
  e ->> 'ErrorMessage' as error_message,
  e ->> 'ErrorTimestamp' as error_timestamp
from
  aws_appstream_image,
  jsonb_array_elements(image_errors) as e;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(e.value, '$.ErrorCode') as error_code,
  json_extract(e.value, '$.ErrorMessage') as error_message,
  json_extract(e.value, '$.ErrorTimestamp') as error_timestamp
from
  aws_appstream_image,
  json_each(image_errors) as e;
```