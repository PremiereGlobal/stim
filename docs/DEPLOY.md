# Deployments with Stim

Stim can be used to deploy to Kubernetes using Vault to configure a deployment environment.

> WARNING: this functionality has not been tested for automated deploys.  Use at your own risk.

## Prerequisites
To use this functionality, the [Docker](https://docs.docker.com/install/) daemon must be installed and running on the machine.

## Usage

`stim deploy`

## Command Line Arguments

| Argument | Description |
| - | - |
| `-f, --deploy-file` | Location of the deployment config file to use.  Defaults to `./stim.deploy.yaml` |
| `-e, --environment` | Environment to deploy. If no value is provided, the user will be prompted. |
| `-i, --instance` | Instance to deploy to. The special value of "all" can be specified to deploy to all environments. If no value is provided, the user will be prompted. |

## Configuration
`stim deploy` is configured with a YAML file (`./stim.deploy.yaml` by default) that provides an inventory of the deployment environments as well as the configuration of those environments.

A basic config which is configured with one instance in one environment might look like:
```
environments:
  - name: stage
    instances:
      - name: us-west-2
        spec:
          kubernetes:
            cluster: blue.mydomin.com
            service-account: admin
          env:
            - name: NAMESPACE
              value: myapp
```

The `spec` section is a special section that can be built hierarchically.  There are three levels which can be set:

1. [Global](#global) level.  This will apply to all environments and instances.
2. [Environment](#environment) level.  This will apply to all instances within an environment.  This will override any conflicting global-level specs.
3. [Instance](#instance) level.  This will apply only to an individual instance.  This will override any conflicting global or environment level specs.

More examples can be found in the [examples directory](../examples).

See below for the details spec of the config file.

## Reserved Environment Variables

The following environment variables are created by `stim deploy` and can be used within the deployment or for debugging.  These are also considered reserved environment variable names and cannot be used in the deployment config.

| Env Var | Description |
| ----- | ----------- |
| `VAULT_ADDR` | Vault address |
| `VAULT_TOKEN` | Vault token |
| `SECRET_CONFIG` | Vault secret config specification |
| `DEPLOY_ENVIRONMENT` | Name of the environment which is being deployed to |
| `DEPLOY_INSTANCE` | Name of the `instance` that is being deployed to |
| `DEPLOY_CLUSTER` | Name of the Kubernetes cluster which is being deployed to |
| `CLUSTER_SERVER` | API endpoint for the Kubernetes cluster |
| `CLUSTER_CA` | Cluster CA for the Kubernetes cluster |
| `USER_TOKEN` | Token used to authenticate against the Kubernetes cluster |

## Config Spec

The deploy config file is a YAML file with the following root structure.

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `deployment` | Configuration for kicking off the deployment | [Deployment](#deployment) | `false` | |
| `global` | Global environment config | [Global](#global) | `false` | |
| `environments` | List of environment specifications | [[]Environment](#environment) | `true` | |

### Deployment

Configuration for kicking off the deployment

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `directory` | Deployment directory (relative to this config file). This directory will be mounted into the deployment container | `string` | `false` | `./` |
| `script` | Deployment script (relative to `directory`).  This is the script that will be executed after the environment is set up | `string` | `false` | `deploy.sh` |
| `container` | Configuration for the deploy container | [Container](#container) | `false` | |

### Container

Configuration for the deploy container

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `repo` | Docker repo | `string` | `false` | `premiereglobal/kube-vault-deploy` |
| `tag` | Docker tag | `string` | `false` | `0.3.1` |

### Global

Global environment config

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `spec` | Environment configuration specification | [Spec](#spec) | `false` | |

### Environment

Environment configuration specification

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `name` | Name of environment | `string` | `true` | |
| `spec` | Environment configuration specification | [Spec](#spec) | `false` | |
| `instances` | Inventory of instances within the environment | [[]Instance](#instance) | `true` | |

### Instance

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `name` | Name of the instance | `string` | `true` | |
| `spec` | Environment configuration specification | [Spec](#spec) | `true` | |

### Spec

The *Spec* represents a set of environment configurations that determine where the deployment happens as well as any environmental and/or secrets parameters.

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `kubernetes` | Kubernetes configuration | [Kubernetes](#kubernetes) | `false` | |
| `env` | Static environment variables | [[]EnvVar](#envvar) | `false` | |
| `secrets` | Secret configuration specification | [[]Secret](#secret) | `false` | |

### Kubernetes

The *Kubernetes* configuration specifies which cluster and auth to use when connecting to Kubernetes.

> Note: Although `kubernetes.cluster` can be set at the global or environment level, it is recommended that each `instance` explicitly call out which cluster it should be deployed to.  This will avoid any hierarchical misconfigurations.

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `cluster` | Name of the cluster to deploy to. This is required to be set somewhere along the hierarchy but not in each instance of this spec. | `string` | `false` | |
| `serviceAccount` | Name of the service account to authenticate with Kubernetes. This is required to be set somewhere along the hierarchy but not in each instance of this spec. | `string` | `false` | |

### EnvVar

The *EnvVar* type represents a shell environment variable consisting of a name and value. Reserved names shown in the [Reserved Environment Variables](#reserved-environment-variables) section are reserved and cannot be used here.

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `name` | Name of the environment variable | `string` | `true` | |
| `value` | Value of the environment variable | `string` | `true` | |

### SecretSpec

The *SecretSpec* type represents a definition of a Vault secret being pulled into an environment variable. See [vault-to-envs](https://github.com/PremiereGlobal/vault-to-envs) for more details.  Reserved names shown in the [Reserved Environment Variables](#reserved-environment-variables) section are reserved and cannot be used here.

| Field | Description | Type | Required | Default |
| ----- | ----------- | ------ | -------- | -------- |
| `secretPath` | The full path within Vault where the secret is stored. | `string` | `false` | |
| `set` | Key-value mappings of environment variable names to secret field names | `map[string]string` | `true` | |
| `version` | The version to pull for Vault kv2 secrets.  Can be negative to "go back" x number of version.  For example, `-1` will pull the last previous version.  | `unsigned int` | `true` | |
| `ttl` | The time-to-live, in seconds, for dynamic secrets. | `int`| `false` | |
