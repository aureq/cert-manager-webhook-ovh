# cert-manager-webhook-ovh

![Version: 0.9.0-alpha.0](https://img.shields.io/badge/Version-0.9.0--alpha.0-informational?style=flat-square) ![AppVersion: 0.9.0-alpha.0](https://img.shields.io/badge/AppVersion-0.9.0--alpha.0-informational?style=flat-square)

OVH DNS cert-manager ACME webhook

**Homepage:** <https://github.com/aureq/cert-manager-webhook-ovh>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Aurélien Requiem |  | <https://github.com/aureq> |

## Source Code

* <https://github.com/aureq/cert-manager-webhook-ovh>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| groupName | string | `"acme.mycompany.example"` | The `groupName` here is used to identify your company or business unit that created this webhook. For example, this may be `acme.mycompany.example`. This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain here is recommended. |
| certManager | object | `{"namespace":"cert-manager","serviceAccountName":"cert-manager"}` | The values below should match your cert-manager deployment. If not, you will get permissions errors. |
| certManager.namespace | string | `"cert-manager"` | namespace in which your cert-manager is deployed |
| certManager.serviceAccountName | string | `"cert-manager"` | cert-manager serverAccount name (default: cert-manager) |
| issuers | list | `[{"acmeServerUrl":"https://acme-v02.api.letsencrypt.org/directory","cnameStrategy":"None","create":false,"disableAccountKeyGeneration":false,"email":"acme@example.net","externalAccountBinding":{"enabled":false,"keyID":"","keySecretRef":{"key":"","name":""}},"kind":"ClusterIssuer","name":"le-prod","namespace":"default","ovhAuthentication":{"applicationConsumerKey":"","applicationKey":"","applicationSecret":"","oauth2ClientID":"","oauth2ClientSecret":""},"ovhAuthenticationMethod":"application","ovhAuthenticationRef":{"applicationConsumerKeyRef":{"key":"applicationConsumerKey","name":""},"applicationKeyRef":{"key":"applicationKey","name":""},"applicationSecretRef":{"key":"applicationSecret","name":""},"oauth2ClientIDRef":{"key":"oauth2ClientID","name":""},"oauth2ClientSecretRef":{"key":"oauth2ClientSecret","name":""}},"ovhEndpointName":"ovh-eu","profile":"classic"}]` | One or more issuers to create |
| issuers[0] | object | `{"acmeServerUrl":"https://acme-v02.api.letsencrypt.org/directory","cnameStrategy":"None","create":false,"disableAccountKeyGeneration":false,"email":"acme@example.net","externalAccountBinding":{"enabled":false,"keyID":"","keySecretRef":{"key":"","name":""}},"kind":"ClusterIssuer","name":"le-prod","namespace":"default","ovhAuthentication":{"applicationConsumerKey":"","applicationKey":"","applicationSecret":"","oauth2ClientID":"","oauth2ClientSecret":""},"ovhAuthenticationMethod":"application","ovhAuthenticationRef":{"applicationConsumerKeyRef":{"key":"applicationConsumerKey","name":""},"applicationKeyRef":{"key":"applicationKey","name":""},"applicationSecretRef":{"key":"applicationSecret","name":""},"oauth2ClientIDRef":{"key":"oauth2ClientID","name":""},"oauth2ClientSecretRef":{"key":"oauth2ClientSecret","name":""}},"ovhEndpointName":"ovh-eu","profile":"classic"}` | Name of this issuer |
| issuers[0].create | bool | `false` | When true this issuer is created. |
| issuers[0].kind | string | `"ClusterIssuer"` | The type of issuer to create. It can either be `ClusterIssuer` or `Issuer`. If `Issuer` is specified, then the `namespace` is also required. See https://cert-manager.io/docs/concepts/issuer/ for more information. |
| issuers[0].namespace | string | `"default"` | If `kind` is `Issuer`, then indicate the namespace in which this issuer should be deployed into. |
| issuers[0].email | string | `"acme@example.net"` | Email to use when registering your account with Let's Encrypt. |
| issuers[0].disableAccountKeyGeneration | bool | `false` | Enables or disables generating a new ACME account key. If `true`, the Issuer resource will not request a new account but will expect the account key to be supplied via an existing secret. If `false`, the cert-manager system will generate a new ACME account key for the Issuer. Defaults to `false` when unspecified. See https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEIssuer |
| issuers[0].profile | string | `"classic"` | If the ACME server supports profiles, you can specify the profile name here. For more details, see https://cert-manager.io/docs/configuration/acme/#acme-certificate-profiles Possible profiles are: `classic`, `tlsserver`, `shortlived`. |
| issuers[0].cnameStrategy | string | `"None"` | Follow CNAME records or not. 2 options: - `None` (default): Don't follow CNAME records - `Follow`: Follow CNAME records See https://cert-manager.io/docs/configuration/acme/dns01/#delegated-domains-for-dns01 for more information. |
| issuers[0].externalAccountBinding | object | `{"enabled":false,"keyID":"","keySecretRef":{"key":"","name":""}}` | Define the EAB (external account binding) key using `secretRef`. This is optional and only required if you want to use it. See https://cert-manager.io/docs/configuration/acme/#external-account-bindings for more information. |
| issuers[0].externalAccountBinding.enabled | bool | `false` | When `enabled` is `true`, the external account binding configuration is applied. |
| issuers[0].externalAccountBinding.keyID | string | `""` | The key ID or account ID of which your external account binding is indexed by the external account manager |
| issuers[0].externalAccountBinding.keySecretRef | object | `{"key":"","name":""}` | The name and key of a secret containing a base 64 encoded URL string of your external account symmetric MAC key. |
| issuers[0].externalAccountBinding.keySecretRef.name | string | `""` | Name of the Kubernetes secret containing the EAB HMAC key. |
| issuers[0].externalAccountBinding.keySecretRef.key | string | `""` | The key name in the secret above that holds the actual EAB HMAC key value. |
| issuers[0].ovhEndpointName | string | `"ovh-eu"` | The endpoint name of the OVH API. It must be one of the following: `ovh-eu`, `ovh-ca`, `kimsufi-eu`, `kimsufi-ca`, `soyoustart-eu`, `soyoustart-ca`, `runabove-ca`. See https://docs.certifytheweb.com/docs/dns/providers/ovh/ for more information. |
| issuers[0].ovhAuthenticationMethod | string | `"application"` | Define how the webhook should authenticate with the OVH API. It can be either `application` or `oauth2`. When `ovhAuthenticationMethod` is `application`, then use:   - `ovhAuthentication.applicationKey`, `ovhAuthentication.applicationSecret`, and `ovhAuthentication.applicationConsumerKey`   or   - `ovhAuthenticationRef.applicationKeyRef`, `ovhAuthenticationRef.applicationSecretRef` and `ovhAuthenticationRef.applicationConsumerKeyRef`.  When `ovhAuthenticationMethod` is `oauth2`, then use:   - `ovhAuthentication.oauth2ClientID` and `ovhAuthentication.oauth2ClientSecret`   or   - `ovhAuthenticationRef.oauth2ClientIDRef` and `ovhAuthenticationRef.oauth2ClientSecretRef`. |
| issuers[0].ovhAuthentication | object | `{"applicationConsumerKey":"","applicationKey":"","applicationSecret":"","oauth2ClientID":"","oauth2ClientSecret":""}` | Define your credentials below and the chart will create the necessary secret for you. If `ovhAuthenticationMethod` is `application`, then provide:   - `applicationKey`, `applicationSecret`, and `applicationConsumerKey` If `ovhAuthenticationMethod` is `oauth2`, then provide:   - `oauth2ClientID` and `oauth2ClientSecret` |
| issuers[0].ovhAuthentication.oauth2ClientID | string | `""` | The OVH OAuth 2 client ID. Leave empty if you are using an existing secret. |
| issuers[0].ovhAuthentication.oauth2ClientSecret | string | `""` | The OVH OAuth 2 client secret. Leave empty if you are using an existing secret. |
| issuers[0].ovhAuthentication.applicationKey | string | `""` | The OVH application key. Leave empty if you are using an existing secret. |
| issuers[0].ovhAuthentication.applicationSecret | string | `""` | The OVH application secret. Leave empty if you are using an existing secret. |
| issuers[0].ovhAuthentication.applicationConsumerKey | string | `""` | Your OVH consumer key. Leave empty if you are using an existing secret. |
| issuers[0].ovhAuthenticationRef | object | `{"applicationConsumerKeyRef":{"key":"applicationConsumerKey","name":""},"applicationKeyRef":{"key":"applicationKey","name":""},"applicationSecretRef":{"key":"applicationSecret","name":""},"oauth2ClientIDRef":{"key":"oauth2ClientID","name":""},"oauth2ClientSecretRef":{"key":"oauth2ClientSecret","name":""}}` | Define the details of a secret already containing the OVH credentials. If `ovhAuthenticationMethod` is `application`, then provide:   - `applicationKeyRef`, `applicationSecretRef` and `applicationConsumerKeyRef`. If `ovhAuthenticationMethod` is `oauth2`, then provide:   - `oauth2ClientIDRef` and `oauth2ClientSecretRef`. |
| issuers[0].ovhAuthenticationRef.oauth2ClientIDRef | object | `{"key":"oauth2ClientID","name":""}` | The secret reference to an existing OVH OAuth 2 client ID. |
| issuers[0].ovhAuthenticationRef.oauth2ClientIDRef.name | string | `""` | Name of the Kubernetes secret containing the OVH OAuth 2 client ID |
| issuers[0].ovhAuthenticationRef.oauth2ClientIDRef.key | string | `"oauth2ClientID"` | The key name in the secret above that holds the actual OVH OAuth 2 client ID value |
| issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef | object | `{"key":"oauth2ClientSecret","name":""}` | The secret reference to an existing OVH Auth 2 client secret. |
| issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef.name | string | `""` | Name of the Kubernetes secret containing the OVH Auth 2 client secret |
| issuers[0].ovhAuthenticationRef.oauth2ClientSecretRef.key | string | `"oauth2ClientSecret"` | The key name in the secret above that holds the actual OVH Auth 2 client secret value |
| issuers[0].ovhAuthenticationRef.applicationKeyRef | object | `{"key":"applicationKey","name":""}` | The secret reference to an existing OVH application key. |
| issuers[0].ovhAuthenticationRef.applicationKeyRef.name | string | `""` | Name of the Kubernetes secret containing the OVH Application Key |
| issuers[0].ovhAuthenticationRef.applicationKeyRef.key | string | `"applicationKey"` | The key name in the secret above that holds the actual OVH application key value |
| issuers[0].ovhAuthenticationRef.applicationSecretRef | object | `{"key":"applicationSecret","name":""}` | The secret reference to an existing OVH application secret. |
| issuers[0].ovhAuthenticationRef.applicationSecretRef.name | string | `""` | Name of the Kubernetes secret containing the OVH Application Secret |
| issuers[0].ovhAuthenticationRef.applicationSecretRef.key | string | `"applicationSecret"` | The key name in the secret above that holds the actual OVH application secret value |
| issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef | object | `{"key":"applicationConsumerKey","name":""}` | The secret reference to an existing OVH consumer key |
| issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef.name | string | `""` | Name of the Kubernetes secret containing the OVH consumer Key |
| issuers[0].ovhAuthenticationRef.applicationConsumerKeyRef.key | string | `"applicationConsumerKey"` | The key name in the secret above that holds the actual OVH consumer key value |
| rbac | object | `{"roleType":"Role"}` | RBAC configuration for the webhook deployment |
| rbac.roleType | string | `"Role"` | If the secret of the issuer is in a different namespace, change `Role` (default) to `ClusterRole` so the webhook is able to access this secret. |
| image.registry | string | `"ghcr.io"` | The registry to use. |
| image.repository | string | `"aureq/cert-manager-webhook-ovh"` | The repository location on the registry. |
| image.tag | string | `""` | The tag to use from the registry. It can either be a version number like `1.0.0` or a digest like `sha256:abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789`. |
| image.pullPolicy | string | `"IfNotPresent"` | The image pull policy as defined in the Kubernetes documentation. |
| image.pullSecrets | list | `[]` | The image pull secrets to use if your registry requires it. |
| pod.annotations | object | `{}` | Extra annotations for the Pod spec. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ for more information. |
| pod.podLabels | object | `{}` | Labels to add to the Pod. See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ |
| pod.podAnnotations | object | `{}` | Annotations to add to the Pod. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ |
| pod.replicas | int | `1` | Number of replicas in this deployment. |
| pod.environment | object | `{}` | Use this option to add environment variables relevant to your deployment. These fields will be passed on to the container when Chart is deployed. |
| pod.securityContext | object | `{"container":{"allowPrivilegeEscalation":false,"privileged":false,"readOnlyRootFilesystem":true},"pod":{"runAsGroup":1000,"runAsUser":1000}}` | Define the security context for both the `Pod` and the `Container`. See the Kubernetes API [reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/latest/#securitycontext-v1-core) and the Kubernetes [documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for more information. |
| pod.securityContext.pod | object | `{"runAsGroup":1000,"runAsUser":1000}` | Pod Security Context. |
| pod.securityContext.container | object | `{"allowPrivilegeEscalation":false,"privileged":false,"readOnlyRootFilesystem":true}` | Container Security Context. |
| pod.port | int | `8443` | The port used by the webhook to receive incoming traffic from the service. It may be useful to change it if the hostNetwork mode is activated to use an available port. |
| pod.hostNetwork | bool | `false` | Specifies if the webhook should be started in `hostNetwork` mode. Required for use in some managed Kubernetes clusters (such as AWS EKS) with custom CNI (such as Calico), because control-plane managed by AWS cannot communicate with pods' IP CIDR and admission webhooks are not working. |
| pod.priorityClassName | string | `""` | The optional priority class to be used for the cert-manager-webhook-ovh pod. |
| pod.resources | object | `{"limits":{},"requests":{}}` | We usually recommend not to specify default resources and to leave this as a conscious choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, adjust them as necessary as shown below. |
| pod.resources.requests | object | `{}` | Resource Requests. |
| pod.resources.limits | object | `{}` | Resource Limits. |
| pod.selectors.nodeSelector | object | `{}` | Node selector. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/ for more information. |
| pod.selectors.affinity | object | `{}` | Affinity selectors. See https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/ for more information. |
| nameOverride | string | `""` | Override the name of the deployed resources |
| fullnameOverride | string | `""` | Override the full name of the deployed resources |
| service | object | `{"annotations":{},"port":443,"type":"ClusterIP"}` | Service configuration for the webhook deployment. |
| service.type | string | `"ClusterIP"` | The service type to use. |
| service.port | int | `443` | The port used by the service. |
| service.annotations | object | `{}` | Extra annotations for the service. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ for more information. |
| extraObjects | list | `[]` | Create dynamic manifests via values.  For example: extraObjects:   - |     apiVersion: v1     kind: ConfigMap     metadata:       name: '{{ template "cert-manager.fullname" . }}-extra-configmap' |

