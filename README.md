## Packaging Calculator
This is a monorepo where each project sits inside cmd directory.

To run test for packaging calculator cd in cmd/packagingcalculator and run
```bash
make test
```

## Getting Started with the monorepo
Install dependencies
```bash
make config
```

Before commiting run `make generate` to generate mocks.

### Running packagingcalculator locally
```bash
make dev_run
```

### Visit  [http://localhost](http://localhost)

### Sample Request
```bash
curl --location 'http://localhost:80/v1/products/packagecaclulator?items=1010'
```

### Sample insert request to change pack sizes
```bash
curl --location --request POST 'http://localhost:80/v1/products/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "pack_sizes": [
        250,500,1000,2000,5000
    ]
}'
```