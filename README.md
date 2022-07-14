# cloud-cli

Aptible `cloud-cli` application that allows interaction with the cloud api.

## Installation

There are currently two ways to install the cli.

### Github releases

Whenever we release a new version of the cli, we will publish a github release
and upload our official binaries for distribution.

[Github releases](https://github.com/aptible/cloud-cli/releases)

### Git

You're welcome to checkout our `main` branch and install the library using go.
Our goal is to ensure that `main` is **always** functional.

## Authentication

Currently the only way to login with aptible is via our [vintage
cli](https://github.com/aptible/aptible-cli).  Once you `login` the token saved
in your filesystem will be read for the new cli.

The more manual alternative is to:
- log into the web app
- open the developer tools
- watch a network request to our apis
- copy the token from the request
- paste it into ~/.aptible/tokens.json

```json
{"auth.aptible.com": "<TOKEN>"}
```

## Config

The goal of the configuration file is to allow the end-user to set some
defaults for all cli interaction.

Create a file `~/.aptible.yml`

```yml
org: "2253ae98-d65a-4180-aceb-8419b7416677"
env: "ENV-ID"
api-domain: "cloud-api.sandbox.aptible-cloud-staging.com"
debug: false
```

Nothing is required inside this file and cli arguments will take precedence.
For example if we saved the yaml file above and then ran:

```bash
aptible asset ls --org="1234"
```

`org=1234` would be used in this case because we overwrite whatever values are
stored in the file.
