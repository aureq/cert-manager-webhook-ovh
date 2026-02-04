# OVH Webhook for Cert Manager

![OVH Webhook for Cert Manager](assets/images/cert-manager-webhook-ovh.svg "OVH Webhook for Cert Manager")

This is a webhook solver for [OVH](http://www.ovh.com) DNS. In short, if your domain has its DNS servers hosted with OVH, you can solve DNS challenges using Cert Manager and OVH Webhook for Cert Manager.

Please star this repository to help others find it.

## Features

- Solve DNS01 challenges using OVH DNS servers.
- Support both multiple Cert Manager `ClusterIssuer` and `Issuer`.
- Support both application based and OAuth2 based authentication.
- Store OVH credentials in a secret per issuer, or use secret references.
- Helm chart repository for ease and simplicity.
- Role based access control, across namespace.
- Support for optional External Account Binding (EAB).
- JSON schema for `values.yaml` automatic validation.
- Full [documentation](tree/main/charts/cert-manager-webhook-ovh) of `values.yaml`.

## Documentation

The documentation is available at [https://aureq.github.io/cert-manager-webhook-ovh/](https://aureq.github.io/cert-manager-webhook-ovh/)

## Artifact details

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/cert-manager-webhook-ovh)](https://artifacthub.io/packages/search?repo=cert-manager-webhook-ovh)

## Support

We would like to inform our users that this repository is maintained by a volunteer, and as such, it is a best effort support, and not a commercial service.
While we strive to provide the best possible support to our users, we cannot guarantee immediate or comprehensive responses to all queries.
Therefore, we advise that users seeking professional support should contact the author directly.

We also want to emphasize that any form of abuse or entitlement towards our volunteer maintainers will not be tolerated.
Our volunteers work hard to provide support to the community, and we expect all users to treat them with respect and appreciation.
We appreciate your understanding and cooperation in maintaining a positive and productive environment for everyone involved in this community-driven project.

## Release workflow

- Run `make helm-unittest`
- Update `charts/cert-manager-webhook-ovh/Chart.yaml`
- Run `make helm-docs helm-schema`
- Commit charts/cert-manager-webhook-ovh/values.schema.json and charts/cert-manager-webhook-ovh/README.md changes
- Prepare `CHANGELOG.md` for `x.y.z`
- Commit all changes
- Push all commits
- run `bash .github/scripts/changelog.sh x.y.z`

## Maintainers

- [@aureq](https://github.com/aureq)

## Contributors

- [@munnerz](https://github.com/munnerz)
- [@Diaphteiros](https://github.com/Diaphteiros)
- [@baarde](https://github.com/baarde)
- Xaver Baun
- lcavajani
- Ricardo Pchevuzinske Katz
- [@MattiasGees](https://github.com/MattiasGees)
- Jean-Marc Andre
- [@IDerr](https://github.com/IDerr)
- Robin KERDILES
- Julian Stiller
- [@julienkosinski](https://github.com/julienkosinski)
- [@aegaeonit](https://github.com/aegaeonit)
- [@TartanLeGrand](https://github.com/TartanLeGrand)
- [@Zcool85](https://github.com/Zcool85)
- [@Yethal](https://github.com/Yethal)
- [Benjamin Maisonnas](https://github.com/Benzhaomin)
- [Kebree](https://github.com/Kebree)
- [Alissia01](https://github.com/Alissia01)
- [Mathieu Sensei](https://github.com/hyu9a)
- [Sebastien MALOT](https://github.com/smalot)
- [Rémy Jacquin](https://github.com/remyj38)
- [flodakto](https://github.com/flodakto)
- [Aurelie Vache](https://github.com/scraly)
- [Karol Stoiński](https://github.com/KarolStoinski)
- [Pierre Mahot](https://github.com/pierremahot)
- [Thomas Coudert](https://github.com/thcdrt)
