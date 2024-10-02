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
    {{- if eq "application" .authenticationMethod }}
      {{- if and .applicationConsumerKey .applicationKey .applicationSecret }}
        {{- true -}}
      {{- else }}
        {{- fail "Error: 'applicationConsumerKey', 'applicationKey', and 'applicationSecret' must all be provided for 'application' authentication method." }}
      {{- end }}
    {{- else if eq "oauth2" .authenticationMethod }}
      {{- if and .oauth2ClientId .oauth2ClientSecret }}
        {{- true -}}
      {{- else }}
        {{- fail "Error: 'oauth2ClientId' and 'oauth2ClientSecret' must both be provided for 'oauth2' authentication method." }}
      {{- end }}
    {{- else }}
      {{- fail "Error: Invalid 'authenticationMethod'. It must be either 'application' or 'oauth2'." }}
    {{- end }}
  {{- end -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationRefAvail" -}}
  {{- if . -}}
    {{- if eq "application" .authenticationMethod }}
      {{- if and (not .applicationConsumerKeyRef) (not .applicationKeyRef) (not .applicationSecretRef) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, 'applicationConsumerKeyRef', 'applicationKeyRef' and 'applicationSecretRef' need to be provided for 'application' authentication method." }}
      {{- end }}

      {{- if or (not .applicationConsumerKeyRef.name) (not .applicationConsumerKeyRef.key) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.applicationConsumerKeyRef.name' and 'ovhAuthenticationRef.applicationConsumerKeyRef.key'" }}
      {{- end }}
      {{- if or (not .applicationKeyRef.name) (not .applicationKeyRef.key) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.applicationKeyRef.name' and 'ovhAuthenticationRef.applicationKeyRef.key'" }}
      {{- end }}
      {{- if or (not .applicationSecretRef.name) (not .applicationSecretRef.key) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.applicationSecretRef.name' and 'ovhAuthenticationRef.applicationSecretRef.key'" }}
      {{- end }}

    {{- else if eq "oauth2" .authenticationMethod }}
      {{- if and (not .oauth2ClientIdRef) (not .oauth2ClientSecretRef) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, 'oauth2ClientIdRef' and 'oauth2ClientSecretRef' need to be provided for 'oauth2' authentication method." }}
      {{- end }}

      {{- if or (not .oauth2ClientIdRef.name) (not .oauth2ClientIdRef.key) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.oauth2ClientIdRef.name' and 'ovhAuthenticationRef.oauth2ClientIdRef.key'" }}
      {{- end }}
      {{- if or (not .oauth2ClientSecretRef.name) (not .oauth2ClientSecretRef.key) }}
        {{- fail "Error: When 'ovhAuthenticationRef' is used, you need to provide 'ovhAuthenticationRef.oauth2ClientSecretRef.name' and 'ovhAuthenticationRef.oauth2ClientSecretRef.key'" }}
      {{- end }}

    {{- else }}
      {{- fail "Error: Invalid 'authenticationMethod'. It must be either 'application' or 'oauth2'." }}
    {{- end }}
    {{- true -}}
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
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "cert-manager-webhook-ovh.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cert-manager-webhook-ovh.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}