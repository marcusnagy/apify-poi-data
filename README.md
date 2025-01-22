# APIFY POI Data

---

## Table of Contents
- [Generate Server Certs for TLS](#generate-server-certs-for-tls)
- [Environment Variables](#environment-variables)
- [Setup and Run](#setup-and-run)
- [Makefile Commands](#makefile-commands)
- [Endpoints](#endpoints)

---

## Generate Server Certs for TLS

* `-x509`: Outputs a self-signed certificate.
* `-newkey rsa:4096`: Generates a new RSA key of 4096 bits.
* `-keyout server.key`: Outputs the private key to server.key.
* `-out server.crt`: Outputs the certificate to server.crt.
* `-days 365`: Sets the certificate validity to 365 days.
* `-nodes`: Prevents encrypting the private key.
* `-subj "/CN=localhost"`: Sets the Common Name (CN) to localhost.
```bash
# Generate a private key for the root CA
openssl genrsa -out rootCA.key 4096

# Create the root CA certificate
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.crt \
  -subj "/CN=RootCA"
```

```bash
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"
```

```bash
# Generate the server key and certificate signing request (CSR)
openssl req -new -newkey rsa:4096 -keyout server.key -out server.csr -nodes \
  -subj "/CN=poi-backend" \
  -addext "subjectAltName=DNS:localhost,DNS:poi-backend,IP:127.0.0.1"

# Sign the server CSR with the root CA
openssl x509 -req -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial \
  -out server.crt -days 365 -sha256 -extfile <(printf "subjectAltName=DNS:localhost,DNS:poi-backend,IP:127.0.0.1")
```

```bash
# Make pem
openssl x509 -in rootCA.crt -out rootCA.pem -outform PEM
```

## Environment Variables

Create a `.env` file in the root directory with the following content:

```
APIFY_KEY=""
APIFY_ACTOR_EXTRACTOR_ID=""
APIFY_ACTOR_SCRAPER_ID=""
DATABASE_USER="postgres"
DATABASE_PASSWORD="postgres"
```

## Setup and Run

1. **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. **Load environment variables:**
    ```sh
    source load_env.sh
    ```

3. **Start Docker containers:**
    ```sh
    make docker-up
    ```

4. **Run database migrations:**
    ```sh
    make migrate-up
    ```

5. **Generate Protobuf code:**
    ```sh
    make buf-generate
    ```

6. **Run the application:**
    ```sh
    make run
    ```

## Makefile Commands

- **Run the application:**
    ```sh
    make run
    ```

- **Start Docker containers:**
    ```sh
    make docker-up
    ```

- **Stop Docker containers:**
    ```sh
    make docker-down
    ```

- **Generate Protobuf code:**
    ```sh
    make buf-generate
    ```

- **Remove Protobuf code:**
    ```sh
    make buf-remove
    ```

- **Run database migrations:**
    ```sh
    make migrate-up
    ```

- **Rollback database migrations:**
    ```sh
    make migrate-down
    ```

- **Check migration version:**
    ```sh
    make migrate-version
    ```

- **Check migration status:**
    ```sh
    make migrate-status
    ```

- **Create a new migration:**
    ```sh
    make migration-new name=<migration_name>
    ```

- **Generate SQLC queries:**
    ```sh
    make sqlc
    ```

## Endpoints

### POI Service

- **List POIs in a Box:**
    ```
    GET /v1/poi/box
    ```

- **List POIs by H3 Cells:**
    ```
    GET /v1/poi/h3
    ```

- **List POIs Along Route with Category:**
    ```
    GET /v1/poi/route/category
    ```

### Maps Service

- **Search Google Maps Scraper:**
    ```
    POST /v1/maps/search/scraper
    ```

- **Search Google Maps Extractor:**
    ```
    POST /v1/maps/search/extractor
    ```

- **Insert Apify Dataset Items:**
    ```
    POST /v1/maps/dataset/insert
    ```

### Tripadvisor Service

- **Search Tripadvisor:**
    ```
    POST /v1/tripadvisor/search
    ```

---