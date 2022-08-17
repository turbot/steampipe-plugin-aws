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
