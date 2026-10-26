# WDI — serwer WWW z rozwiązaniami 8 zadań (Go, stdlib, jeden binarny plik)
#
# Budowa:  docker build -t wdi:latest .
# Uruchomienie (compose):  docker compose up --build -d   → http://localhost:8080
#
# Bazy: golang:1.26 (seria 1.25 EOL 2026-08-19), alpine:3.23 (3.20 EOL 2026-04-01).
# Runtime bez apk add: serwer nie wykonuje żądań TLS (brak http.Client w kodzie),
# a healthcheck używa wget z busybox — obraz ~15 MB i etap runtime buduje się bez sieci.

# ---------- etap 1: kompilacja ----------
FROM golang:1.26 AS build
WORKDIR /src

# Najpierw same metadane modułu (lepsze cache'owanie warstwy).
# Projekt ma zero zależności zewnętrznych, więc COPY go.sum jest zbędne.
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Statyczny binarny plik (CGO wyłączone, bez debug info) — wchodzi do etapu 2.
# touch -d @0: deterministyczny timestamp pliku (powtarzalność warstwy).
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" \
      -o /out/wdi-server ./cmd/wdi-server \
 && touch -d @0 /out/wdi-server

# ---------- etap 2: runtime ----------
FROM alpine:3.23

# Użytkownik nie-root; wget do healthchecku daje busybox.
RUN adduser -D -u 10001 wdi

USER wdi
COPY --link --chown=0:0 --from=build /out/wdi-server /usr/local/bin/wdi-server

# Wszystkie zasoby (frontend, "Pan Tadeusz", słownik) są wbudowane w binarny plik
# przez go:embed — kontener nie potrzebuje żadnego pliku poza samym binarkiem.
EXPOSE 8080
ENV PORT=8080
# Zdrowie obserwowalne także bez compose ($PORT rozwija shell w kontenerze).
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
  CMD wget -qO- "http://127.0.0.1:${PORT}/api/tasks" > /dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/wdi-server"]
