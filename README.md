# crawljob-api

A lightweight REST API written in Go that generates `.crawljob` files for [JDownloader](https://jdownloader.org/). Drop a download URL, get a crawljob file picked up automatically by JDownloader.

Built to run as a Docker container.

---

## How it works

1. You send a `POST /jobs` request with a download URL
2. The API validates the URL (scheme, allowed domains)
3. A `.crawljob` file is generated and dropped into a watched folder
4. JDownloader picks it up and starts the download automatically

---

## Getting Started

### Run with Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e DESTINATION_FOLDER=/mnt/downloads \
  -e CRAWLJOB_FOLDER=/mnt/crawljobs \
  -v /your/download/path:/mnt/downloads \
  -v /your/crawljob/path:/mnt/crawljobs \
  frostbyte0x/crawljob-api:latest
```

### Build locally

```bash
git clone https://github.com/FrostByte0x/crawljob-api
cd crawljob-api
go run main.go
```

---

## Configuration

| Variable | Description | Default |
|---|---|---|
| `DESTINATION_FOLDER` | Download destination folder | `.` (current dir) |
| `CRAWLJOB_FOLDER` | Folder watched by JDownloader | `.` (current dir) |

---

## API Reference

### `POST /jobs`

Submit a download URL.

**Request**
```json
{
  "url": "https://1fichier.com/yourfile"
}
```

**Responses**

| Code | Description |
|---|---|
| `201 Created` | Job file successfully created |
| `400 Bad Request` | Invalid URL or malformed body |
| `405 Method Not Allowed` | Only POST is accepted |

---

## Allowed Domains
This can be changed in the Dockerfile configuration 

Currently restricted to:
- `1fichier.com`
- `mega.nz`

> Contact the server owner or set your own domain list to extend this.

---

## Project Structure

```
crawljob-api/
├── main.go           # Server entrypoint
├── handler/
│   ├── job.go        # HTTP handler
│   └── validator.go  # URL validation
├── model/
│   ├── crawljob.go   # CrawlJob model + file generation
│   └── utils.go      # Helpers
└── Dockerfile
```

---

## License

MIT
