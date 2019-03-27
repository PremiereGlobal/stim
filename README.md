# Stim
[![Build][Build-Status-Image]][Build-Status-Url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

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
├── pkg/
│   ├── pagerduty/
│   ├── utils/
│   ├── vault/
│   ├── ...
├── stim/
├── stimpacks/
│   ├── deploy/
│   ├── vault/
│   ├── ...
├── scripts/
├── vendor/
```

* `pkg/` The components in this directory should be developed as stand-alone packages that can be consumed not only by Stim but also externally.  They are generally wrappers around existing APIs (for example Vault) that simplify basic functionality.
* `stim/` This component is the core of the Stim application.  It is what every `stimpack` interfaces with to talk with the core Stim application.  Stim initializes components as-needed by the stimpacks.  For instance, if a stimpack needs access to Pagerduy, Stim will call Vault, get the API key for Pagerduty and instantiate a new instance of Pagerduty for the stimpack to use. Stim also allows stimpacks to attach cli commands and add configuration parameters.
* `stimpacks/` Stimpacks are pluggable extensions of the main Stim application.  They interface directly with the Stim api and can add commands and configuration to the cli.  They generally contain opionated functions for configuring developer workstations, building applications, testing, and deployments.

### Developing Stimpacks
See comments in `stimpacks/vault` for details
TODO: More docs here

### Developing Re-usable Packages
Guidelines:
* Don't log, just return errors and let the consumer deal with it


[Build-Status-Url]: https://travis-ci.org/ReadyTalk/stim
[Build-Status-Image]: https://travis-ci.org/ReadyTalk/stim.svg?branch=master
[reportcard-url]: https://goreportcard.com/badge/github.com/readytalk/stim
[reportcard-image]: https://goreportcard.com/badge/github.com/readytalk/stim
[godoc-url]: https://godoc.org/github.com/ReadyTalk/stim
[godoc-image]: https://godoc.org/github.com/ReadyTalk/stim?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
