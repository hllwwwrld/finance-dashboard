# Общие переменные для Makefile.local и Makefile.prod
DC      ?= docker compose
DOCKER  ?= docker
LAN_IP  ?= $(shell ipconfig getifaddr en0 2>/dev/null || hostname -I 2>/dev/null | awk '{print $$1}')
