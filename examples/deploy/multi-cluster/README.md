# Multi-cluster deployment with Stim

In this example we're deploying and instance of Grafana to two environments: `stage` and `production`.  This differs slightly from the [basic example](../basic).  Here, our production environment has two instances, in two different Kubernetes clusters.  We're also pulling dynamic secrets from Vault for our database credentials.

To run this example, simply run `stim deploy` from this directory.
