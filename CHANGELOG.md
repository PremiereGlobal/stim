# Stim Changelog

## 0.0.5 (in-development)

### Features
* Added profile support for AWS logins. There are two new parameters for the `stim aws login` command:
  * `-p, --use-profiles` tells stim to use profiles which will save credentials as a new aws profile in `~/.aws/credentials`. The profile name will be in the format `<account>/<role>` to allow the profile to be reused in the future.  This option can also be set in the config file with `aws.use-profiles=true`.
  * `-d, --default-profile` tells stim to also set the aws default profile to the newly created creds. This option can also be set in the config file with `aws.default-profile=true`.

### Improvements
* Configurable auth methods.  A config entry of `auth.method` or command line parameter `--auth-method` can be used to specify the authentication method (ex: ldap, github, etc.).
* Prompt interface would trigger terminal bell when using arrow key to navigate the menu. This has been removed.

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
