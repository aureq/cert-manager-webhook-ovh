![OVH Webhook for Cert Manager](https://raw.githubusercontent.com/aureq/cert-manager-webhook-ovh/main/assets/images/cert-manager-webhook-ovh.svg "OVH Webhook for Cert Manager")

This is a webhook solver for [OVH](http://www.ovh.com) DNS. In short, if your domain has its DNS servers hosted with OVH, you can solve DNS challenges using Cert Manager and OVH Webhook for Cert Manager.

## Requirements

- A kubernetes cluster
- Helm v3 or higher (https://helm.sh/docs/intro/install/)
- Cert Manager 1.9 or higher installed (https://cert-manager.io/docs/installation/)
- An [OVH](https://www.ovh.com) account

## Preparation

Before you install anything, 2 mains tasks need to be completed.

### OVH API Keys

Obtaining API keys from your OVH account (in which your DNS zones are hosted) will allow this webhook to perform the necessary operations to help resolve Let's Encrypt [DNS01 challenges](https://letsencrypt.org/docs/challenge-types/#dns-01-challenge). In response to the DNS01 challenges, Let's Encrypt will issue a valid TLS certiticate for your applications to use.

1. Go to [api.ovh.com](https://api.ovh.com/console/) console and log-in using your credentials
2. Create a new [OVH API key](https://api.ovh.com/createToken/index.cgi?GET=/domain/zone/*&PUT=/domain/zone/*&POST=/domain/zone/*&DELETE=/domain/zone/*
) for this webhook (This needs to be done once per OVH account)
    * Application name: cert-manager-webhook-ovh (or anything you'd like)
    * Application description: API Keys for Cert Manager Webhook OVH (or anything you'd like)
    * Validity: Unlimited
    * Rights: (pre-populated)
      * `GET /domaine/zone/*`
      * `PUT /domaine/zone/*`
      * `POST /domaine/zone/*`
      * `DELETE /domaine/zone/*`
    * Restrict IPs: Leave blank or restrict as you need.

Securely take note of the `ApplicationKey`, `ApplicationSecret` and `ConsumerKey`.

### Helm chart repository

1. Add a new Helm repository

    ```
    helm repo add cert-manager-webhook-ovh-charts https://aureq.github.io/cert-manager-webhook-ovh/
    ```

2. Refresh the repository information

    ```
    helm repo update
    ```

3. Search for all available charts in this repository

    ```
    helm search repo cert-manager-webhook-ovh-charts --versions
    ```

    Or list the latest development/unstable versions

    ```
    helm search repo cert-manager-webhook-ovh-charts --versions --devel
    ```

## Configuration

The configuration is done via `values.yam` and for complete details, you should refer to the [repository](https://github.com/aureq/cert-manager-webhook-ovh/blob/main/charts/cert-manager-webhook-ovh/values.yaml).

* `groupName` The GroupName here is used to identify your company or business unit that created this webhook.
* `certManager`
  * `namespace: cert-manager`: namespace where cert-manager is installed
  * `serviceAccountName: cert-manager` name of the cert-manager service account

* `issuers` A list of issuers as defined below

For each issuer
* `name` Name of your issuer
* `create` When set to `true`, the issuer is created
* `kind` Can either be `ClusterIssuer` or `Issuer`. See documentation below.
* `namespace` If kind is `Issuer`, then indicate the namespace in which this issuer should be deployed into.
* `acmeServerUrl` which Acme Server URL to use.
* `email` An email address when registering an account with the Acme server.
* `ovhEndpointName` The endpoint name of the OVH API.
* `ovhAuthentication` (cannot be use when `ovhAuthenticationRef` is used)
  * `ovhAuthentication.applicationKey` Your OVH application key.
  * `ovhAuthentication.applicationSecret` Your OVH application secret.
  * `ovhAuthentication.consumerKey` Your OVH consumer key.
* `ovhAuthenticationRef` (cannot be use when `ovhAuthentication` is used)
  * `applicationKeyRef`
    * `name` Name of the Kubernetes secret
    * `key` The key name in the secret above that holds the actual value
  * `applicationSecretRef`
    * `name` Name of the Kubernetes secret
    * `key` The key name in the secret above that holds the actual value
  * `consumerKeyRef`
    * `name` Name of the Kubernetes secret
    * `key` The key name in the secret above that holds the actual value

### Issuer vs ClusterIssuer

`Issuers`, and `ClusterIssuers`, are Kubernetes resources that represent certificate authorities (CAs) that are able to generate signed certificates by honoring certificate signing requests.

An `Issuer` is a namespaced resource, and it is not possible to issue certificates from an `Issuer` in a different namespace. This means you will need to create an `Issuer` in each namespace you wish to obtain `Certificates` in.

If you want to create a single `Issuer` that can be consumed in multiple namespaces, you should consider creating a `ClusterIssuer` resource. This is almost identical to the `Issuer` resource, however is non-namespaced so it can be used to issue `Certificates` across all namespaces.

If you are using cert-manager with an ingress controller, using a `ClusterIssuer` is recommended.

See cert-manager [documentation](https://cert-manager.io/docs/concepts/issuer/) for more details.

### Secret vs Secret References

It is usually safe to provide your OVH API keys as part of your `values.yaml` file or on the command line. However, it may be needed to separate the domain of responsibility between Ops (in charge of deploying the chart) and Security (in charge of obtaining and deploying the OVH API keys).

When providing your OVH API keys directly, this chart stores the provided OVH API keys in a secret (format shown  below) and then leverages secret references. The values of `applicationKey`, `applicationSecret` and `consumerKey` are base64 encoded.

If you decide to use secret references, simply indicate which secret name to use, and which key name in the secret holds the correct value.

If you are using a `ClusterIssuer`, then the secret should be stored in the same namespace as this webhook (usually `cert-manager`). And when using an `Issuer`, the secret should be stored in the same namespace as the issuer itself.

Example of a secret generated by this Helm Chart.
```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: ovh-credentials
  namespace: cert-manager
data:
  applicationKey: YW5BcHBsaWNhdGlvbktleQ==          # anApplicationKey
  applicationSecret: YW5BcHBsaWNhdGlvblNlY3JldA==   # anApplicationSecret
  consumerKey: YUNvbnN1bWVyS2V5                     # aConsumerKey
```

### Proxy support

If your Kubernetes cluster requires a proxy to access the internet, you have the ability to set environment variables to help the webhook. This is done in `values.yaml`.

```yaml
# Use this field to add environment variables relevant to this webhook.
# These fields will be passed on to the container when Chart is deployed.
environment:
  # Use these variables to configure the HTTP_PROXY environment variables
  HTTP_PROXY: "http://proxy:8080"
  HTTPS_PROXY: "http://proxy:8080"
  NO_PROXY: "10.0.0.0/8,127.0.0.0/8,172.16.0.0/12,192.168.0.0/16,*.svc,*.cluster.local,*.svc.cluster.local,169.254.169.254,127.0.0.1,localhost,localhost.localdomain"
```

### Split DNS support

If you plan to use cert-manager to generate certificates for domains served via split DNS you need to ensure that your DNS server returns correct SOA record for your domain. Otherwise the certificate generation will fail with following error message:
> Error presenting challenge: OVH API call failed: GET /domain/zone/com/status - HTTP Error 404: "This service does not exist"

In order to verify the record:

  ```
  dig yourdomain.com soa
  ```

Example correct response:

  ```
  ; <<>> DiG 9.18.1-1ubuntu1.3-Ubuntu <<>> ovh.com soa
  ;; global options: +cmd
  ;; Got answer:
  ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 57237
  ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

  ;; OPT PSEUDOSECTION:
  ; EDNS: version: 0, flags:; udp: 65494
  ;; QUESTION SECTION:
  ;ovh.com.			IN	SOA

  ;; ANSWER SECTION:
  ovh.com.		86400	IN	SOA	dns.ovh.net. tech.ovh.net. 2023020402 86400 3600 3600000 600

  ;; Query time: 52 msec
  ;; SERVER: 127.0.0.53#53(127.0.0.53) (UDP)
  ;; WHEN: Sat Feb 04 15:13:03 CET 2023
  ;; MSG SIZE  rcvd: 88
  ```

Example incorrect response, notice the missing answer section:

  ```
  ; <<>> DiG 9.18.1-1ubuntu1.3-Ubuntu <<>> ovh.com soa
  ;; global options: +cmd
  ;; Got answer:
  ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 57237
  ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

  ;; OPT PSEUDOSECTION:
  ; EDNS: version: 0, flags:; udp: 65494
  ;; QUESTION SECTION:
  ;ovh.com.			IN	SOA

  ;; Query time: 52 msec
  ;; SERVER: 127.0.0.53#53(127.0.0.53) (UDP)
  ;; WHEN: Sat Feb 04 15:13:03 CET 2023
  ;; MSG SIZE  rcvd: 88
  ```

## Installation

To install the cert-manager-webhook-ovh chart:

    helm upgrade --install --namespace cert-manager -f values.yaml cm-webhook-ovh cert-manager-webhook-ovh-charts/cert-manager-webhook-ovh

To uninstall the chart:

    helm uninstall --namespace cert-manager cm-webhook-ovh