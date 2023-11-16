---
title: "Table: aws_appstream_image - Query AWS AppStream Images using SQL"
description: "Allows users to query AWS AppStream Images to gain insights into their properties, states, and associated metadata."
---

# Table: aws_appstream_image - Query AWS AppStream Images using SQL

The `aws_appstream_image` table in Steampipe provides information about images within AWS AppStream. This table allows DevOps engineers to query image-specific details, including the image's name, ARN, state, platform, and associated metadata. Users can utilize this table to gather insights on images, such as their visibility, status, and the applications they are associated with. The schema outlines the various attributes of the AppStream Image, including the image ARN, creation time, visibility status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_appstream_image` table, you can use the `.inspect aws_appstream_image` command in Steampipe.

**Key columns**:

- `name`: The name of the image. This can be used to join with other tables that reference the image by name.
- `arn`: The Amazon Resource Number (ARN) of the image. This is a unique identifier that can be used to join with other tables that reference the image by ARN.
- `state`: The current state of the image. This can be useful when joining with other tables to understand the status of the image.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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

### List private images

```sql
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

```sql
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

### Get the permission model of the images

```sql
select
  name,
  arn,
  image_permissions ->> 'AllowFleet' as allow_fleet,
  image_permissions ->> 'AllowImageBuilder' as allow_image_builder
from
  aws_appstream_image;
```

### Get error details of failed images

```sql
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
