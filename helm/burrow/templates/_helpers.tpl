{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "hsc.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "hsc.fullname" -}}
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
{{- define "hsc.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Formulate the how the seeds feed is populated.
*/}}
{{- define "hsc.seeds" -}}
{{- if (and .Values.peer.ingress.enabled (not (eq (len .Values.peer.ingress.hosts) 0))) -}}
{{- $host := index .Values.peer.ingress.hosts 0 -}}
{{- range $index, $val := $.Values.validators -}}
{{- $addr := $val.nodeAddress | lower -}}
{{- $node := printf "%03d" $index -}}
tcp://{{ $addr }}@{{ $node }}.{{ $host }}:{{ $.Values.config.Tendermint.ListenPort }},
{{- end -}}
{{- if not (eq (len .Values.chain.extraSeeds) 0) -}}
{{- range .Values.chain.extraSeeds -}},{{ . }}{{- end -}}
{{- end -}}
{{- else -}}
{{- range $index, $val := $.Values.validators -}}
{{- $addr := $val.nodeAddress | lower -}}
{{- $node := printf "%03d" $index -}}
tcp://{{ $addr }}@{{ template "hsc.fullname" $ }}-peer-{{ $node }}:{{ $.Values.config.Tendermint.ListenPort }},
{{- end -}}
{{- if not (eq (len .Values.chain.extraSeeds) 0) -}}
{{- range .Values.chain.extraSeeds -}},{{ . }}{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "hsc.image" -}}
{{ printf "%s:%s" .Values.image.repository .Values.image.tag }}
{{- end -}}

{{- define "contracts.image" -}}
{{ printf "%s:%s" .Values.contracts.image.repository .Values.contracts.image.tag }}
{{- end -}}

{{- define "restore.image" -}}
{{ printf "%s:%s" .Values.restore.image.repository .Values.restore.image.tag }}
{{- end -}}