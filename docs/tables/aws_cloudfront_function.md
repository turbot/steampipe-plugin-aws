---
title: "Table: aws_cloudfront_function - Query AWS CloudFront Functions using SQL"
description: "Allows users to query AWS CloudFront Functions to retrieve detailed information about each function, including its ARN, stage, status, and more."
---

# Table: aws_cloudfront_function - Query AWS CloudFront Functions using SQL

The `aws_cloudfront_function` table in Steampipe provides information about functions within AWS CloudFront. This table allows DevOps engineers to query function-specific details, including the function's ARN, stage, status, and associated metadata. Users can utilize this table to gather insights on functions, such as their status, the events they are associated with, and more. The schema outlines the various attributes of the CloudFront function, including the function ARN, creation timestamp, last modified timestamp, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_function` table, you can use the `.inspect aws_cloudfront_function` command in Steampipe.

### Key columns:

- `name`: The name of the function. This column can be used to join this table with other tables for more detailed insights.
- `arn`: The Amazon Resource Name (ARN) of the function. It is a unique identifier for the function and can be used to join this table with other tables.
- `stage`: The function's stage, either `DEVELOPMENT` or `LIVE`. This information is useful for understanding the function's deployment status.