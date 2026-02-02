IMAGE_NAME := "aureq/cert-manager-webhook-ovh"
IMAGE_TAG := "latest"

.PHONY: test build helm-test rendered-manifest.yaml install-helm-docs schema install-helm-schema

OUT := $(shell pwd)/_out
TEST_ASSET_ETCD := $(OUT)/kubebuilder/bin/etcd
TEST_ASSET_KUBE_APISERVER := $(OUT)/kubebuilder/bin/kube-apiserver
TEST_ASSET_KUBECTL := $(OUT)/kubebuilder/bin/kubectl

test:
	@test -d "$(OUT)" || mkdir -p "$(OUT)"
	@bash ./scripts/fetch-test-binaries.sh
	TEST_ASSET_ETCD="$(TEST_ASSET_ETCD)" \
		TEST_ASSET_KUBE_APISERVER="$(TEST_ASSET_KUBE_APISERVER)" \
		TEST_ASSET_KUBECTL="$(TEST_ASSET_KUBECTL)" \
		go test -v .

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
	@helm-docs --chart-search-root charts/cert-manager-webhook-ovh/ --sort-values-order=file

helm-schema: install-helm-schema
	@helm-schema --chart-search-root charts/cert-manager-webhook-ovh/ --add-schema-reference --keep-full-comment

install-helm-unittest:
	@helm plugin install https://github.com/helm-unittest/helm-unittest

install-helm-docs:
	@go install github.com/norwoodj/helm-docs/cmd/helm-docs@latest

install-helm-schema:
	@go install github.com/dadav/helm-schema/cmd/helm-schema@latest
