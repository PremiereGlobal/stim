# Stim Changelog

## 0.3.1
### Improvements
* changed thottle to 5 times and 2/4/8/16/32 sec
### Bugfix
* still fail if throttling happens more then 5 times


## 0.3.0
### Improvements
* Allow vault to retry getting AWS creditials on throttle message: `vault.retryOnThrottle:bool`
* Start work to commonize all configs to using more map style config based on stimpack name `stimpack.key`, IE `vault.address` `kube.config.namespace`

## 0.2.1
### Bugfix
* Fixed defaulted value for prompt list.

## 0.2.0
### Bugfix
* fixed build/publish for versions
* added filters for cluster and service account
* moved version to only output version nothing else.  Removed unneeded else statement
* added semver requirement, and checking.  Deploy added global.requiredVersion and global.minimumVersion to the stim.deploy.yaml, so you can enforce versions.  Fixed minor issue with logging Fatal, showing as Info
* Adding validation to check for no arg for completion cmd
* hard set to alpine:3.11 since alpine:3.12 does not support py2 anymore
* tab completion functionality for bash and zsh
* Allow all config options to also be envVars using subpaths as _ which then get translated to - in viper config
* Upped version of `kube-vault-deploy` to fix checksum issue downloading kubectl >= 1.18

## 0.1.7

### **Deprecations**
* For `stim deploy`, the `HELM_VERSION` environment variable for specifying Helm versions is now deprecated.  Please use the `.spec.tools.helm` configuration for specifying the helm version to use. See [deploy docs](https://github.com/PremiereGlobal/stim/blob/master/docs/DEPLOY.md) for more details.
* For `stim deploy`, auto-detection of Helm v2 versions is now deprecated.  Please use the `.spec.tools.helm` configuration for specifying the helm version to use. See [deploy docs](https://github.com/PremiereGlobal/stim/blob/master/docs/DEPLOY.md) for more details.

### Bugfix
* Fixed a bug where running `stim deploy` via the default Docker method would not respect the `spec.tools.helm` version.

## 0.1.6

### Bugfix
* Fixed bug in deploy config when setting variables in the global section did not propagate into subsections, like environments.

## 0.1.5

### Bugfix
* `stim deploy` - Fixed issue with cascading tool versions not working correctly. Certain instance tool versions would overwrite the tool version for other instances.
* Changed displayed username in AWS web console for the STS federated login. The username was 'stim-user' and now matches the LDAP sAMAccountName short name.

## 0.1.4

### Improvements
* Added new Docker image with the format `premiereglobal/stim:v0.1.4-deploy`.  This container provides additional features for deployments with stim
  * Contains additional deploy utilities such as: `bash`, `curl`, `zip`, `jq`, `less`, `python`, `git`, `yq`, and `aws`
  * Entrypoint is `bash` to make it easier to run with custom commands
* Updated Docker tagging
  * Master builds will now have the Docker tag `master` and `master-deploy`, with the `stim version` of `stim/master`
  * Tagged releases will continue having the tag scheme `v0.1.4` and `v0.1.4-deploy`, but will also update the latest tags `latest` and `latest-deploy`.  Doing `stim version` on any `latest` image will now reveal the actual version `stim/vX.X.X`, instead of `stim/master`.

## 0.1.3

### Improvements
* Added `--method` to `stim deploy` allowing optional `shell` deployment.  Will auto-select the best options if left empty.
  * `shell` deploy will be auto-selected if it detects that it is running in a container
* Added configuration option for cache location. See [docs/CONFIG.md](docs/CONFIG.md)

### Bugfix
* Fixed Pagerduty request to include all CLI parameters (for example source, component, etc.) as those were not actually being sent to Pagerduty previously

## 0.1.2

### Improvement
* Added configuration option: `vault-username-skip-prompt`. If set to `true` stim will skip prompting the user for their username if the config option `vault-username` is set.
  * Note: As part of this change, stim will no longer default to the system username. First-time users will be required to enter their username.
* Changing the deployment binary cache directory to be within the `~/.stim` dir

### Bugfix
* Fixed bug where deployments would not utilize the VAULT_TOKEN env variable correctly
* Fixed bug where having an empty stim cache directory (`~/.kube-vault-deploy/cache`) would cause a deployment failure

## 0.1.1

### Bugfix
* Fixed and issue that was causing environment-level env vars not to override global env vars in `stim deploy`

## 0.1.0
* Changing release versions to proper notation
* Fixed bug in `stim deploy` where it wasn't correctly resolving the full path for deployment volumes

## 0.0.7
* Added the `stim deploy` subcommand.  See [docs/DEPLOY.md](docs/DEPLOY.md) for details.
* Enabled `info` logging by default

## 0.0.6

### Features
* Added ttl support for AWS logins.  There are two new parameters for the `stim aws login` command:
  * `-t, --ttl` tells stim to set to AWS credentials ttl to the given value (defaults to `8h`)
  * `-b, --web-ttl` tells stim to set the AWS web console ttl to the given value (defaults to `1h`).  This value must be between `15m` and `36h`

## 0.0.5

### Features
* Added profile support for AWS logins. There are two new parameters for the `stim aws login` command:
  * `-p, --use-profiles` tells stim to use profiles which will save credentials as a new aws profile in `~/.aws/credentials`. The profile name will be in the format `<account>/<role>` to allow the profile to be reused in the future.  This option can also be set in the config file with `aws.use-profiles=true`.
  * `-d, --default-profile` tells stim to also set the aws default profile to the newly created creds. This option can also be set in the config file with `aws.default-profile=true`.

### Improvements
* Configurable auth methods.  A config entry of `auth.method` or command line parameter `--auth-method` can be used to specify the authentication method (ex: ldap, github, etc.).
* Prompt interface would trigger terminal bell when using arrow key to navigate the menu. This has been removed.
* Changed AWS login wait login to be linear vs exponential
* Vault token duration messages improved to show actual TTL vs requested when getting a new token
* Improved build fail detection

### Other
* Removed `bash` installer for now until it can be fixed

## 0.0.4
### BREAKING CHANGES
* Changing parameter for `stim aws login`. Changed -m to -a to be clear about account name.

### Features
* Added STS support.  This allows `stim aws login` to provision STS credentials with the IAM credentials it received from Vault.  This increases the utility as you can now provision web console access that is limited to the user's IAM credential role.

### Improvements
* Updated the logger for more robustness and readability
* AWS login now has -o option to print generated URL and not Launch

### Bug Fixes
* Fixed issue with `kube config` wherein the `--namespace` argument  was not being used correctly
