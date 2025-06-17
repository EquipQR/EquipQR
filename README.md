# EquipQR

## Environment Setup

Before running the application, you’ll need to configure environment variables. For convenience, a `.env.example` file is provided at the root of the project. It contains default values that mirror how EquipQR is run on-site.

### `.env.example` Overview (Short Explanation)

* **Backend Configuration**

  * `APP_HOST`, `APP_PORT`: Where the backend server listens.
  * `SSL_CERT_PATH`, `SSL_KEY_PATH`: Path to TLS cert/key (required for HTTPS).
  * `CORS_ALLOW_ORIGINS`, `CORS_ALLOW_HEADERS`: CORS policy for API access.

* **Frontend Configuration**

  * `VITE_HOST`, `VITE_PORT`: Where the frontend dev server runs.
  * `VITE_SSL_KEY_PATH`, `VITE_SSL_CERT_PATH`: TLS cert/key used by Vite in dev mode.

* **PostgreSQL Configuration**

  * `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`: Database credentials.
  * `POSTGRES_HOST`, `POSTGRES_PORT`: Database host and port (typically `localhost:5432`).

To initialize your environment:

```bash
cp .env.example .env
```

---

## SSL Certificates for HTTPS (Required)

EquipQR runs over **HTTPS by default**, even in development. This is critical for:

* **Camera Access**: Browsers **require a secure origin (HTTPS)** to access camera APIs used for QR scanning.
* **Progressive Web App (PWA) Installation**: Required for full PWA support on most platforms.

### Creating SSL Certificates with `mkcert`

We recommend using [`mkcert`](https://github.com/FiloSottile/mkcert) to generate trusted development certificates:

1. **Install mkcert**

   ```bash
   brew install mkcert           # macOS
   sudo pacman -S mkcert         # Arch Linux
   choco install mkcert          # Windows
   ```

2. **Create Local CA and Certs**

   ```bash
   mkcert -install
   mkcert localhost 127.0.0.1 ::1
   ```

3. **Place the certs in your certs directory and update `.env` paths accordingly**

---

## Docker Setup

### Running with Docker Compose

After setting up your `.env` file and SSL certificates:

1. **Start Services**

   ```bash
   docker compose up -d
   ```

2. **View Logs**

   ```bash
   docker compose logs -f
   ```

The `docker-compose.yml` in the root directory launches backend, frontend, and database containers.

---

## Backend

### Running the Backend (Development)

1. **Navigate to Backend Directory**

   ```bash
   cd repos/backend
   ```

2. **Run the Server**

   ```bash
   go run cmd/equipqr/main.go
   ```

   The backend is live at `https://localhost:8080` (or the configured port). Ensure the `SSL_CERT_PATH` and `SSL_KEY_PATH` are valid.

### Building & Running in Production

1. **Build the Binary**

   ```bash
   go build -o ./equipqr cmd/equipqr/main.go
   ```

2. **Run the Binary**

   ```bash
   ./equipqr
   ```

---

## Frontend

### Running the Frontend (Development)

1. **Navigate to Frontend Directory**

   ```bash
   cd repos/frontend
   ```

2. **Install Dependencies**

   ```bash
   npm install
   ```

3. **Start the Dev Server**

   ```bash
   npm run dev
   ```

   Visit the frontend at `https://localhost:3000` (or as specified in `.env`).

> **Note**: Because QR scanning requires camera access, the frontend must be served over **HTTPS**, or browser permissions will fail. You **must** use a trusted SSL certificate on your system. See the SSL section above.

---

## Progressive Web App (PWA) Support

EquipQR is built as a **Progressive Web App**, allowing it to be installed like a native app:

* On **mobile devices**, tap the **“Add to Home Screen”** prompt in supported browsers.
* On **desktop browsers** (e.g., Chrome or Edge), click the **install icon** in the address bar.

This allows full-screen offline-friendly access, QR scanning, and native-like performance on all platforms.
