# storj.io implementation

---

## Requirements

- [Golang](https://golang.org/doc/install) installed

## Runnning

1. Have a [storj](https://storj.io/) account S3 credentials ready.

    - [docs.storj.io/dcs/getting-started](https://docs.storj.io/dcs/getting-started)

2. Get the credentials a put them in a `.env` file based off of `.env.template`.

3. Run the server

    ```bash
    go run main.go
    ```

    Endpoints:
    list all avaliable buckets `GET http://localhost:8080/buckets`
    list all objects in bucket `http://localhost:8080/buckets/{bucketName}`
