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
Returns true if ovhAuthentication is correctly set when ovhAuthenticationMethod is "application".
*/}}
{{- define "cert-manager-webhook-ovh.applicationOvhAuthenticationAvail" -}}
  {{- if . -}}
    {{- with .ovhAuthentication -}}
      {{- if and (.applicationConsumerKey) (.applicationKey) (.applicationSecret) -}}
        {{- eq "true" "true" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Returns true if ovhAuthentication is correctly set when ovhAuthenticationMethod is "oauth2".
*/}}
{{- define "cert-manager-webhook-ovh.oauth2OvhAuthenticationAvail" -}}
  {{- if . -}}
    {{- with .ovhAuthentication -}}
      {{- if and (.oauth2ClientID) (.oauth2ClientSecret) -}}
        {{- eq "true" "true" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Returns true if ovhAuthentication is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationAvail" -}}
  {{- if . -}}
    {{- if eq .ovhAuthenticationMethod "application" -}}
      {{ include "cert-manager-webhook-ovh.applicationOvhAuthenticationAvail" . }}
    {{- else if eq .ovhAuthenticationMethod "oauth2" -}}
      {{ include "cert-manager-webhook-ovh.oauth2OvhAuthenticationAvail" . }}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set when ovhAuthenticationMethod is "application".
*/}}
{{- define "cert-manager-webhook-ovh.applicationOvhAuthenticationRefAvail" -}}
  {{- if .ovhAuthenticationRef -}}
    {{- if eq .ovhAuthenticationMethod "application" -}}
      {{- with .ovhAuthenticationRef -}}
        {{- if or (not .applicationConsumerKeyRef) (not .applicationKeyRef) (not .applicationSecretRef) }}
          {{- fail "Error: When 'ovhAuthenticationRef' is used, 'applicationConsumerKeyRef', 'applicationKeyRef' and 'applicationSecretRef' need to be provided." }}
        {{- end }}
        {{- if or (not .applicationConsumerKeyRef.name) (not .applicationConsumerKeyRef.key) }}
          {{- eq "true" "false" -}}
        {{- end }}
        {{- if or (not .applicationKeyRef.name) (not .applicationKeyRef.key) }}
          {{- eq "true" "false" -}}
        {{- end }}
        {{- if or (not .applicationSecretRef.name) (not .applicationSecretRef.key) }}
          {{- eq "true" "false" -}}
        {{- end }}
      {{- end -}}
    {{- end -}}{{/* end if ovhAuthenticationMethod */}}
    {{- eq "true" "true" -}}
  {{- end -}}{{/* end if ovhAuthenticationRef */}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set when ovhAuthenticationMethod is "oauth2".
*/}}
{{- define "cert-manager-webhook-ovh.oauth2OvhAuthenticationRefAvail" -}}
  {{- if .ovhAuthenticationRef -}}
    {{- if eq .ovhAuthenticationMethod "oauth2" -}}
      {{- with .ovhAuthenticationRef -}}
        {{- if or (not .oauth2ClientIDRef) (not .oauth2ClientSecretRef) }}
          {{- eq "true" "false" -}}
        {{- end }}
        {{- if or (not .oauth2ClientIDRef.name) (not .oauth2ClientIDRef.key) }}
          {{- eq "true" "false" -}}
        {{- end }}
        {{- if or (not .oauth2ClientSecretRef.name) (not .oauth2ClientSecretRef.key) }}
          {{- eq "true" "false" -}}
        {{- end }}
      {{- end -}}
    {{- end -}}{{/* end if ovhAuthenticationMethod */}}
    {{- eq "true" "true" -}}
  {{- end -}}{{/* end if ovhAuthenticationRef */}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationRefAvail" -}}
  {{- if .ovhAuthenticationRef -}}
    {{- if eq .ovhAuthenticationMethod "application" -}}
      {{ include "cert-manager-webhook-ovh.applicationOvhAuthenticationRefAvail" . }}
    {{- else if eq .ovhAuthenticationMethod "oauth2" -}}
      {{ include "cert-manager-webhook-ovh.oauth2OvhAuthenticationRefAvail" . }}
    {{- end -}}{{/* end if ovhAuthenticationMethod */}}
    {{- eq "true" "true" -}}
  {{- end -}}{{/* end if ovhAuthenticationRef */}}
{{- end -}}

{{/*
Returns true if externalAccountBinding is correctly set.
*/}}
{{- define "cert-manager-webhook-ovh.isExternalAccountBindingAvail" -}}
  {{- if .externalAccountBinding.enabled -}}
    {{- if and (.externalAccountBinding.keyID) (.externalAccountBinding.keySecretRef.name) (.externalAccountBinding.keySecretRef.key) -}}
      {{- eq "true" "true" -}}
    {{- end -}}
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

{{/*
Return the proper image name
*/}}
{{- define "cert-manager-webhook-ovh.image" -}}
    {{- $registryName := default "ghcr.io" .Values.image.registry -}}
    {{- $repositoryName := default "authelia/authelia" .Values.image.repository -}}
    {{- $tag := .Values.image.tag | default .Chart.AppVersion  | toString -}}
    {{- if hasPrefix "sha256:" $tag }}
    {{- printf "%s/%s@%s" $registryName $repositoryName $tag -}}
    {{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- end -}}