<!-- markdownlint-configure-file
{
    "MD004": false,
    "MD010": false,
    "MD033": {
        "allowed_elements": [
            "table",
            "thead",
            "th",
            "tbody",
            "td",
            "tr",
            "pre"
        ]
    },
    "MD060": false
}
-->

# cert-manager-webhook-ovh

![Version: 0.9.2](https://img.shields.io/badge/Version-0.9.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.9.2](https://img.shields.io/badge/AppVersion-0.9.2-informational?style=flat-square)

![OVH Webhook for Cert Manager](https://raw.githubusercontent.com/aureq/cert-manager-webhook-ovh/main/assets/images/cert-manager-webhook-ovh.svg "OVH Webhook for Cert Manager")

This is a webhook solver for [OVH](http://www.ovh.com) DNS. In short, if your domain has its DNS servers hosted with OVH, you can solve DNS challenges using Cert Manager and OVH Webhook for Cert Manager.

**Homepage:** <https://github.com/aureq/cert-manager-webhook-ovh>

## ⚠️ Development Status

**This chart is in early development (pre-1.0).** For this chart we use 0.major.minor semantic versioning, which means:
- **No stability guarantees** - Anything may change in any release including a breaking change
- **API may change without notice** - Chart structure and values can change significantly
- **It may not be suitable for production use** - Use at your own risk, and additional caution is recommended
- **Always pin to specific versions** - Avoid version ranges to prevent unexpected breaking changes during updates

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

With the OVHcloud API:

1. Go to [api.ovh.com/console/](https://api.ovh.com/console/) console
2. Click `Authentication` link on the left handside to log-in using your OVH credentials
3. In the top left corner, select `v1` and then select `/me` API
4. On the left panel, search for `POST /me/api/oauth2/client` [↗️](https://api.ovh.com/console/?section=%2Fme&branch=v1#post-/me/api/oauth2/client)
5. Create a new service account with the following request body and click on the `Execute` button.

    ```json
    {
      "callbackUrls": [],
      "description": "Service account for OVH cert-manager webhook",
      "flow": "CLIENT_CREDENTIALS",
      "name": "cert-manager"
    }
    ```

6. Take note of both `ClientId` and `ClientSecret` and save them in a **secure** location. Be carefull, you will not be able to retrieve the client secret later. You'll need to delete and create a new service account.
7. Navigate to `GET /me/api/oauth2/client/{clientId}` [↗️](https://api.ovh.com/console/?section=%2Fme&branch=v1#get-/me/api/oauth2/client/-clientId-)
8. Use the `ClientId` to retrieve the details of the service account. Take note of `identity`.

OR, with the `ovhcloud` CLI:

```bash
ovhcloud account api oauth2 client create --description Service account for OVH cert-manager webhook --flow CLIENT_CREDENTIALS --name cert-manager
```

You should obtain a response like this:
```plain
✅ OAuth2 client created successfully (client ID: xxxxxxxx, client secret: xxxxxxxx)
```

Take note of both `ClientId` and `ClientSecret` and save them in a **secure** location. Be carefull, you will not be able to retrieve the client secret later. You'll need to delete and create a new service account.

Use the `ClientId` to retrieve the details of the service account:

```bash
ovhcloud account api oauth2 client get EU.590d00c34ddd55c3
```

You should obtain a response like this:
```json
{
  "callbackUrls": null,
  "clientId": "xxxxxxxxxx",
  "createdAt": "2025-11-26T14:30:04.492Z",
  "description": "Service",
  "flow": "CLIENT_CREDENTIALS",
  "identity": "urn:v1:eu:identity:credential:xxxxxxxxxx/oauth2-xxxxxxxxxx",
  "name": "cert-manager"
}
```

Take note of the value of the `identity` field.

Now, you can create the policy to grant permissions on your domain to your service account.

With the OVHcloud API:

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
            "action": "dnsZone:apiovh:record/get"
          },
          {
            "action": "dnsZone:apiovh:record/delete"
          },
          {
            "action": "dnsZone:apiovh:status/get"
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

The configuration is done via [`values.yaml`](https://github.com/aureq/cert-manager-webhook-ovh/blob/main/charts/cert-manager-webhook-ovh/values.yaml).

The documentation below is automatically generated from the `values.yaml` file.

## Values

<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>groupName</td>
			<td>string</td>
			<td><pre lang="json">
"acme.mycompany.example"
</pre>
</td>
			<td>The `groupName` here is used to identify your company or business unit that created this webhook. For example, this may be `acme.mycompany.example`. This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain here is recommended.</td>
		</tr>
		<tr>
			<td>certManager</td>
			<td>object</td>
			<td><pre lang="json">
{
  "namespace": "cert-manager",
  "serviceAccountName": "cert-manager"
}
</pre>
</td>
			<td>The values below should match your cert-manager deployment. If not, you will get permissions errors.</td>
		</tr>
		<tr>
			<td>certManager.namespace</td>
			<td>string</td>
			<td><pre lang="json">
"cert-manager"
</pre>
</td>
			<td>namespace in which your cert-manager is deployed (default: cert-manager)</td>
		</tr>
		<tr>
			<td>certManager.serviceAccountName</td>
			<td>string</td>
			<td><pre lang="json">
"cert-manager"
</pre>
</td>
			<td>cert-manager serverAccount name (default: cert-manager)</td>
		</tr>
		<tr>
			<td>issuers</td>
			<td>list</td>
			<td><pre lang="json">
[
  {
    "acmeServerUrl": "https://acme-v02.api.letsencrypt.org/directory",
    "cnameStrategy": "None",
    "create": false,
    "disableAccountKeyGeneration": false,
    "email": "acme@example.net",
    "externalAccountBinding": {
      "enabled": false,
      "keyID": "",
      "keySecretRef": {
        "key": "",
        "name": ""
      }
    },
    "kind": "ClusterIssuer",
    "name": "le-prod",
    "namespace": "default",
    "ovhAuthentication": {
      "applicationConsumerKey": "",
      "applicationKey": "",
      "applicationSecret": "",
      "oauth2ClientID": "",
      "oauth2ClientSecret": ""
    },
    "ovhAuthenticationMethod": "application",
    "ovhAuthenticationRef": {
      "applicationConsumerKeyRef": {
        "key": "applicationConsumerKey",
        "name": ""
      },
      "applicationKeyRef": {
        "key": "applicationKey",
        "name": ""
      },
      "applicationSecretRef": {
        "key": "applicationSecret",
        "name": ""
      },
      "oauth2ClientIDRef": {
        "key": "oauth2ClientID",
        "name": ""
      },
      "oauth2ClientSecretRef": {
        "key": "oauth2ClientSecret",
        "name": ""
      }
    },
    "ovhEndpointName": "ovh-eu",
    "profile": "classic"
  }
]
</pre>
</td>
			<td>One or more issuers to create</td>
		</tr>
		<tr>
			<td>issuers[0]</td>
			<td>object</td>
			<td><pre lang="json">
{
  "acmeServerUrl": "https://acme-v02.api.letsencrypt.org/directory",
  "cnameStrategy": "None",
  "create": false,
  "disableAccountKeyGeneration": false,
  "email": "acme@example.net",
  "externalAccountBinding": {
    "enabled": false,
    "keyID": "",
    "keySecretRef": {
      "key": "",
      "name": ""
    }
  },
  "kind": "ClusterIssuer",
  "name": "le-prod",
  "namespace": "default",
  "ovhAuthentication": {
    "applicationConsumerKey": "",
    "applicationKey": "",
    "applicationSecret": "",
    "oauth2ClientID": "",
    "oauth2ClientSecret": ""
  },
  "ovhAuthenticationMethod": "application",
  "ovhAuthenticationRef": {
    "applicationConsumerKeyRef": {
      "key": "applicationConsumerKey",
      "name": ""
    },
    "applicationKeyRef": {
      "key": "applicationKey",
      "name": ""
    },
    "applicationSecretRef": {
      "key": "applicationSecret",
      "name": ""
    },
    "oauth2ClientIDRef": {
      "key": "oauth2ClientID",
      "name": ""
    },
    "oauth2ClientSecretRef": {
      "key": "oauth2ClientSecret",
      "name": ""
    }
  },
  "ovhEndpointName": "ovh-eu",
  "profile": "classic"
}
</pre>
</td>
			<td>Name of this issuer</td>
		</tr>
		<tr>
			<td>issuers[0].create</td>
			<td>bool</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>When true this issuer is created.</td>
		</tr>
		<tr>
			<td>issuers[0].kind</td>
			<td>string</td>
			<td><pre lang="json">
"ClusterIssuer"
</pre>
</td>
			<td>The type of issuer to create. It can either be `ClusterIssuer` or `Issuer`. If `Issuer` is specified, then the `namespace` is also required. See https://cert-manager.io/docs/concepts/issuer/ for more information.</td>
		</tr>
		<tr>
			<td>issuers[0].namespace</td>
			<td>string</td>
			<td><pre lang="json">
"default"
</pre>
</td>
			<td>If `kind` is `Issuer`, then indicate the namespace in which this issuer should be deployed into.</td>
		</tr>
		<tr>
			<td>issuers[0].email</td>
			<td>string</td>
			<td><pre lang="json">
"acme@example.net"
</pre>
</td>
			<td>Email to use when registering your account with Let's Encrypt.</td>
		</tr>
		<tr>
			<td>issuers[0].disableAccountKeyGeneration</td>
			<td>bool</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>Enables or disables generating a new ACME account key. If `true`, the Issuer resource will not request a new account but will expect the account key to be supplied via an existing secret. If `false`, the cert-manager system will generate a new ACME account key for the Issuer. Defaults to `false` when unspecified. See https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEIssuer</td>
		</tr>
		<tr>
			<td>issuers[0].profile</td>
			<td>string</td>
			<td><pre lang="json">
"classic"
</pre>
</td>
			<td>If the ACME server supports profiles, you can specify the profile name here. For more details, see https://cert-manager.io/docs/configuration/acme/#acme-certificate-profiles Possible profiles are: `classic`, `tlsserver`, `shortlived`.</td>
		</tr>
		<tr>
			<td>issuers[0].cnameStrategy</td>
			<td>string</td>
			<td><pre lang="json">
"None"
</pre>
</td>
			<td>Follow CNAME records or not. 2 options: - `None` (default): Don't follow CNAME records - `Follow`: Follow CNAME records See https://cert-manager.io/docs/configuration/acme/dns01/#delegated-domains-for-dns01 for more information.</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding</td>
			<td>object</td>
			<td><pre lang="json">
{
  "enabled": false,
  "keyID": "",
  "keySecretRef": {
    "key": "",
    "name": ""
  }
}
</pre>
</td>
			<td>Define the EAB (external account binding) key using `secretRef`. This is optional and only required if you want to use it. See https://cert-manager.io/docs/configuration/acme/#external-account-bindings for more information.</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding.enabled</td>
			<td>bool</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>When `enabled` is `true`, the external account binding configuration is applied.</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding.keyID</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The key ID or account ID of which your external account binding is indexed by the external account manager</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding.keySecretRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "",
  "name": ""
}
</pre>
</td>
			<td>The name and key of a secret containing a base 64 encoded URL string of your external account symmetric MAC key.</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding.keySecretRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the EAB HMAC key.</td>
		</tr>
		<tr>
			<td>issuers[0].externalAccountBinding.keySecretRef.key</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The key name in the secret above that holds the actual EAB HMAC key value.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhEndpointName</td>
			<td>string</td>
			<td><pre lang="json">
"ovh-eu"
</pre>
</td>
			<td>The endpoint name of the OVH API. It must be one of the following: `ovh-eu`, `ovh-ca`, `kimsufi-eu`, `kimsufi-ca`, `soyoustart-eu`, `soyoustart-ca`, `runabove-ca`. See https://docs.certifytheweb.com/docs/dns/providers/ovh/ for more information.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationMethod</td>
			<td>string</td>
			<td><pre lang="json">
"application"
</pre>
</td>
			<td>Define how the webhook should authenticate with the OVH API. It can be either `application` or `oauth2`. When `ovhAuthenticationMethod` is `application`, then use:   - `ovhAuthentication.applicationKey`, `ovhAuthentication.applicationSecret`, and `ovhAuthentication.applicationConsumerKey`   or   - `ovhAuthenticationRef.applicationKeyRef`, `ovhAuthenticationRef.applicationSecretRef` and `ovhAuthenticationRef.applicationConsumerKeyRef`.  When `ovhAuthenticationMethod` is `oauth2`, then use:   - `ovhAuthentication.oauth2ClientID` and `ovhAuthentication.oauth2ClientSecret`   or   - `ovhAuthenticationRef.oauth2ClientIDRef` and `ovhAuthenticationRef.oauth2ClientSecretRef`.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication</td>
			<td>object</td>
			<td><pre lang="json">
{
  "applicationConsumerKey": "",
  "applicationKey": "",
  "applicationSecret": "",
  "oauth2ClientID": "",
  "oauth2ClientSecret": ""
}
</pre>
</td>
			<td>Define your credentials below and the chart will create the necessary secret for you. If `ovhAuthenticationMethod` is `application`, then provide:   - `applicationKey`, `applicationSecret`, and `applicationConsumerKey` If `ovhAuthenticationMethod` is `oauth2`, then provide:   - `oauth2ClientID` and `oauth2ClientSecret`</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication.oauth2ClientID</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The OVH OAuth 2 client ID. Leave empty if you are using an existing secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication.oauth2ClientSecret</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The OVH OAuth 2 client secret. Leave empty if you are using an existing secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication.applicationKey</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The OVH application key. Leave empty if you are using an existing secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication.applicationSecret</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The OVH application secret. Leave empty if you are using an existing secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthentication.applicationConsumerKey</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Your OVH consumer key. Leave empty if you are using an existing secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "applicationConsumerKeyRef": {
    "key": "applicationConsumerKey",
    "name": ""
  },
  "applicationKeyRef": {
    "key": "applicationKey",
    "name": ""
  },
  "applicationSecretRef": {
    "key": "applicationSecret",
    "name": ""
  },
  "oauth2ClientIDRef": {
    "key": "oauth2ClientID",
    "name": ""
  },
  "oauth2ClientSecretRef": {
    "key": "oauth2ClientSecret",
    "name": ""
  }
}
</pre>
</td>
			<td>Define the details of a secret already containing the OVH credentials. If `ovhAuthenticationMethod` is `application`, then provide:   - `applicationKeyRef`, `applicationSecretRef` and `applicationConsumerKeyRef`. If `ovhAuthenticationMethod` is `oauth2`, then provide:   - `oauth2ClientIDRef` and `oauth2ClientSecretRef`.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientIDRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "oauth2ClientID",
  "name": ""
}
</pre>
</td>
			<td>The secret reference to an existing OVH OAuth 2 client ID.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientIDRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the OVH OAuth 2 client ID</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientIDRef.key</td>
			<td>string</td>
			<td><pre lang="json">
"oauth2ClientID"
</pre>
</td>
			<td>The key name in the secret above that holds the actual OVH OAuth 2 client ID value</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "oauth2ClientSecret",
  "name": ""
}
</pre>
</td>
			<td>The secret reference to an existing OVH Auth 2 client secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the OVH Auth 2 client secret</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef.key</td>
			<td>string</td>
			<td><pre lang="json">
"oauth2ClientSecret"
</pre>
</td>
			<td>The key name in the secret above that holds the actual OVH Auth 2 client secret value</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationKeyRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "applicationKey",
  "name": ""
}
</pre>
</td>
			<td>The secret reference to an existing OVH application key.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationKeyRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the OVH Application Key</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationKeyRef.key</td>
			<td>string</td>
			<td><pre lang="json">
"applicationKey"
</pre>
</td>
			<td>The key name in the secret above that holds the actual OVH application key value</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationSecretRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "applicationSecret",
  "name": ""
}
</pre>
</td>
			<td>The secret reference to an existing OVH application secret.</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationSecretRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the OVH Application Secret</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationSecretRef.key</td>
			<td>string</td>
			<td><pre lang="json">
"applicationSecret"
</pre>
</td>
			<td>The key name in the secret above that holds the actual OVH application secret value</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef</td>
			<td>object</td>
			<td><pre lang="json">
{
  "key": "applicationConsumerKey",
  "name": ""
}
</pre>
</td>
			<td>The secret reference to an existing OVH consumer key</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef.name</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Name of the Kubernetes secret containing the OVH consumer Key</td>
		</tr>
		<tr>
			<td>issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef.key</td>
			<td>string</td>
			<td><pre lang="json">
"applicationConsumerKey"
</pre>
</td>
			<td>The key name in the secret above that holds the actual OVH consumer key value</td>
		</tr>
		<tr>
			<td>rbac</td>
			<td>object</td>
			<td><pre lang="json">
{
  "roleType": "Role"
}
</pre>
</td>
			<td>RBAC configuration for the webhook deployment</td>
		</tr>
		<tr>
			<td>rbac.roleType</td>
			<td>string</td>
			<td><pre lang="json">
"Role"
</pre>
</td>
			<td>If the secret of the issuer is in a different namespace, change `Role` (default) to `ClusterRole` so the webhook is able to access this secret.</td>
		</tr>
		<tr>
			<td>image.registry</td>
			<td>string</td>
			<td><pre lang="json">
"ghcr.io"
</pre>
</td>
			<td>The registry to use.</td>
		</tr>
		<tr>
			<td>image.repository</td>
			<td>string</td>
			<td><pre lang="json">
"aureq/cert-manager-webhook-ovh"
</pre>
</td>
			<td>The repository location on the registry.</td>
		</tr>
		<tr>
			<td>image.tag</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The tag to use from the registry. It can either be a version number like `1.0.0` or a digest like `sha256:abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789`.</td>
		</tr>
		<tr>
			<td>image.pullPolicy</td>
			<td>string</td>
			<td><pre lang="json">
"IfNotPresent"
</pre>
</td>
			<td>The image pull policy as defined in the Kubernetes documentation.</td>
		</tr>
		<tr>
			<td>image.pullSecrets</td>
			<td>list</td>
			<td><pre lang="json">
[]
</pre>
</td>
			<td>The image pull secrets to use if your registry requires it.</td>
		</tr>
		<tr>
			<td>pod.annotations</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Extra annotations for the Pod spec. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ for more information.</td>
		</tr>
		<tr>
			<td>pod.podLabels</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Labels to add to the Pod. See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/</td>
		</tr>
		<tr>
			<td>pod.podAnnotations</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Annotations to add to the Pod. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/</td>
		</tr>
		<tr>
			<td>pod.replicas</td>
			<td>int</td>
			<td><pre lang="json">
1
</pre>
</td>
			<td>Number of replicas in this deployment.</td>
		</tr>
		<tr>
			<td>pod.environment</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Use this option to add environment variables relevant to your deployment. These fields will be passed on to the container when Chart is deployed.</td>
		</tr>
		<tr>
			<td>pod.securityContext</td>
			<td>object</td>
			<td><pre lang="json">
{
  "container": {
    "allowPrivilegeEscalation": false,
    "privileged": false,
    "readOnlyRootFilesystem": true
  },
  "pod": {
    "runAsGroup": 1000,
    "runAsUser": 1000
  }
}
</pre>
</td>
			<td>Define the security context for both the `Pod` and the `Container`. See the Kubernetes API [reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/latest/#securitycontext-v1-core) and the Kubernetes [documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for more information.</td>
		</tr>
		<tr>
			<td>pod.securityContext.pod</td>
			<td>object</td>
			<td><pre lang="json">
{
  "runAsGroup": 1000,
  "runAsUser": 1000
}
</pre>
</td>
			<td>Pod Security Context.</td>
		</tr>
		<tr>
			<td>pod.securityContext.container</td>
			<td>object</td>
			<td><pre lang="json">
{
  "allowPrivilegeEscalation": false,
  "privileged": false,
  "readOnlyRootFilesystem": true
}
</pre>
</td>
			<td>Container Security Context.</td>
		</tr>
		<tr>
			<td>pod.port</td>
			<td>int</td>
			<td><pre lang="json">
8443
</pre>
</td>
			<td>The port used by the webhook to receive incoming traffic from the service. It may be useful to change it if the hostNetwork mode is activated to use an available port.</td>
		</tr>
		<tr>
			<td>pod.hostNetwork</td>
			<td>bool</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>Specifies if the webhook should be started in `hostNetwork` mode. Required for use in some managed Kubernetes clusters (such as AWS EKS) with custom CNI (such as Calico), because control-plane managed by AWS cannot communicate with pods' IP CIDR and admission webhooks are not working.</td>
		</tr>
		<tr>
			<td>pod.priorityClassName</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>The optional priority class to be used for the cert-manager-webhook-ovh pod.</td>
		</tr>
		<tr>
			<td>pod.resources</td>
			<td>object</td>
			<td><pre lang="json">
{
  "limits": {},
  "requests": {}
}
</pre>
</td>
			<td>We usually recommend not to specify default resources and to leave this as a conscious choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, adjust them as necessary as shown below.</td>
		</tr>
		<tr>
			<td>pod.resources.requests</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Resource Requests.</td>
		</tr>
		<tr>
			<td>pod.resources.limits</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Resource Limits.</td>
		</tr>
		<tr>
			<td>pod.selectors.nodeSelector</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Node selector. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/ for more information.</td>
		</tr>
		<tr>
			<td>pod.selectors.affinity</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Affinity selectors. See https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/ for more information.</td>
		</tr>
		<tr>
			<td>nameOverride</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Override the name of the deployed resources</td>
		</tr>
		<tr>
			<td>fullnameOverride</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Override the full name of the deployed resources</td>
		</tr>
		<tr>
			<td>service</td>
			<td>object</td>
			<td><pre lang="json">
{
  "annotations": {},
  "port": 443,
  "type": "ClusterIP"
}
</pre>
</td>
			<td>Service configuration for the webhook deployment.</td>
		</tr>
		<tr>
			<td>service.type</td>
			<td>string</td>
			<td><pre lang="json">
"ClusterIP"
</pre>
</td>
			<td>The service type to use.</td>
		</tr>
		<tr>
			<td>service.port</td>
			<td>int</td>
			<td><pre lang="json">
443
</pre>
</td>
			<td>The port used by the service.</td>
		</tr>
		<tr>
			<td>service.annotations</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Extra annotations for the service. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ for more information.</td>
		</tr>
		<tr>
			<td>extraObjects</td>
			<td>list</td>
			<td><pre lang="json">
[]
</pre>
</td>
			<td>Create dynamic manifests via values.  For example: extraObjects:   - |     apiVersion: v1     kind: ConfigMap     metadata:       name: '{{ template "cert-manager.fullname" . }}-extra-configmap'</td>
		</tr>
	</tbody>
</table>

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

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Aurélien Requiem |  | <https://github.com/aureq> |

## Source Code

* <https://github.com/aureq/cert-manager-webhook-ovh>
