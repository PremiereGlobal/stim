# Stim Changelog

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
