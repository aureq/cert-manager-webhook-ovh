{{/* vim: set filetype=mustache: */}}

{{- /* 🔥 function should NOT `fail` */ -}}

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
  {{- printf "%s-selfsign" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.rootCAIssuer" -}}
  {{- printf "%s-ca" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.rootCACertificate" -}}
  {{- printf "%s-ca" (include "cert-manager-webhook-ovh.fullname" .) }}
{{- end -}}

{{- define "cert-manager-webhook-ovh.servingCertificate" -}}
  {{- printf "%s-webhook-tls" (include "cert-manager-webhook-ovh.fullname" .) }}
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

{{/* --------------------------------------------------------------------- */}}
{{/* ----------------------- Direct Authentication  ---------------------- */}}
{{/* --------------------------------------------------------------------- */}}

{{/*
Returns true if ovhAuthentication.application* is correctly set.
No checks done on ovhAuthenticationMethod, this function only checks if the required fields for application authentication are set in ovhAuthentication.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh._applicationOvhAuthenticationAvail" -}}
  {{- $result := "false" -}}
  {{- if .ovhAuthentication -}}
    {{- if and .ovhAuthentication.applicationKey .ovhAuthentication.applicationSecret .ovhAuthentication.applicationConsumerKey -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthentication is correctly set when ovhAuthenticationMethod is "application".
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.applicationOvhAuthenticationAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "application" -}}
    {{- if eq (include "cert-manager-webhook-ovh._applicationOvhAuthenticationAvail" .) "true" -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}{{/* end if ovhAuthenticationMethod */}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthentication.oauth2* is correctly set.
No checks done on ovhAuthenticationMethod, this function only checks if the required fields for oauth2 authentication are set in ovhAuthentication.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh._oauth2OvhAuthenticationAvail" -}}
  {{- $result := "false" -}}
  {{- if .ovhAuthentication -}}
    {{- if and .ovhAuthentication.oauth2ClientID .ovhAuthentication.oauth2ClientSecret -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}


{{/*
Returns true if ovhAuthentication is correctly set when ovhAuthenticationMethod is "oauth2".
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.oauth2OvhAuthenticationAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "oauth2" -}}
    {{- if eq (include "cert-manager-webhook-ovh._oauth2OvhAuthenticationAvail" .) "true" -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}{{/* end if ovhAuthenticationMethod */}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthentication is correctly set depending on the ovhAuthenticationMethod.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "application" -}}
    {{- if and (eq (include "cert-manager-webhook-ovh.applicationOvhAuthenticationAvail" .) "true") (not (eq (include "cert-manager-webhook-ovh.oauth2OvhAuthenticationAvail" .) "true")) }}
      {{- $result = "true" -}}
    {{- end -}}
  {{- else if eq .ovhAuthenticationMethod "oauth2" -}}
    {{- if and (eq (include "cert-manager-webhook-ovh.oauth2OvhAuthenticationAvail" .) "true") (not (eq (include "cert-manager-webhook-ovh.applicationOvhAuthenticationAvail" .) "true")) }}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}

{{/* --------------------------------------------------------------------- */}}
{{/* -------------------- Authentication by reference -------------------- */}}
{{/* --------------------------------------------------------------------- */}}

{{/*
Returns true if ovhAuthenticationRef.application* is correctly set.
No checks done on ovhAuthenticationMethod, this function only checks if the required fields for application authentication are set in ovhAuthenticationRef.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh._applicationOvhAuthenticationRefAvail" -}}
  {{- $result := "false" -}}
  {{- if .ovhAuthenticationRef -}}
    {{- if and .ovhAuthenticationRef.applicationKeyRef .ovhAuthenticationRef.applicationSecretRef .ovhAuthenticationRef.applicationConsumerKeyRef -}}
      {{- if and .ovhAuthenticationRef.applicationKeyRef.name .ovhAuthenticationRef.applicationKeyRef.key .ovhAuthenticationRef.applicationSecretRef.name .ovhAuthenticationRef.applicationSecretRef.key .ovhAuthenticationRef.applicationConsumerKeyRef.name .ovhAuthenticationRef.applicationConsumerKeyRef.key -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set when ovhAuthenticationMethod is "application".
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.applicationOvhAuthenticationRefAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "application" -}}
    {{- if eq (include "cert-manager-webhook-ovh._applicationOvhAuthenticationRefAvail" .) "true" -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}{{/* end if ovhAuthenticationMethod */}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef.oauth2* is correctly set.
No checks done on ovhAuthenticationMethod, this function only checks if the required fields for oauth2 authentication are set in ovhAuthenticationRef.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh._oauth2OvhAuthenticationRefAvail" -}}
  {{- $result := "false" -}}
  {{- if .ovhAuthenticationRef -}}
    {{- if and .ovhAuthenticationRef.oauth2ClientIDRef .ovhAuthenticationRef.oauth2ClientSecretRef -}}
      {{- if and .ovhAuthenticationRef.oauth2ClientIDRef.name .ovhAuthenticationRef.oauth2ClientIDRef.key .ovhAuthenticationRef.oauth2ClientSecretRef.name .ovhAuthenticationRef.oauth2ClientSecretRef.key -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set when ovhAuthenticationMethod is "oauth2".
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.oauth2OvhAuthenticationRefAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "oauth2" -}}
    {{- if eq (include "cert-manager-webhook-ovh._oauth2OvhAuthenticationRefAvail" .) "true" -}}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}{{/* end if ovhAuthenticationMethod */}}
  {{- eq $result "true" -}}
{{- end -}}

{{/*
Returns true if ovhAuthenticationRef is correctly set depending on the ovhAuthenticationMethod.
$arg1: The issuer values
*/}}
{{- define "cert-manager-webhook-ovh.isOvhAuthenticationRefAvail" -}}
  {{- $result := "false" -}}
  {{- if eq .ovhAuthenticationMethod "application" -}}
    {{- if and (eq (include "cert-manager-webhook-ovh.applicationOvhAuthenticationRefAvail" .) "true") (not (eq (include "cert-manager-webhook-ovh.oauth2OvhAuthenticationRefAvail" .) "true")) }}
      {{- $result = "true" -}}
    {{- end -}}
  {{- else if eq .ovhAuthenticationMethod "oauth2" -}}
    {{- if and (eq (include "cert-manager-webhook-ovh.oauth2OvhAuthenticationRefAvail" .) "true") (not (eq (include "cert-manager-webhook-ovh.applicationOvhAuthenticationRefAvail" .) "true")) }}
      {{- $result = "true" -}}
    {{- end -}}
  {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}

{{/* --------------------------------------------------------------------- */}}
{{/* ----------------- Single authentication per issuer ------------------ */}}
{{/* --------------------------------------------------------------------- */}}

{{/* Returns true if only one authentication method is set, regardless of the method. */}}
{{- define "cert-manager-webhook-ovh.isMultipleAuthenticationMethodSet" -}}
  {{- $result := "false" -}}
    {{/* Check if direct authentication is set for application */}}
    {{- if and .ovhAuthentication .ovhAuthentication.applicationKey .ovhAuthentication.applicationSecret .ovhAuthentication.applicationConsumerKey -}}
      {{- if or
        (and .ovhAuthentication.oauth2ClientID .ovhAuthentication.oauth2ClientSecret)
        (and .ovhAuthenticationRef .ovhAuthenticationRef.applicationKeyRef .ovhAuthenticationRef.applicationKeyRef.name .ovhAuthenticationRef.applicationKeyRef.key .ovhAuthenticationRef.applicationSecretRef .ovhAuthenticationRef.applicationSecretRef.name .ovhAuthenticationRef.applicationSecretRef.key .ovhAuthenticationRef.applicationConsumerKeyRef .ovhAuthenticationRef.applicationConsumerKeyRef.name .ovhAuthenticationRef.applicationConsumerKeyRef.key)
        (and .ovhAuthenticationRef .ovhAuthenticationRef.oauth2ClientIDRef .ovhAuthenticationRef.oauth2ClientIDRef.name .ovhAuthenticationRef.oauth2ClientIDRef.key .ovhAuthenticationRef.oauth2ClientSecretRef .ovhAuthenticationRef.oauth2ClientSecretRef.name .ovhAuthenticationRef.oauth2ClientSecretRef.key) -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{/* Check if direct authentication is set for oauth2 */}}
    {{- else if and .ovhAuthentication .ovhAuthentication.oauth2ClientID .ovhAuthentication.oauth2ClientSecret -}}
      {{- if or
        (and .ovhAuthentication.applicationKey .ovhAuthentication.applicationSecret .ovhAuthentication.applicationConsumerKey)
        (and .ovhAuthenticationRef .ovhAuthenticationRef.applicationKeyRef .ovhAuthenticationRef.applicationKeyRef.name .ovhAuthenticationRef.applicationKeyRef.key .ovhAuthenticationRef.applicationSecretRef .ovhAuthenticationRef.applicationSecretRef.name .ovhAuthenticationRef.applicationSecretRef.key .ovhAuthenticationRef.applicationConsumerKeyRef .ovhAuthenticationRef.applicationConsumerKeyRef.name .ovhAuthenticationRef.applicationConsumerKeyRef.key)
        (and .ovhAuthenticationRef .ovhAuthenticationRef.oauth2ClientIDRef .ovhAuthenticationRef.oauth2ClientIDRef.name .ovhAuthenticationRef.oauth2ClientIDRef.key .ovhAuthenticationRef.oauth2ClientSecretRef .ovhAuthenticationRef.oauth2ClientSecretRef.name .ovhAuthenticationRef.oauth2ClientSecretRef.key) -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{/* Check if authentication by reference is set for application */}}
    {{- else if and .ovhAuthenticationRef .ovhAuthenticationRef.applicationKeyRef .ovhAuthenticationRef.applicationKeyRef.name .ovhAuthenticationRef.applicationKeyRef.key .ovhAuthenticationRef.applicationSecretRef .ovhAuthenticationRef.applicationSecretRef.name .ovhAuthenticationRef.applicationSecretRef.key .ovhAuthenticationRef.applicationConsumerKeyRef .ovhAuthenticationRef.applicationConsumerKeyRef.name .ovhAuthenticationRef.applicationConsumerKeyRef.key -}}
      {{- if or
        (and .ovhAuthentication .ovhAuthentication.applicationKey .ovhAuthentication.applicationSecret .ovhAuthentication.applicationConsumerKey)
        (and .ovhAuthentication .ovhAuthentication.oauth2ClientID .ovhAuthentication.oauth2ClientSecret)
        (and .ovhAuthenticationRef.oauth2ClientIDRef .ovhAuthenticationRef.oauth2ClientIDRef.name .ovhAuthenticationRef.oauth2ClientIDRef.key .ovhAuthenticationRef.oauth2ClientSecretRef .ovhAuthenticationRef.oauth2ClientSecretRef.name .ovhAuthenticationRef.oauth2ClientSecretRef.key) -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{/* Check if authentication by reference is set for oauth2 */}}
    {{- else if and .ovhAuthenticationRef .ovhAuthenticationRef.oauth2ClientIDRef .ovhAuthenticationRef.oauth2ClientIDRef.name .ovhAuthenticationRef.oauth2ClientIDRef.key .ovhAuthenticationRef.oauth2ClientSecretRef .ovhAuthenticationRef.oauth2ClientSecretRef.name .ovhAuthenticationRef.oauth2ClientSecretRef.key -}}
      {{- if or
        (and .ovhAuthentication .ovhAuthentication.applicationKey .ovhAuthentication.applicationSecret .ovhAuthentication.applicationConsumerKey)
        (and .ovhAuthentication .ovhAuthentication.oauth2ClientID .ovhAuthentication.oauth2ClientSecret)
        (and .ovhAuthenticationRef.applicationKeyRef .ovhAuthenticationRef.applicationKeyRef.name .ovhAuthenticationRef.applicationKeyRef.key .ovhAuthenticationRef.applicationSecretRef .ovhAuthenticationRef.applicationSecretRef.name .ovhAuthenticationRef.applicationSecretRef.key .ovhAuthenticationRef.applicationConsumerKeyRef .ovhAuthenticationRef.applicationConsumerKeyRef.name .ovhAuthenticationRef.applicationConsumerKeyRef.key) -}}
        {{- $result = "true" -}}
      {{- end -}}
    {{- end -}}
  {{- eq $result "true" -}}
{{- end -}}