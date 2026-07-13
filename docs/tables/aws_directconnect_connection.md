---
title: "Steampipe Table: aws_directconnect_connection - Query AWS Direct Connect Connections using SQL"
description: "Allows users to query AWS Direct Connect Connection resources for detailed information."
folder: "DX"
---

# Table: aws_directconnect_connection - Query AWS Direct Connect Connections using SQL

AWS Direct Connect (DX) is a service that allows creating dedicated fiber connections with other networks, such as on-premises data centers, corporate networks, or cloud service providers.

A DX Connection represents a single physical fiber connection between an AWS Direct Connect router and the customer router.

## Table Usage Guide

The `aws_directconnect_connection` table in Steampipe provides you with information about DX Connections. This table allows you to query connection details, including its location, port speed, encryption configuration, and state.

## Examples

### Basic Connection info
Gain insights into the locations, port speed, and status:

```sql+postgres
select
  connection_id,
  connection_name,
  location,
  bandwidth,
  connection_state
from
  aws_directconnect_connection;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  location,
  bandwidth,
  connection_state
from
  aws_directconnect_connection;
```


### List all connections that are down
Determine which connections are down:


```sql+postgres
select
  connection_id,
  connection_name,
  location,
  bandwidth,
  connection_state
from
  aws_directconnect_connection
where
  connection_state = 'down';
```

```sql+sqlite
select
  connection_id,
  connection_name,
  location,
  bandwidth,
  connection_state
from
  aws_directconnect_connection
where
  connection_state = 'down';
```




### Understand the number of connections per location
Determine how many connections you have in which DX location:

```sql+postgres
select
  location,
  count(connection_id) as connections
from
  aws_directconnect_connection
group by
  location
order by connections desc;
```

```sql+sqlite
select
  location,
  count(connection_id) as connections
from
  aws_directconnect_connection
group by
  location
order by connections desc;
```
