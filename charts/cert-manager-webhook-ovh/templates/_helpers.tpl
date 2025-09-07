{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cert-manager-webhook-ovh.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cert-manager-webhook-ovh.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cert-manager-webhook-ovh.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cert-manager-webhook-ovh.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.rootCAIssuer" -}}
{{ printf "%s-ca" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.rootCACertificate" -}}
{{ printf "%s-ca" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.servingCertificate" -}}
{{ printf "%s-webhook-tls" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{/*
Returns true if ovhAuthentication is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationAvail" -}}
  {{- if . -}}
    {{- if and (.consumerKey) (.applicationKey) (.applicationSecret) -}}
      {{- eq "true" "true" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationRefAvail" -}}
  {{- if . -}}
    {{- if or (not .consumerKeyRef) (not .applicationKeyRef) (not .applicationSecretRef) }}
      {{- fail "Error: When 'ovhAuthenticationRef' is used, 'consumerKeyRef', 'applicationKeyRef' and 'applicationSecretRef' need to be provided." }}
    {{- end }}
    {{- if or (not .consumerKeyRef.name) (not .consumerKeyRef.key) }}
      {{ fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.consumerKeyRef.name' and 'ovhAuthenticationRef.consumerKeyRef.key'" }}
    {{- end }}
    {{- if or (not .applicationKeyRef.name) (not .applicationKeyRef.key) }}
      {{ fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.applicationKeyRef.name' and 'ovhAuthenticationRef.applicationKeyRef.key'" }}
    {{- end }}
    {{- if or (not .applicationSecretRef.name) (not .applicationSecretRef.key) }}
      {{ fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.applicationSecretRef.name' and 'ovhAuthenticationRef.applicationSecretRef.key'" }}
    {{- end }}
    {{- eq "true" "true" -}}
  {{- end -}}
{{- end -}}

{{/*
Returns the number of Issuer/ClusterIssuer to create
*/}}
{{- define "cert-manager-webhook-ovh.isIssuerToCreate" -}}
  {{- $issuerCount := 0 }}
  {{- range $.Values.issuers }}
    {{- if .create }}
      {{- $issuerCount = $issuerCount | add1 -}}
    {{- end }}{{/* end if .create */}}
  {{- end }}{{/* end range */}}
  {{- $issuerCount }}
{{- end }}{{/* end define */}}

{{/*
Common/recommended labels: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
*/}}
{{- define "cert-manager-webhook-ovh.labels" -}}
helm.sh/chart: {{ include "cert-manager-webhook-ovh.chart" . }}
app.kubernetes.io/component: webhook
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: cert-manager
{{ include "cert-manager-webhook-ovh.selectorLabels" . }}
{{- if or .Chart.AppVersion .Values.image.tag }}
app.kubernetes.io/version: {{ (.Values.image.tag | default .Chart.AppVersion | split "@")._0 | quote }}
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "cert-manager-webhook-ovh.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cert-manager-webhook-ovh.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
