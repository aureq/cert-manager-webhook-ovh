GO ?= $(shell which go)
OS ?= $(shell $(GO) env GOOS)
ARCH ?= $(shell $(GO) env GOARCH)

ENVTEST_K8S_VERSION=1.35.0

IMAGE_NAME := "aureq/cert-manager-webhook-ovh"
IMAGE_TAG := "latest"

.PHONY: tests go-test setup-envtest envtest build helm-test rendered-manifest.yaml install-helm-docs schema install-helm-schema

OUT := $(shell pwd)/_out
TEST_ASSET_ETCD := $(OUT)/kubebuilder/bin/etcd
TEST_ASSET_KUBE_APISERVER := $(OUT)/kubebuilder/bin/kube-apiserver
TEST_ASSET_KUBECTL := $(OUT)/kubebuilder/bin/kubectl

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p "$(LOCALBIN)"

ENVTEST ?= $(LOCALBIN)/setup-envtest

ENVTEST_VERSION ?= $(shell v='$(call gomodver,sigs.k8s.io/controller-runtime)'; \
	[ -n "$$v" ] || { echo "Set ENVTEST_VERSION manually (controller-runtime replace has no tag)" >&2; exit 1; }; \
	printf '%s\n' "$$v" | sed -E 's/^v?([0-9]+)\.([0-9]+).*/release-\1.\2/')

ENVTEST_K8S_VERSION ?= $(shell v='$(call gomodver,k8s.io/api)'; \
	[ -n "$$v" ] || { echo "Set ENVTEST_K8S_VERSION manually (k8s.io/api replace has no tag)" >&2; exit 1; }; \
  	printf '%s\n' "$$v" | sed -E 's/^v?[0-9]+\.([0-9]+).*/1.\1/')

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink -- "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
	set -e; \
	package=$(2)@$(3) ;\
	echo "Installing $${package}" ;\
	rm -f "$(1)" ;\
	GOBIN="$(LOCALBIN)" go install $${package} ;\
	mv "$(LOCALBIN)/$$(basename "$(1)")" "$(1)-$(3)" ;\
} ;\
ln -sf "$$(realpath "$(1)-$(3)")" "$(1)"
endef

define gomodver
$(shell go list -m -f '{{if .Replace}}{{.Replace.Version}}{{else}}{{.Version}}{{end}}' $(1) 2>/dev/null)
endef

go-test: setup-envtest
	@TEST_ASSET_ETCD=$(LOCALBIN)/k8s/$(ENVTEST_K8S_VERSION)-$(OS)-$(ARCH)/etcd \
		TEST_ASSET_KUBE_APISERVER=$(LOCALBIN)/k8s/$(ENVTEST_K8S_VERSION)-$(OS)-$(ARCH)/kube-apiserver \
		TEST_ASSET_KUBECTL=$(LOCALBIN)/k8s/$(ENVTEST_K8S_VERSION)-$(OS)-$(ARCH)/kubectl \
		$(GO) test -v .


setup-envtest: envtest ## Download the binaries required for ENVTEST in the local bin directory.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@"$(ENVTEST)" use $(ENVTEST_K8S_VERSION) --bin-dir "$(LOCALBIN)" -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

clean:
	@chmod -R u+w $(LOCALBIN) $(OUT) 2>/dev/null || true
	@rm -rf $(LOCALBIN) $(OUT)

tests: go-test helm-unittest

prepare: helm-schema helm-docs

local-build:
	CGO_ENABLED=0 go build -trimpath -o cert-manager-webhook-ovh .

build:
	@test -z "$$HTTP_PROXY" -a -z "$$HTTPS_PROXY" || docker buildx build \
		--progress=plain \
		--compress \
		--output type=image,oci-mediatypes=true,compression=estargz,force-compression=true,push=false \
		--build-arg "HTTP_PROXY=$$HTTP_PROXY" \
		--build-arg "HTTPS_PROXY=$$HTTPS_PROXY" \
		-t "$(IMAGE_NAME):$(IMAGE_TAG)" .
	@test ! -z "$$HTTP_PROXY" -o ! -z "$$HTTPS_PROXY" || docker buildx build \
		--progress=plain \
		--compress \
		--output type=image,oci-mediatypes=true,compression=estargz,force-compression=true,push=false \
		-t "$(IMAGE_NAME):$(IMAGE_TAG)" .

rendered-manifest.yaml:
	@test -d "$(OUT)" || mkdir -p "$(OUT)"
	@helm template \
	    cert-manager-webhook-ovh \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=$(IMAGE_TAG) \
        charts/cert-manager-webhook-ovh > "$(OUT)/rendered-manifest.yaml"

helm-unittest: install-helm-unittest
	@helm unittest charts/cert-manager-webhook-ovh/

helm-docs: install-helm-docs
	@helm-docs --chart-search-root=charts/ --template-files=./_templates.gotmpl --template-files=README.md.gotmpl --sort-values-order=file

helm-schema: install-helm-schema
	@helm-schema --chart-search-root charts/cert-manager-webhook-ovh/ --add-schema-reference --keep-full-comment

install-helm-unittest:
	@helm plugin list | grep ^unittest >/dev/null 2>&1 || helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false

install-helm-docs:
	@go install github.com/norwoodj/helm-docs/cmd/helm-docs@latest

install-helm-schema:
	@go install github.com/dadav/helm-schema/cmd/helm-schema@latest
