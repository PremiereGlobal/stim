# Basic deployment with Stim

In this example we're deploying an instance of Grafana to two environments: `stage` and `production`.  We have different admin passwords for each environment that we're getting from Vault.  Also, in this example, we're deploying to the same Kubernetes cluster, just a different namespace separating the two environments.

To run this example, simply run `stim deploy` from this directory
