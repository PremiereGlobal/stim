# Stim [![Build Status](https://travis-ci.org/ReadyTalk/stim.svg?branch=master)](https://travis-ci.org/ReadyTalk/stim)

## Running in Docker
Stim is available in Docker.  To use, simply run

```
docker run readytalk/stim <stim-command>
```

Stim natively supports configuration via environment variables. So, for example, to log into Vault and map the token to your home directory, run

```
docker run \
  -it \
  -e VAULT_ADDR=https://my-vault-domain:8200 \
  -v $HOME/.vault-token:/root/vault-token \
  readytalk/stim vault login
```

## Developing with Stim

### Project Structure
The project is broken down into 4 major componenets: `api`, `command`, `packages`, and `stimpacks`. More explaination below.

```
├── api/
├── cmd/
├── pkg/
│   ├── pagerduty/
│   ├── utils/
│   ├── vault/
│   ├── ...
├── stimpacks/
│   ├── deploy/
│   ├── vault/
│   ├── ...
├── vendor/
```

* `pkg/` The components in this directory should be developed as stand-alone packages that can be consumed not only by Stim but also externally.  They are generally wrappers around existing APIs (for example Vault) that simplify basic functionality.
* `api/` This component is the core of the Stim application.  It is what every `stimpack` interfaces with to talk with the core Stim application.  The API initializes components as-needed by the stimpacks.  For instance, if a stimpack needs access to Pagerduy, the API will call Vault, get the API key for Pagerduty and instantiate a new instance of Pagerduty for the stimpack to use. The API also allows stimpacks to attach cli commands and add configuration parameters.
* `cmd/` This component configures the root-level cli command and creates the API component.  It's also responsible for setting up logging and calling all the top-level stimpacks.
* `stimpacks/` Stimpacks are pluggable extensions of the main Stim application.  They interface directly with the Stim api and can add commands and configuration to the cli.  They generally contain opionated functions for configuring developer workstations, building applications, testing, and deployments.

### Developing Stimpacks
TODO

### The Stim API
TODO

### Developing Re-usable Packages
Guidelines:
* Don't log, just return errors and let the consumer deal with it
