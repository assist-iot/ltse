{{/*
Expand the name of the chart.
*/}}
{{- define "enabler.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "enabler.fullname" -}}
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
{{- define "enabler.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Name of the component api.
*/}}
{{- define "api.name" -}}
{{ include "enabler.name" . }}-api
{{- end }}

{{/*
Component api labels
*/}}
{{- define "api.labels" -}}
helm.sh/chart: {{ include "enabler.chart" . }}
{{ include "api.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Component api selector labels
*/}}
{{- define "api.selectorLabels" -}}
app.kubernetes.io/name: {{ include "enabler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
enabler: {{ .Chart.Name }}
component: {{ .Values.api.name }}
isMainInterface: "yes"
tier: {{ .Values.api.tier}}
{{- end }}

{{/*
Name of the component postgrest.
*/}}
{{- define "postgrest.name" -}}
{{ include "enabler.name" . }}-postgrest
{{- end }}

{{/*
Component postgrest labels
*/}}
{{- define "postgrest.labels" -}}
helm.sh/chart: {{ include "enabler.chart" . }}
{{ include "postgrest.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Component postgrest selector labels
*/}}
{{- define "postgrest.selectorLabels" -}}
app.kubernetes.io/name: {{ include "enabler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
enabler: {{ .Chart.Name }}
component: {{ .Values.postgrest.name }}
isMainInterface: "no"
tier: {{ .Values.postgrest.tier}}
{{- end }}

{{/*
Name of the component postgres.
*/}}
{{- define "postgres.name" -}}
{{ include "enabler.name" . }}-postgres
{{- end }}

{{/*
Component postgres labels
*/}}
{{- define "postgres.labels" -}}
helm.sh/chart: {{ include "enabler.chart" . }}
{{ include "postgres.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Component postgres selector labels
*/}}
{{- define "postgres.selectorLabels" -}}
app.kubernetes.io/name: {{ include "enabler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
enabler: {{ .Chart.Name }}
component: {{ .Values.postgres.name }}
isMainInterface: "no"
tier: {{ .Values.postgres.tier}}
{{- end }}

{{/*
Name of the component elastic.
*/}}
{{- define "elastic.name" -}}
{{ include "enabler.name" . }}-elastic
{{- end }}

{{/*
Component elastic labels
*/}}
{{- define "elastic.labels" -}}
helm.sh/chart: {{ include "enabler.chart" . }}
{{ include "elastic.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Component elastic selector labels
*/}}
{{- define "elastic.selectorLabels" -}}
app.kubernetes.io/name: {{ include "enabler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
enabler: {{ .Chart.Name }}
component: {{ .Values.elastic.name }}
isMainInterface: "no"
tier: {{ .Values.elastic.tier}}
{{- end }}

