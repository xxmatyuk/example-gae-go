# JWT token generation for GCP service account

## Supported algorithms are RS256 and ES256. Usage example:

```
export ALG=<YOUR_ALGORITHM>
export IAP_AUDIENCE_STRING=<YOUR_AUDIENCE_STRING>
export SERVICE_ACCOUNT_JSON_PATH=<PATH_TO_SERVICE_ACCOUNT_CREDS_JSON_FILE>
go run tokengen.go
```
