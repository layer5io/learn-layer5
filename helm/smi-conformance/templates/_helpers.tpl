{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "smi-conformance.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "smi-conformance.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "smi-conformance.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "smi-conformance.labels" -}}
helm.sh/chart: {{ include "smi-conformance.chart" . }}
{{ include "smi-conformance.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "smi-conformance.selectorLabels" -}}
app.kubernetes.io/name: {{ include "smi-conformance.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "smi-conformance.serviceAccount" -}}
{{- if .Values.serviceAccount }}
{{- default (include "smi-conformance.fullname" .) .Values.serviceAccount }}
{{- else }}
{{- default "default" .Values.serviceAccount }}
{{- end }}
{{- end }}

{{/*
Create the name of the namespace to use
*/}}
{{- define "smi-conformance.namespace" -}}
{{- if .Values.namespace }}
{{- default (include "smi-conformance.fullname" .) .Values.namespace }}
{{- else }}
{{- default "default" .Values.namespace }}
{{- end }}
{{- end }}
