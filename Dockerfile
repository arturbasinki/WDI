# WDI — serwer WWW z rozwiązaniami 8 zadań (Go, stdlib, jeden binarny plik)
#
# Budowa:  docker build -t wdi:latest .
# Uruchomienie (compose):  docker compose up --build -d   → http://localhost:8080

# ---------- etap 1: kompilacja ----------
FROM golang:1.25.6 AS build
WORKDIR /src

# Najpierw same metadane modułu (lepsze cache'owanie warstwy).
# Projekt ma zero zależności zewnętrznych, więc COPY go.sum jest zbędne.
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Statyczny binarny plik (CGO wyłączone, bez debug info) — wchodzi do etapu 2.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" \
      -o /out/wdi-server ./cmd/wdi-server

# ---------- etap 2: runtime ----------
FROM alpine:3.20

# busybox wget (healthcheck) + certyfikaty systemowe + użytkownik nie-root
RUN apk add --no-cache ca-certificates wget \
 && adduser -D -u 10001 wdi

USER wdi
COPY --from=build /out/wdi-server /usr/local/bin/wdi-server

# Wszystkie zasoby (frontend, "Pan Tadeusz", słownik) są wbudowane w binarny plik
# przez go:embed — kontener nie potrzebuje żaden pliku poza samym binarkiem.
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/local/bin/wdi-server"]
