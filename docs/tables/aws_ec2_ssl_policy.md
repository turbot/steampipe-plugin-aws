---
title: "Table: aws_ec2_ssl_policy - Query AWS EC2 SSL Policies using SQL"
description: "Allows users to query AWS EC2 SSL Policies to retrieve detailed information about SSL policies used in AWS EC2 Load Balancers."
---

# Table: aws_ec2_ssl_policy - Query AWS EC2 SSL Policies using SQL

The `aws_ec2_ssl_policy` table in Steampipe provides information about SSL policies used in AWS Elastic Compute Cloud (EC2) Load Balancers. This table allows developers and cloud architects to query SSL policy-specific details, including the policy name, the SSL protocols, and the cipher suite configurations. Users can utilize this table to gather insights on the SSL policies, such as enabled SSL protocols, preferred cipher suites, and more. The schema outlines the various attributes of the SSL policy, including the policy name, the SSL protocols, the SSL ciphers, and the server order preference.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_ssl_policy` table, you can use the `.inspect aws_ec2_ssl_policy` command in Steampipe.

**Key columns**:

- `name`: The name of the SSL policy. This is the primary key of the table and can be used to join with other tables that reference SSL policies.
- `ssl_protocols`: The list of SSL protocols enabled for the policy. This can be used to determine the SSL protocol configurations for Load Balancers.
- `ssl_ciphers`: The set of SSL ciphers enabled for the policy. This is useful for analyzing the cipher suite configurations for Load Balancers.

## Examples

### Basic info

```sql
select
  name,
  ssl_protocols
from
  aws_ec2_ssl_policy;
```


### List load balancer listeners that use an SSL policy with weak ciphers

```sql
select
  arn,
  ssl_policy
from
  aws_ec2_load_balancer_listener listener
join 
  aws_ec2_ssl_policy ssl_policy
on
  listener.ssl_policy = ssl_policy.Name
where
  ssl_policy.ciphers @> '[{"Name":"DES-CBC3-SHA"}]';
```
