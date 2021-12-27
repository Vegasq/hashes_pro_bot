## Build

```bash
docker buildx build . -t hpbot
```

## Run

```bash
docker run --restart unless-stopped -d -t -i -e TELEGRAM_TOKEN=<SECRET TELEGRAM TOKEN> hpbot
```
