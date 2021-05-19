## Invoices API

This API is responsible for processing a call log in order to identify the invoice corresponding to a particular customer. For this, it takes into account all calls made by the client or those received by reverse charge. National and international calls are identified and have a differential price. In addition, the Friends API is used to make a particular discount in these cases.

### Test

```bash
go test ./... -v -race
```

### Build & Run

```bash
cd cmd/api
go build
./api
```

### Available resource

```bash
curl --location --request POST 'http://localhost:8080/invoices/process' \
--form 'input.csv=@"/Users/your/file/path/input.csv"' \
--form 'clientPhone="+5491167930920"' \
--form 'invoicePeriod="2020-01-01/2020-12-30"'
```

