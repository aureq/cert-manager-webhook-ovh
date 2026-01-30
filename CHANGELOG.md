# Changelog

## 0.8.1-alpha.1

### Noteworthy changes

- 🐛 fix ACME challenge subdomain parsing for multi-level subdomains. [#75](https://github.com/aureq/cert-manager-webhook-ovh/pull/75) (by [Karol Stoiński](https://github.com/KarolStoinski))
- 🌿 add helm schema support
- 📝 add helpful post-install notes to Helm chart
- 📄 improved documentation so readers have instructions on how to perform the setup with the [OVH cli](https://github.com/ovh/ovhcloud-cli)

### Dependencies

- ⏩ upgrade alpine to 3.23
- ⏩ upgrade golang to 1.25.5
- ⏩ upgrade golang toolchain to 1.25.5
- ⏩ upgrade github.com/cert-manager/cert-manager to v1.19.2
- ⏩ upgrade k8s.io/api to v0.34.3
- ⏩ upgrade k8s.io/apiextensions-apiserver to v0.34.3
- ⏩ upgrade k8s.io/apimachinery to v0.34.3
- ⏩ upgrade k8s.io/client-go to v0.34.3

## 0.8.1-alpha.0

### Noteworthy changes

- 🐛 trim `"` around TXT records when checking value before deciding to delete (by [flodakto](https://github.com/flodakto))
- 🌿 add warning log when skipping TXT record deletion due to values mismatch
- 🌿 improve error logging throughout the entire webhook
- 🌿 add charmbracelet/log v0.4.2 to perform structured logging functions

### Dependencies

- ⏩ add charmbracelet/log v0.4.2
- ⏩ upgrade golang to 1.25.0
- ⏩ upgrade golang toolchain to 1.25.0
- ⏩ upgrade github.com/cert-manager/cert-manager to v1.19.1
- ⏩ upgrade k8s.io/api to v0.34.1
- ⏩ upgrade k8s.io/apimachinery to v0.34.1
- ⏩ upgrade k8s.io/client-go to v0.34.1

## 0.8.0

### Major features

- 🎉 add support for OAuth2 authentication when communicating with OVH API (by [Rémy Jacquin](https://github.com/remyj38))
- 🎉 add support for ACME profiles paving the way to shortlived certificates ([docs](https://letsencrypt.org/docs/profiles/))

### Breaking changes

- ⚠️ renamed configuration key `consumerKey` to `applicationConsumerKey` to prevent confusion with OAuth2 authentication. See [documentation](https://aureq.github.io/cert-manager-webhook-ovh/#configuration).
- ⚠️ renamed `ConsumerKeyRef` to `ApplicationConsumerKeyRef` to prevent confusion with OAuth2 authentication. See [documentation](https://aureq.github.io/cert-manager-webhook-ovh/#configuration).
- ⚠️ due to the breaking changes described above, the new value for `configVersion` field is `0.0.2`.

### Noteworthy changes

- 🎉 add unit tests to validate the Chart (by [Rémy Jacquin](https://github.com/remyj38))
- 🎉 add unit tests to validate the webhook (by [Rémy Jacquin](https://github.com/remyj38))
- 🌿 add 4 end-to-end tests to fully test the webhook (combination of app/oauth2 and direct/secret references)
- 🌿 improve some error messages when `configVersion` is invalid or missing
- 📄 document how to configure and use the new OAuth2 authentication (by [Rémy Jacquin](https://github.com/remyj38), with modification from [Aurélien Requiem](https://github.com/aureq))
- 📄 expand documentation on how to create IAM policy
- 🐛 add missing permissions for the IAM policy

### Dependencies

- ⏩ upgrade golang to 1.24.7
- ⏩ upgrade golang toolchain to 1.24.7

## 0.8.0-alpha.1

### Major features

- 🎉 add support for OAuth2 authentication when communicating with OVH API (by [Rémy Jacquin](https://github.com/remyj38))

### Breaking changes

- ⚠️ renamed configuration key `consumerKey` to `applicationConsumerKey` to prevent confusion with OAuth2 authentication. See [documentation](https://aureq.github.io/cert-manager-webhook-ovh/#configuration).
- ⚠️ renamed `ConsumerKeyRef` to `ApplicationConsumerKeyRef` to prevent confusion with OAuth2 authentication. See [documentation](https://aureq.github.io/cert-manager-webhook-ovh/#configuration).
- ⚠️ due to the breaking changes described above, the new value for `configVersion` field is `0.0.2`.

### Noteworthy changes

- 🎉 add unit tests to validate the Chart (by [Rémy Jacquin](https://github.com/remyj38))
- 🎉 add unit tests to validate the webhook (by [Rémy Jacquin](https://github.com/remyj38))
- 🌿 improve some error messages when `configVersion` is invalid or missing
- 📄 document how to configure and use the new OAuth2 authentication (by [Rémy Jacquin](https://github.com/remyj38), with modification from [Aurélien Requiem](https://github.com/aureq))

### Dependencies

- ⏩ upgrade golang to 1.24.7
- ⏩ upgrade golang toolchain to 1.24.7

## 0.7.6

### Noteworthy changes

- 🎉 both the build image and base image are using Alpine Linux 3.22
- 🎉 allow creation of extra manifests via `values.yaml` (by [Rémy Jacquin](https://github.com/remyj38))
- 🐛 fix `app.kubernetes.io/version` label when using SHA in image tag (by [Rémy Jacquin](https://github.com/remyj38))

### Dependencies

- ⏩ upgrade golang build image to 1.24-alpine3.22
- ⏩ upgrade alpine base image to 3.22
- ⏩ upgrade github.com/cert-manager/cert-manager to v1.18.2
- ⏩ upgrade github.com/ovh/go-ovh to v1.9.0

## 0.7.5

### Noteworthy changes

- 🎉 add support for External Account Binding (by [Sebastien MALOT](https://github.com/smalot))

### Dependencies

- ⏩ upgrade github.com/cert-manager/cert-manager to v1.17.3
- ⏩ upgrade github.com/ovh/go-ovh to v1.8.0
- ⏩ upgrade k8s.io/api to v0.33.2
- ⏩ upgrade k8s.io/apimachinery to v0.33.2
- ⏩ upgrade k8s.io/client-go to v0.33.2

## 0.7.4

### Noteworthy changes

- 🧹 maintenance release, updated dependenies only.
- 🐛 fix minor casing issue in `Dockerfile`
- 🙈 ignore .aider* files

### Dependencies

- ⏩ upgrade github.com/cert-manager/cert-manager to v1.17.2
- ⏩ upgrade github.com/ovh/go-ovh to v1.7.0
- ⏩ upgrade k8s.io/api to v0.32.5
- ⏩ upgrade k8s.io/apiextensions-apiserver to v0.32.5
- ⏩ upgrade k8s.io/apimachinery to v0.32.5
- ⏩ upgrade k8s.io/client-go to v0.32.5
- ⏩ upgrade golang to 1.24
- ⏩ upgrade alpine to 3.21

## 0.7.3

### Noteworthy changes

- 🧹 maintenance release, only updated dependenies.

### Dependencies

- ⏩ update golang.org/x/net v0.33.0 to address [CVE-2024-45338](https://github.com/advisories/GHSA-w32m-9786-jp63)

## 0.7.2

### Noteworthy changes

- 🧹 maintenance release, only updated dependenies.

### Dependencies

- ⏩ update golang.org/x/crypto v0.31.0

## 0.7.1

### Noteworthy changes

- 🧹 maintenance release, only updated dependenies.

### Dependencies

- ⏩ update go 1.23.3
- ⏩ update github.com/cert-manager/cert-manager v1.16.2
- ⏩ update k8s.io/api v0.31.3
- ⏩ update k8s.io/apiextensions-apiserver v0.31.3
- ⏩ update k8s.io/apimachinery v0.31.3
- ⏩ update k8s.io/client-go v0.31.3

## 0.7.0

### Noteworthy changes

- ✨ Add new `configVersion` to assist with breaking change
- 🌿 Prefix Helm Chart error messages with 'Error:'
- 🐛 Address minor typography issues in documentation.
- 🌿 support adding customer labels to pod
- 📄 slightly improve documentation in values.yaml

### Dependencies

- ⏩ Use Alpine to 3.20 and Golang 1.23 as build image
- ⏩ Use Alpine to 3.20 as base image
- ⏩ Use Go 1.23.0 to build webhook
- ⏩ Bump github.com/cert-manager/cert-manager 1.14.1 to 1.15.3
- ⏩ Bump github.com/ovh/go-ovh from 1.4.3 to 1.6.0
- ⏩ Bump k8s.io/api from 0.29.1 to 0.30.1
- ⏩ Bump k8s.io/apiextensions-apiserver from 0.29.1 to 0.30.1
- ⏩ Bump golang.org/x/net from 0.20.0 to 0.23.0

## 0.7.0-alpha.3

### Noteworthy changes

- 🌿 rename schemaVersion to configVersion
- 🐛 fix error when handling commented configVersion
- 🌿 improve version check

## 0.7.0-alpha.2

### Noteworthy changes

- 🌿 support adding customer labels to pod
- 📄 slightly improve documentation in values.yaml

## 0.7.0-alpha.1

### Noteworthy changes

- ✨ Add new `schemaVersion` to assist with breaking change
- 🌿 Prefix error messages with 'Error:'
- 🐛 Address minor typography issues in documentation.

### Dependencies

- ⏩ Use Alpine to 3.20 and Golang 1.23 as build image
- ⏩ Use Alpine to 3.20 as base image
- ⏩ Use Go 1.23.0 to build webhook
- ⏩ Bump github.com/cert-manager/cert-manager 1.14.1 to 1.15.3
- ⏩ Bump github.com/ovh/go-ovh from 1.4.3 to 1.6.0
- ⏩ Bump k8s.io/api from 0.29.1 to 0.30.1
- ⏩ Bump k8s.io/apiextensions-apiserver from 0.29.1 to 0.30.1
- ⏩ Bump golang.org/x/net from 0.20.0 to 0.23.0

## 0.6.0

### Noteworthy changes

- ⚠️ Separate `securityContext` for both `container` and `pod`. See `values.yaml` for more details. See [#32](https://github.com/aureq/cert-manager-webhook-ovh/pull/32). Authored by [Mathieu Sensei](https://github.com/hyu9a).
- ✨ Support `podAnnotations`. See [#32](https://github.com/aureq/cert-manager-webhook-ovh/pull/32). Authored by [Mathieu Sensei](https://github.com/hyu9a).
- 🌿 Comment out `image.tag` as it's not needed unless someone wants to override the container image version

### Dependencies

- ⏩ Use Alpine to 3.19.1 as base image
- ⏩ Use Go 1.21.6 to build webhook
- ⏩ Bump github.com/cert-manager/cert-manager 1.13.0 to 1.14.1
- ⏩ Bump github.com/ovh/go-ovh from 1.4.2 to 1.4.3
- ⏩ Bump golang.org/x/crypto from 0.14.0 to 0.18.0
- ⏩ Bump golang.org/x/net from 0.17.0 to 0.20.0
- ⏩ Bump k8s.io/api from 0.29.0 to 0.29.1
- ⏩ Bump k8s.io/apiextensions-apiserver from 0.29.0 to 0.29.1

## 0.5.2

### Dependencies

- ⏩ Bump google.golang.org/grpc from 1.58.2 to 1.58.3. See [Dependabot](https://github.com/aureq/cert-manager-webhook-ovh/pull/34)
- ⏩ Bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.35.0 to 0.46.0. See [Dependabot](https://github.com/aureq/cert-manager-webhook-ovh/pull/35)
- ⏩ Bump go.opentelemetry.io/otel/exporters/otlp/otlptrace from 1.19.0 to 1.20.0
- ⏩ Bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc from 1.19.0 to 1.20.0
- ⏩ Bump go.opentelemetry.io/otel/sdk from 1.19.0 to 1.20.0
- ⏩ Bump golang.org/x/sys from 0.13.0 to 0.14.0

### Known issues

- 🔥 Alpine 3.18.4 is vulnerable to the following CVEs. Should be fixed in [3.18.5 release](https://gitlab.alpinelinux.org/groups/alpine/-/milestones/5#tab-issues).
  - [CVE-2023-5363](https://avd.aquasec.com/nvd/cve-2023-5363)
  - [CVE-2023-5678](https://avd.aquasec.com/nvd/cve-2023-5678)
  - [CVE-2023-5363](https://avd.aquasec.com/nvd/cve-2023-5363)
  - [CVE-2023-5678](https://avd.aquasec.com/nvd/cve-2023-5678)

## 0.5.1

### Dependencies

- ⏩ bump go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp to v0.45.0 to address CVE-2023-45142. See [Dependabot](https://github.com/aureq/cert-manager-webhook-ovh/security/dependabot/6).
- ⏩ bump golang.org/x/net from 0.15.0 to 0.17.0. See [Dependabot PR](https://github.com/aureq/cert-manager-webhook-ovh/pull/33).

## 0.5.0

### Noteworthy changes

- ✨ add support for `readOnlyRootFilesystem` on the deployment (thanks @Benzhaomin)
- ✨ add Deployment annotation support (thanks @Benzhaomin)
- ✨ add ref link about `nodeSelector`, `tolerations`, `affinity` and `annotations`
- ✨ choose rbac role type (default `Role`) (thanks @Alissia01)
- 📄 document 3 more configuration entries in `values.yaml`
- 🌿 make this chart compatible with helm 3 by settings `apiVersion` to `v2`
- 🌿 drop `v` in `appVersion` and `version` fields, set `"0.5.0"`
- 🌿 udpate `image.tag` value to use SemVer 2.0 and set its values to `"0.5.0"`
- 🐛 typo fix
- ⏩ update k8s.io/apiserver to v0.28.2 due to security (dependabot)

### Dependencies

- ⏩ upgrade github.com/cert-manager/cert-manager to v1.13.0
- ⏩ build with go 1.20
- ⏩ upgrade k8s dependencies to 0.28.1
- ⏩ use alpine 3.18 as base image
- ⏩ update dependency for github.com/ovh/go-ovh to v1.4.2
- ⏩ Bump google.golang.org/grpc from 1.51.0 to 1.53.0

## 0.5.0-alpha.2

### Noteworthy changes

- ✨ add support for `readOnlyRootFilesystem` on the deployment (thanks @Benzhaomin)
- 🐛 typo fix
- ✨ add annotation support (thanks @Benzhaomin)
- ✨ add ref link about `nodeSelector`, `tolerations`, `affinity` and `annotations`
- ✨ choose rbac role type (default `Role`)
- ⏩ build with go 1.20
- ⏩ upgrade k8s dependencies to 0.28.1
- ⏩ upgrade github.com/cert-manager/cert-manager to v1.13.0
- ⏩ use alpine 3.18 as base image

## 0.5.0-alpha.1

### Noteworthy changes

- ⏩ Bump google.golang.org/grpc from 1.51.0 to 1.53.0
- 📄 document 3 more configuration entries in `values.yaml`
- 🌿 make this chart compatible with helm 3 by settings `apiVersion` to `v2`
- 🌿 drop `v` in `appVersion` and `version` fields, set `"0.5.0"`
- 🌿 udpate `image.tag` value to use SemVer 2.0 and set its values to `"0.5.0"`
- ⏩ update dependency for github.com/ovh/go-ovh to v1.4.2

## v0.4.2

### Noteworthy changes

- ✨ build images for amd64, arm64 and armv7 architectures
- 🐙 add issue templates for bugs and feature requests
- 🤖 configure dependabot to get alerts on vulnerabilities
- 📄 add disclaimer about support and code of conduct
- ✨ integration with [artifacthub.io](https://artifacthub.io/packages/helm/cert-manager-webhook-ovh/cert-manager-webhook-ovh)
- 📄 minor inconsistency fix in README.md
- 📄 add steps to make a release
- ⏩ update cert-manager dependency to v1.11.0
- ⏩ update k8s dependency to v0.26.0
- ⏩ build image using Go 1.19.7
- ⏩ upgrade alpine to 3.17
- ⏩ update Chart.yaml and `values.yaml` to use latest container image

## v0.4.2-alpha.1

### Noteworthy changes

- 📄 minor consistency fix in README.md
- ✨ start work to integrade with artifacthub.io

## v0.4.2-alpha.0

### Noteworthy changes

- ⏩ update cert-manager dependency to v1.11.0
- ⏩ update k8s dependency to v0.26.0
- ✨ build image using Go 1.19.5
- ✨ initial work to build arm64 and armv7 images

## v0.4.1

### Noteworthy changes

- 🐛 include minutes and seconds in certificates duration fields. see [argoproj/argo-cd#6008](https://github.com/argoproj/argo-cd/issues/6008) for details. via [@aegaeonit](https://github.com/aegaeonit)
- ✨ optimize Dockerfile for better builds
- ✨ explicitly use Alpine 3.16 throughout the Dockerfile
- ✨ run the webhook as `nobody`/`nogroup`
- ✨ reduce container image from 107MB down to 56.2MB
- ✨ add CNAME strategy to issuers in [#8](https://github.com/aureq/cert-manager-webhook-ovh/pull/8). Thanks ([@Zcool85](https://github.com/Zcool85))
- ✨ build image using Go 1.19.4

## v0.4.0

### Major features

- ⚠️ breaking changes ahead if comming from previous version
- 📄 documentation and helm chart hosted at [https://aureq.github.io/cert-manager-webhook-ovh/](https://aureq.github.io/cert-manager-webhook-ovh/)
- ✨ deploy multiple `Issuer` (namespaced) and `ClusterIssuer` via chart
- ✨ either specify your OVH credentials, or use an existing secret
- ✨ OVH credential are all stored in a secret (ApplicationKey, ApplicaitonSecret, ConsumerKey)
- ✨ deploy necessary permissions to access the OVH credentials
- ✨ role based access control to access secrets across namespaces
- 🚀 publish container image on GitHub Container Registry
- 🚀 publish Helm Chart on GitHub pages
- ⬆️ upgrade dependencies to reduce warnings
- ✨ drop root privileges
- ✨ add support for HTTP/HTTPS proxy

### Noteworthy changes

- 🚀 use kubernetes recommended labels
- ✨ move some helm logic in _helpers.tpl
- ✨ completely rework `values.yaml` to support creating issuers and ovh credentials
- ✨ create role and bind it so the webhook can access necessary secrets
- ⬆️ upgrade dependencies to reduce warnings
  - cert-manager `v1.5.3` to `v1.9.1`
  - go-ovh `v1.1.0` to `v1.3.0`
  - client-go `v0.22.1` to `v0.24.2`
- build webhook using golang `1.18`
- ✨ add image pull secrets to helm chart by Julian Stiller)
- 🐛 fix base64 encoded secrets by [@julienkosinski](https://github.com/julienkosinski)
- 🔥 drop root privilges (missing attribution)
- 🐛 fix how security context is checked
- ✨ add RBAC (missing attribution)
- ⬆️ upgrade to Alpine Linux 3.16 container image
- 🐛 fix `Makefile` references and enable HTTP proxy to local build environment
- ✨ set `CAP_NET_BIND_SERVICE` to binary to bind on privileged ports without root privileges (missing attribution)
- 🐛 add `libpcap` to container image
- ✨ create logo based on cert-manager logo and [icons8](https://icons8.com/icon/92/link)
- ✨ more fields populated in `Chart.yaml`
- 🌱 some ground work to automate the release process via GitHub Actions and GitHub packages

## v0.4.0-alpha.1

### Major features

- ⚠️ breaking changes ahead
- ✨ major helm chart improvements
- ✨ deploy multiple `Issuer` (namespaced) and `ClusterIssuer` via chart
- ✨ either specify your OVH credentials, or use an existing secret
- ✨ OVH credential are all stored in a secret (ApplicationKey, ApplicaitonSecret, ConsumerKey)
- ✨ deploy necessary permissions to access the OVH credentials
- ✨ role based access control to access secrets across namespaces

### Note worthy changes

- ✨ move some helm logic in _helpers.tpl
- ✨ completely rework `values.yaml` to support creating issuers and ovh credentials
- ✨ create role and bind it so the webhook can access necessary secrets
- 📄 documentation and helm chart hosted at [https://aureq.github.io/cert-manager-webhook-ovh/](https://aureq.github.io/cert-manager-webhook-ovh/)

## v0.4.0-alpha.0

### Major features

- 🚀 publish container image on GitHub Container Registry
- 🚀 publish Helm Chart on GitHub pages
- ⬆️ upgrade dependencies to reduce warnings
- ✨ drop root privileges
- 🌱 some ground work to automate the release process via GitHub Actions

### Noteworthy changes

- ✨ add support for HTTP proxy
- ⬆️ upgrade dependencies to reduce warnings
  - cert-manager `v1.5.3` to `v1.9.1`
  - go-ovh `v1.1.0` to `v1.3.0`
  - client-go `v0.22.1` to `v0.24.2`
- build webhook using golang `1.18`
- ✨ add image pull secrets to helm chart by Julian Stiller)
- 🐛 fix base64 encoded secrets by [@julienkosinski](https://github.com/julienkosinski)
- 🔥 drop root privilges (missing attribution)
- 🐛 fix how security context is checked
- ✨ add RBAC (missing attribution)
- ⬆️ upgrade to Alpine Linux 3.16 container image
- 🐛 fix `Makefile` references and enable HTTP proxy to local build environment
- ✨ set `CAP_NET_BIND_SERVICE` to binary to bind on privileged ports without root privileges (missing attribution)
- 🐛 add `libpcap` to container image
- ✨ create logo based on cert-manager logo and [icons8](https://icons8.com/icon/92/link)
- ✨ more fields populated in `Chart.yaml`
- 🌱 some ground work to automate the release process via GitHub Actions and GitHub packages

## 0.3.0

- Initial work by [@baarde](https://github.com/baarde)
- [cert-manager-webhook-ovh](https://github.com/baarde/cert-manager-webhook-ovh/)
- Commit [`ab4d192`](https://github.com/baarde/cert-manager-webhook-ovh/commit/ab4d192358ed7048091e1788e7256fc4fbf5e767)
