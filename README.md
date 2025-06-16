# EquipQR

## Environment Setup

---

### Frontend

#### Running the Frontend (Development)

To run the frontend in development mode, follow these steps:

1. **Ensure Environment Setup**

   * Make sure that all dependencies and environment variables are correctly configured.

2. **Install Dependencies**

   * Navigate to the frontend project directory:

     ```bash
     cd repos/frontend
     ```
   * Install the necessary dependencies using `npm`:

     ```bash
     npm install
     ```

3. **Start the Development Server**

   * Once the dependencies are installed, start the development server:

     ```bash
     npm run dev
     ```

   This will launch the frontend application, and you should be able to access it in your browser at the designated development URL (typically `http://localhost:3000`).

---

### Backend

#### Running the Backend (Development)

To run the backend in development mode, use the following steps:

1. **Ensure Environment Setup**

   * Navigate to the backend directory:

     ```bash
     cd repos/backend
     ```
   * Ensure that all necessary environment variables and dependencies are set up correctly for local development.

2. **Run the Backend with `go run`**

   * Start the backend in development mode using the `go run` command:

     ```bash
     go run cmd/equipqr/main.go
     ```

   By default, the backend will be accessible at `http://localhost:8080`. This can be customized if needed by adjusting the relevant environment variables or configuration files.

#### Building and Starting the Backend (Production)

To build and run the backend in production mode:

1. **Build the Backend**

   * First, build the Go binary:

     ```bash
     cd repos/backend
     go build -o ./equipqr cmd/equipqr/main.go
     ```

2. **Run the Backend**

   * After building the binary, start the backend service:

     ```bash
     ./equipqr
     ```

   The backend will be live at `http://localhost:8080` by default, unless otherwise specified in the environment setup.

Ensure that your environment is properly configured with necessary configuration files and environment variables for a smooth deployment.
