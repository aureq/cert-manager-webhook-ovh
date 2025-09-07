# OVH Webhook for Cert Manager

![OVH Webhook for Cert Manager](https://raw.githubusercontent.com/aureq/cert-manager-webhook-ovh/main/assets/images/cert-manager-webhook-ovh.svg "OVH Webhook for Cert Manager")

This is a webhook solver for [OVH](http://www.ovh.com) DNS. In short, if your domain has its DNS servers hosted with OVH, you can solve DNS challenges using Cert Manager and OVH Webhook for Cert Manager.

## Requirements

- A kubernetes cluster
- Helm v3 or higher ([https://helm.sh/docs/intro/install/](https://helm.sh/docs/intro/install/))
- Cert Manager 1.9 or higher installed ([https://cert-manager.io/docs/installation/](https://cert-manager.io/docs/installation/))
- An [OVH](https://www.ovh.com) account

## Preparation

Before you install anything, 2 mains tasks need to be completed.

### OVH API Keys

Obtaining API keys from your OVH account (in which your DNS zones are hosted) will allow this webhook to perform the necessary operations to help resolve Let's Encrypt [DNS01 challenges](https://letsencrypt.org/docs/challenge-types/#dns-01-challenge). In response to the DNS01 challenges, Let's Encrypt will issue a valid TLS certiticate for your applications to use.

1. Go to [api.ovh.com](https://api.ovh.com/console/) console and log-in using your credentials
2. Create a new [OVH API key](https://api.ovh.com/createToken/index.cgi?GET=/domain/zone/*&PUT=/domain/zone/*&POST=/domain/zone/*&DELETE=/domain/zone/*
) for this webhook (This needs to be done once per OVH account)

    - Application name: cert-manager-webhook-ovh (or anything you'd like)
    - Application description: API Keys for Cert Manager Webhook OVH (or anything you'd like)
    - Validity: Unlimited
    - Rights: (pre-populated)
      - `GET /domain/zone/*`
      - `PUT /domain/zone/*`
      - `POST /domain/zone/*`
      - `DELETE /domain/zone/*`
    - Restrict IPs: Leave blank or restrict as you need.

Take note of the `ApplicationKey`, `ApplicationSecret` and `ConsumerKey` and save them in a **secure** location.

### OVH OAuth2 Client

In order to use OVH IAM, you need to generate a OAuth2 Client, create a policy and link the policy to the client.

The following instructions are based on the [official documentation](https://help.ovhcloud.com/csm/en-manage-service-account?id=kb_article_view&sysparm_article=KB0059343).

1. Go to [api.ovh.com/console/](https://api.ovh.com/console/) console
2. Click `Authentication` link on the left handside to log-in using your OVH credentials
3. In the top left corner, select `v1` and then select `/me` API
4. On the left panel, search for `POST /me/api/oauth2/client` [↗️](https://api.ovh.com/console/?section=%2Fme&branch=v1#post-/me/api/oauth2/client)
5. Create a new service account with the following request body and click `Execute`

    ```json
    {
      "callbackUrls": [],
      "description": "Service account for OVH cert-manager webhook",
      "flow": "CLIENT_CREDENTIALS",
      "name": "cert-manager"
    }
    ```

6. Take note of both `ClientId` and `clientSecret` and save them in a **secure** location. Be carefull, you will not be able to retrieve the client secret later. You'll need to delete and create a new service account.
7. Navidate to `GET /me/api/oauth2/client/{clientId}` [↗️](https://api.ovh.com/console/?section=%2Fme&branch=v1#get-/me/api/oauth2/client/-clientId-)
8. Use the `ClientId` to retrieve the details of the service account. Take note of `identity`.

Now, you can create the policy to grant permissions on your domain to your service account.

1. In the top left corner, select `v2` and then select `/iam` API
2. Search for `POST /iam/policy` [↗️](https://api.ovh.com/console/?section=%2Fiam&branch=v2#post-/iam/policy)
3. Create a new IAM policy with the following request body. Adjust the `urn` to restrict the policy to one or more specifc domains and click `Execute`.

    ```json
    {
      "description": "Allow cert-manager of create and manage DNS records",
      "identities": [
        "<--- 🔥 INSERT SERVICE ACCOUNT IDENTITY FROM STEP 8 --->"
      ],
      "name": "cert-manager",
      "permissions": {
        "allow": [
          {
            "action": "dnsZone:apiovh:record/create"
          },
          {
            "action": "dnsZone:apiovh:record/delete"
          },
          {
            "action": "dnsZone:apiovh:refresh"
          }
        ]
      },
      "resources": [
        {
            "urn": "urn:v1:eu:resource:dnsZone:*", // ⚠️ all domains
            "urn": "urn:v1:eu:resource:dnsZone:example.net", // restrict the policy to "example.net"
            "urn": "urn:v1:eu:resource:dnsZone:example2.net", // restrict the policy to "example2.net"
        }
      ]
    }
    ```

4. You may take note of the policy Id in the response to more easily manage it in the future.
5. You can now use the OAuth2 Client ID and Client Secret on the webhook.

### Helm chart repository

1. Add a new Helm repository

    ```bash
    helm repo add cert-manager-webhook-ovh-charts https://aureq.github.io/cert-manager-webhook-ovh/
    ```

2. Refresh the repository information

    ```bash
    helm repo update
    ```

3. Search for all available charts in this repository

    ```bash
    helm search repo cert-manager-webhook-ovh-charts --versions
    ```

    Or list the latest development/unstable versions

    ```bash
    helm search repo cert-manager-webhook-ovh-charts --versions --devel
    ```

## Configuration

The configuration is done via `values.yaml` and for complete details, you should refer to the [repository](https://github.com/aureq/cert-manager-webhook-ovh/blob/main/charts/cert-manager-webhook-ovh/values.yaml).

- `configVersion` The config version is used each time a breaking changes introduced in `values.yaml` to prevent broken deployments. To retrieve the latest `values.yaml` which contains the updated configuration and the correct config version, run `helm show values cert-manager-webhook-ovh-charts/cert-manager-webhook-ovh`.
- `groupName` The GroupName here is used to identify your company or business unit that created this webhook.
- `certManager`
  - `namespace: cert-manager`: namespace where cert-manager is installed
  - `serviceAccountName: cert-manager` name of the cert-manager service account

- `issuers` A list of issuers as defined below

For each issuer:

- `name` Name of your issuer
- `create` When set to `true`, the issuer is created
- `kind` Can either be `ClusterIssuer` or `Issuer`. See documentation below.
- `namespace` If kind is `Issuer`, then indicate the namespace in which this issuer should be deployed into.
- `acmeServerUrl` which Acme Server URL to use.
- `email` An email address when registering an account with the Acme server.
- `ovhEndpointName` The endpoint name of the OVH API.
- `ovhAuthenticationMethod` Authentication method (possible values: application or oauth2)
- `ovhAuthentication` (cannot be use when `ovhAuthenticationRef` is used)
  - `ovhAuthentication.applicationKey` Your OVH application key.
  - `ovhAuthentication.applicationSecret` Your OVH application secret.
  - `ovhAuthentication.applicationConsumerKey` Your OVH consumer key.
  - `ovhAuthentication.oauth2ClientID` the OVH OAuth 2 client ID.
  - `ovhAuthentication.oauth2ClientSecret` the OVH OAuth 2 client Secret.
- `ovhAuthenticationRef` (cannot be use when `ovhAuthentication` is used)
  - `applicationKeyRef`
    - `name` Name of the Kubernetes secret
    - `key` The key name in the secret above that holds the actual value
  - `applicationSecretRef`
    - `name` Name of the Kubernetes secret
    - `key` The key name in the secret above that holds the actual value
  - `applicationConsumerKeyRef`
    - `name` Name of the Kubernetes secret
    - `key` The key name in the secret above that holds the actual value
  - `oauth2ClientIDRef`
    - `name` Name of the Kubernetes secret
    - `key` The key name in the secret above that holds the actual value
  - `oauth2ClientSecretRef`
    - `name` Name of the Kubernetes secret
    - `key` The key name in the secret above that holds the actual value

### Issuer vs ClusterIssuer

`Issuers`, and `ClusterIssuers`, are Kubernetes resources that represent certificate authorities (CAs) that are able to generate signed certificates by honoring certificate signing requests.

An `Issuer` is a namespaced resource, and it is not possible to issue certificates from an `Issuer` in a different namespace. This means you will need to create an `Issuer` in each namespace you wish to obtain `Certificates` in.

If you want to create a single `Issuer` that can be consumed in multiple namespaces, you should consider creating a `ClusterIssuer` resource. This is almost identical to the `Issuer` resource, however is non-namespaced so it can be used to issue `Certificates` across all namespaces.

If you are using cert-manager with an ingress controller, using a `ClusterIssuer` is recommended.

See cert-manager [documentation](https://cert-manager.io/docs/concepts/issuer/) for more details.

### Secret vs Secret References

It is usually safe to provide your OVH API keys as part of your `values.yaml` file or on the command line. However, it may be needed to separate the domain of responsibility between Ops (in charge of deploying the chart) and Security (in charge of obtaining and deploying the OVH API keys).

When providing your OVH API keys directly, this chart stores the provided OVH API keys in a secret (format shown  below) and then leverages secret references. The values of `applicationKey`, `applicationSecret` and `applicationConsumerKey` are base64 encoded.

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
  applicationConsumerKey: YUNvbnN1bWVyS2V5          # aConsumerKey
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

### Split Horizon DNS

Please note, using [Split Horizon DNS](https://en.wikipedia.org/wiki/Split-horizon_DNS) is unsupported. Use at your own risks.

If you plan to use cert-manager to generate certificates for domains served via split DNS you need to ensure that your DNS server returns correct SOA record for your domain. Otherwise the certificate generation will fail with following error message:

> Error presenting challenge: OVH API call failed: GET /domain/zone/com/status - HTTP Error 404: "This service does not exist"

In order to verify the record:

  ```bash
  dig example.com soa
  ```

Example correct response:

  ```ini
  ; <<>> DiG 9.18.1-1ubuntu1.3-Ubuntu <<>> example.com soa
  ;; global options: +cmd
  ;; Got answer:
  ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 57237
  ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

  ;; OPT PSEUDOSECTION:
  ; EDNS: version: 0, flags:; udp: 65494
  ;; QUESTION SECTION:
  ;example.com.     IN  SOA

  ;; ANSWER SECTION:
  example.com.    86400 IN  SOA dns.ovh.net. tech.ovh.net. 2023020402 86400 3600 3600000 600

  ;; Query time: 52 msec
  ;; SERVER: 127.0.0.53#53(127.0.0.53) (UDP)
  ;; WHEN: Sat Feb 04 15:13:03 CET 2023
  ;; MSG SIZE  rcvd: 88
  ```

Example incorrect response, notice the missing answer section:

  ```ini
  ; <<>> DiG 9.18.1-1ubuntu1.3-Ubuntu <<>> ovh.com soa
  ;; global options: +cmd
  ;; Got answer:
  ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 57237
  ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

  ;; OPT PSEUDOSECTION:
  ; EDNS: version: 0, flags:; udp: 65494
  ;; QUESTION SECTION:
  ;example.com.         IN  SOA

  ;; Query time: 52 msec
  ;; SERVER: 127.0.0.53#53(127.0.0.53) (UDP)
  ;; WHEN: Sat Feb 04 15:13:03 CET 2023
  ;; MSG SIZE  rcvd: 88
  ```

## Installation

To install the cert-manager-webhook-ovh chart:

```bash
helm upgrade --install --namespace cert-manager -f values.yaml cm-webhook-ovh cert-manager-webhook-ovh-charts/cert-manager-webhook-ovh
```

To uninstall the chart:

```bash
helm uninstall --namespace cert-manager cm-webhook-ovh
```
