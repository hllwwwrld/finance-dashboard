.DEFAULT_GOAL := help

.PHONY: help

help:
	@echo "Локальная разработка (БД в Docker, Go и Next на машине):"
	@echo "  make dev-db            Postgres + миграции"
	@echo "  make dev-backend       API :3000"
	@echo "  make dev-frontend      UI :3001 → http://localhost:3000"
	@echo "  make dev-frontend-lan  UI в Wi‑Fi (LAN_IP из en0 или задайте вручную)"
	@echo "  make dev-all           dev-db + backend + frontend в одном терминале"
	@echo "  make lint / format / deps / bin-deps"
	@echo ""
	@echo "Сервер / Docker (HTTP и HTTPS — см. make -f Makefile.prod help):"
	@echo "  make -f Makefile.prod help"
	@echo "Команда для получения текущего динамического ip машины ipconfig getifaddr en0"

%:
	@$(MAKE) -f Makefile.local $@
