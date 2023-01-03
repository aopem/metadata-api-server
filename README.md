# Metadata API Server

A golang metadata API server that supports CRUD operations with YAML metadata, along with a query capability.

## Usage

Run server on `localhost:8080`

> **NOTE**
> Can update endpoint/port by modifying the `addr` variable as desired in `main.go`.

```bash
cd /metadata-api-server
go run .
```

Run unit tests

```bash
cd /metadata-api-server
go test ./...
```

## Endpoints

-   **PUT** `/metadata`

    -   Create or update a metadata record. An empty ID field indicates a create operation.
    -   Request body format:

    > **IMPORTANT**
    > To update an existing record, make sure to pass the existing record's ID in the `id` field.

    ```yaml
    id: ""
    metadata:
        title: title
        version: version
        maintainers:
            - name: name
            email: valid@email.com
            - name: nametwo
            email: anothervalid@email.com
        company: company
        website: https://validurl.com
        source: source
        license: license
        description: |
            description
    ```

-   **DELETE** `/metadata/:id`

    -   Delete a metadata record using its ID.

-   **GET** `/metadata`

    -   Get a list of all metadata that is stored.

-   **GET** `/metadata/:id`

    -   Get a metadata record using its ID.

-   **PUT** `/metadata/query`
    -   Perform a search for a metadata record. All searches are performed using "AND" semantics, assuming that the resulting record(s) must match all query parameters.
    -   Parameters: `title, version, maintainerName, maintainerEmail, company, website, source, license, description`
