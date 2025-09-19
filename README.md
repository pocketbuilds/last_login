# Last Login

An [xpb](https://github.com/pocketbuilds/xpb) plugin for [Pocketbase](https://pocketbase.io/) that allows for easily configured last_login fields.

Configure datetime fiels in auth collections to be automatically set on an auth request.

## Installation

1. [Install XPB](https://docs.pocketbuilds.com/installing-xpb).
2. [Use the builder](https://docs.pocketbuilds.com/using-the-builder):

```sh
xpb build --with github.com/pocketbuilds/last_login@latest
```

## Setup
1. Create a datetime field in an auth collection.
2. Add the field to the [plugin config](#plugin-config).
3. Restart the pocketbase app.

## Plugin Config

```toml
# pocketbuilds.toml

[last_login]
# Array of datetime fields in auth collections to
#   automatically set on auth request.
#   - format: "<collection_name>.<field_name>"
fields = [
  "my_auth_collection.last_login",
  "users.last_sign_in",
  # etc...
]
```
