# Configuration
Stim can be configured in several ways.  Many options can be configured via environment variables, config files or command-line parameters.

## Order of Presedence
For all stim options:

`CLI Options` *overrides* `Environment Variables` *overrides* `Config File`

## Global
These options configure where stim looks for and stores its core configuration data.
|Env Var| Global CLI Option | Description | Default |
|---|---|---|---|
| `STIM_PATH` | `--path` | Path to the stim directory.  This is the default location for configuration files. | `${HOME}/.stim`|
| `STIM_CACHE_PATH` | `--cache-path` | Path for caching data. See [CACHE.md](CACHE.md) for more details. | `${STIM_PATH}/cache` |
| `STIM_CONFIG_FILE` | `--config` | Path for the global stim configuration file | `${STIM_PATH}/config.yaml`|

### Stim Config File
Additional configuration can be set in the `STIM_CONFIG_FILE`.

| Option | Description | Type | Default |
|---|---|---|---|
| `path` |  | `string` | `token` |
| `cache-path` |  | `string` | `token` |
| `auth.method` | Method to use for authentication.  Currently this would be the Vault auth-backend to use. | `string` | `token` |
| `aws.default-profile` | When fetching AWS credential, set to default AWS profile (in `~/.aws/credentials`). | `bool` | `false` |
| `aws.ttl` | Default ttl to set when fetching AWS credentials. (ex. `24h`) | `duration` | `Vault Default Setting` |
| `aws.use-profiles` | When fetching AWS credential, store the credentials as AWS profile (in `~/.aws/credentials`). | `bool` | `false` |
| `aws.web-ttl` | TTL for AWS web logins. | `duration` | `AWS default` |
| `logging.file.disable` | Option to disable file logging | `boolean` | `false` |
| `logging.file.level` | File logging verbosity | `string` | `info` |
| `logging.file.path` | File logging path | `string` | `info` |
| `pagerduty.vault-apikey-key` | Vault key for the Pagerduty API key | `string` | ` ` |
| `pagerduty.vault-apikey-path` | Vault path for the Pagerduty API key | `string` | ` ` |
| `vault-address` | Address to be used for connecting with Vault | `string` | ` ` |
| `vault-initial-token-duration` | Default token duration to use when authenticating with Vault | `duration` | `Vault Default Setting` |
| `vault-username` | Default username to use when logging into Vault | `string` | `Vault Default Setting` |
| `vault-username-skip-prompt` | Skip the username prompt if `vault-username` is set | `bool` | `false` |
| `verbose` | Use verbose logging | `bool` | `false` |
