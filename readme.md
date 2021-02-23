# Pho - Backend

## Top-Level Directories
- shared
  - portions of the backend shared across different services
- daemon
  - daemon process that:
    - ensures the API is healthy
    - dispatches workers
- api
  - the rest api for the frontend to query
- workers
  - job consumers and producers
