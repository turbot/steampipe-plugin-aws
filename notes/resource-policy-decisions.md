# Decisions and why

The code is here to determine the effective access.

## Issue 1

Priority: (LOW)

Always `allowed_principal_account_ids` have the user account. Remove the user account from `allowed_principal_account_ids` when account is **explicity** denied.

## Issue 2

If `access_level` is **not** `public` then the field `public_statement_ids` should be empty.

## Issue 3

Federated users should be `shared` today.

Need to explore if Federated users is `public`?
Is it public if there is an `audience`.

## Issue 4

The name of `access_level` is causing confusion

## Issue 5

Addition of the fields `private_access_level` and `shared_access_level`.

## Issue 6

In the scope of `"Effect": "Allow"`. Using the condition `aws:PrincipalAccounts`, if there is a wildcard in the value and we are looking at `StringLike` or `NotLike`, we take the view that the `access_level` is public.

## Issue 7

Suppose the case of Allow all and Deny one to many accounts, we take the view that it's `Public` because we allow to a range of accounts.

We do not remove the `*` from Principal Account IDs because of above.

## Issue 8

Write notes about `Public Statement IDs`

## Issue 9

Add `Shared Statement IDs` ...

## Issue 10

Add `Shared` and `Private` access level

## Issue 11

AllowedPrincipalAccountIds should change the wildcard from `*` to `0123456789??`, if this is the principal ID

## Issue 12

If there is no Principal we have 3 types of behaviour:

- No Resource, no principal in table
- Resource but no condition then the Principal is the User Account as this is an IDENTITY_POLICY
- Resource with condition then the Principal is the Conditions value otherwise it is the Principal

## Issue 13

If there is no Action that is valid, in other words, there is no permission, then there is no access, don't report
Action missing will no longer evaluate the statement

## Issue 14

Anything with `IfExists` is assume to not exists when evaulating the resource policy

## Issue 15

Our conditions should check the Pricipals better and evaluate if they should really show.

Conditions should be applied against Pricipals and not appended it Principals as is the behaviour at the moment.

## Issue 16

Denies with conditions <- Work on this tests

## Issue 17

This should not be valid: "aws:SourceArn": ["1234*"]
}

## Issue 18

If public service or public federation then we should give the direct value of the condition in Source Conditions.

## Issue 19

// Functionality is currently incorrect
It looks like SourceOrg and stuff is definately for services only and not Principals.

## Issue 20

PrincipalArn is a condition that works exclusively for the Principal

## Issue 21

It would be interesting to know that your policy has services but not restricting by using `Source*`
This means that the services are open to the world. Good use case.
https://docs.aws.amazon.com/sns/latest/dg/sns-access-policy-use-cases.html

## Issue 22

Case-sensitive matching of the ARN. Each of the six colon-delimited components of the ARN is checked separately and each can include multi-character match wildcards (\*) or single-character match wildcards (?). The ArnEquals and ArnLike condition operators behave identically.
