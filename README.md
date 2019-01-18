# Stim


<!-- # Running in Docker
Stim is available in Docker.  To use, simply run

```
docker run readytalk/stim <stim-command>
```

Stim natively supports configuration via environment variables. So, as an example, to log into Vault and map the token to your home directory, run

```
docker run \
  -it \
  -e VAULT_ADDR=https://my-vault-domain:8200 \
  -v $HOME/.vault-token:/root/vault-token \
  readytalk/stim vault login
``` -->
