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


## Non-Golang Dependencies
- heif-convert
  - (may change because of [this bug][1])
- exiftool
- [epeg][2]


[1]: https://github.com/enricozb/pho-backend/issues/3
[2]: https://github.com/mattes/epeg
