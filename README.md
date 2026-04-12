[![Build and Push Docker image](https://github.com/FrostByte0x/crawljob-api/actions/workflows/docker.yml/badge.svg?event=push)](https://github.com/FrostByte0x/crawljob-api/actions/workflows/docker.yml)

# crawljob-api

A lightweight REST API written in Go that generates `.crawljob` files for [JDownloader](https://jdownloader.org/). Drop a download URL, get a crawljob file picked up automatically by JDownloader.

Built to run as a Docker container.

--- 
## Web Interface
A web interface is available at `/` and `/downloads`. It offers a simple text field and one-click action to submit a download URL. The `/downloads` page lists all files in the download directory and lets you download them directly from the browser.

HTML and CSS courtesy of Claude; the purpose of this project is to write the API, not create web interfaces.

## How it works

1. You send a `POST /jobs` request with a download URL
2. The API validates the URL (scheme, allowed domains)
3. A `.crawljob` file is generated and dropped into a watched folder
4. JDownloader picks it up and starts the download automatically
5. Query `GET /api/files` to list completed downloads
6. Retrieve a specific file with `GET /download?filename=<name>`
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
  ghcr.io/frostbyte0x/crawljob-api:latest
```

### Build locally
This will start the web server on port 8080
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
| `ALLOWED_DOMAINS` | Allowed download domains | 1fichier.com,mega.nz | 
---

## API Reference

### `POST /jobs`

Submit a download URL.

**Request Body**
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

### `GET /api/files`

List all files and directories in the download folder.

**Response Body**
```json
[
  {
    "Name": "movie.mkv",
    "Type": "file",
    "Extension": ".MKV",
    "Size": "4.2 GB"
  },
  {
    "Name": "archive",
    "Type": "dir",
    "Extension": "DIR",
    "Size": "0 B"
  }
]
```

**Responses**

| Code | Description |
|---|---|
| `200 OK` | JSON array of files returned |
| `403 Forbidden` | Download folder cannot be accessed |

---

### `GET /download?filename=<name>`

Stream a file from the download folder to the client.

**Query Parameters**

| Parameter | Description |
|---|---|
| `filename` | Name of the file to download (must be within the download directory) |

**Responses**

| Code | Description |
|---|---|
| `200 OK` | File streamed as attachment |
| `403 Forbidden` | Path traversal attempt or folder inaccessible |
| `404 Not Found` | No filename provided |

---

## Allowed Domains
This can be changed in the Dockerfile configuration using ALLOWED_DOMAINS

Currently restricted to:
- `1fichier.com`
- `mega.nz`

> Contact the server owner or set your own domain list to extend this.

---

## Project Structure

```
crawljob-api/
├── main.go             # Server entrypoint
├── handler/
│   ├── job.go          # HTTP handler
│   ├── download_ui.go  # HTTP handler for /downloads (web interface)
│   ├── validator.go    # URL validation
│   ├── download.go     # API that returns a json array of files in the download directory
│   └── ui.go           # HTTP handler for / (web interface)
├── model/
│   ├── crawljob.go     # CrawlJob model + file generation
│   └── utils.go        # Helpers
└── Dockerfile
```

---

## License

MIT
