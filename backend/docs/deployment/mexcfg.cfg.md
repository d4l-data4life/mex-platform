# MEx configuration

## Overview
| `met` | `idx` | `qry` | `cfg` | `aut` | Go struct field | Type | Secret | Environment variable | Default value | Title |
| ----- | ----- | ----- | ----- | ----- | --------------- | ---- | ------ | -------------------- | ------------- | ----- |
| ✅ | ✅ | ✅ | ✅ | ✅ | .TenantId | string |  |  `MEX_TENANT_ID` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.ReadTimeout | message |  |  `MEX_WEB_READ_TIMEOUT` | `'5s'` | HTTP service read timeout |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.WriteTimeout | message |  |  `MEX_WEB_WRITE_TIMEOUT` | `'5s'` | HTTP service write timeout |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.IdleTimeout | message |  |  `MEX_WEB_IDLE_TIMEOUT` | `'5s'` | HTTP service idle timeout |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.MaxHeaderBytes | int32 |  |  `MEX_WEB_MAX_HEADER_BYTES` | `'2097152'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.MaxBodyBytes | int64 |  |  `MEX_WEB_MAX_BODY_BYTES` | `'2097152'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.ApiHost | string |  |  `MEX_WEB_API_HOST` | `'0.0.0.0:3000'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.GrpcHost | string |  |  `MEX_WEB_GRPC_HOST` | `'0.0.0.0:9000'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.MetricsPath | string |  |  `MEX_WEB_METRICS_PATH` | `'/metrics'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.CaCerts.AdditionalCaCertsFiles | []string |  |  `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_FILES` | `'∅'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.CaCerts.AdditionalCaCertsPem | bytes | 🔒 | ❗ `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_PEM_B64` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.CaCerts.AccessAttempts | uint32 |  |  `MEX_WEB_CA_CERTS_ACCESS_ATTEMPTS` | `'20'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.CaCerts.AccessPause | message |  |  `MEX_WEB_CA_CERTS_ACCESS_PAUSE` | `'2s'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.IpFilter.Enabled | bool |  |  `MEX_WEB_IP_FILTER_ENABLED` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.IpFilter.AllowedIps | []string |  |  `MEX_WEB_IP_FILTER_ALLOWED_IPS` | `'∅'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.RateLimiting.Enabled | bool |  |  `MEX_WEB_RATE_LIMITING_ENABLED` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.RateLimiting.Period | message |  |  `MEX_WEB_RATE_LIMITING_PERIOD` | `'1s'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.RateLimiting.Limit | int64 |  |  `MEX_WEB_RATE_LIMITING_LIMIT` | `'100'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Web.RateLimiting.ClientIpHeader | string |  |  `MEX_WEB_RATE_LIMITING_CLIENT_IP_HEADER` | `'X-Real-Ip'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.User | string |  |  `MEX_DB_USER` | `'postgres'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.Password | string | 🔒 |  `MEX_DB_PASSWORD` | _none_ |  |
| ✅ | ✅ | ✅ |  |  | .Db.Hostname | string |  |  `MEX_DB_HOSTNAME` | `'localhost'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.Port | uint32 |  |  `MEX_DB_PORT` | `'5432'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.Name | string |  |  `MEX_DB_NAME` | `'postgres'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.SearchPath | []string |  |  `MEX_DB_SEARCH_PATH` | `'mex,public'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.SslMode | string |  |  `MEX_DB_SSL_MODE` | `'verify-full'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.ConnectionAttempts | uint32 |  |  `MEX_DB_CONNECTION_ATTEMPTS` | `'10'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.ConnectionPause | message |  |  `MEX_DB_CONNECTION_PAUSE` | `'2s'` |  |
| ✅ | ✅ | ✅ |  |  | .Db.SlowThreshold | message |  |  `MEX_DB_SLOW_THRESHOLD` | `'200ms'` |  |
|  | ✅ | ✅ |  |  | .Solr.Origin | string |  |  `MEX_SOLR_ORIGIN` | `'http://localhost:8983'` |  |
|  | ✅ | ✅ |  |  | .Solr.Collection | string |  |  `MEX_SOLR_COLLECTION` | `'mex'` |  |
|  | ✅ | ✅ |  |  | .Solr.ConfigsetName | string |  |  `MEX_SOLR_CONFIGSET_NAME` | `'mex_d4l'` |  |
|  | ✅ | ✅ |  |  | .Solr.ConnectionAttempts | uint32 |  |  `MEX_SOLR_CONNECTION_ATTEMPTS` | `'10'` |  |
|  | ✅ | ✅ |  |  | .Solr.ConnectionPause | message |  |  `MEX_SOLR_CONNECTION_PAUSE` | `'2s'` |  |
|  | ✅ | ✅ |  |  | .Solr.BasicauthUser | string |  |  `MEX_SOLR_BASICAUTH_USER` | _none_ |  |
|  | ✅ | ✅ |  |  | .Solr.BasicauthPassword | string | 🔒 |  `MEX_SOLR_BASICAUTH_PASSWORD` | _none_ |  |
|  | ✅ | ✅ |  |  | .Solr.IndexBatchSize | uint32 |  |  `MEX_SOLR_INDEX_BATCH_SIZE` | `'100'` |  |
|  | ✅ | ✅ |  |  | .Solr.CommitWithin | message |  |  `MEX_SOLR_COMMIT_WITHIN` | `'1000ms'` |  |
|  | ✅ | ✅ |  |  | .Solr.ReplicationFactor | uint32 |  |  `MEX_SOLR_REPLICATION_FACTOR` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.Hostname | string |  |  `MEX_REDIS_HOSTNAME` | `'localhost'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.Port | uint32 |  |  `MEX_REDIS_PORT` | `'6379'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.Password | string | 🔒 |  `MEX_REDIS_PASSWORD` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.Db | uint32 |  |  `MEX_REDIS_DB` | `'1'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.ConnectionAttempts | uint32 |  |  `MEX_REDIS_CONNECTION_ATTEMPTS` | `'10'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.ConnectionPause | message |  |  `MEX_REDIS_CONNECTION_PAUSE` | `'2s'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.ShutdownGracePeriod | message |  |  `MEX_REDIS_SHUTDOWN_GRACE_PERIOD` | `'200ms'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.UseTls | bool |  |  `MEX_REDIS_USE_TLS` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Redis.PubSubPrefix | string |  |  `MEX_REDIS_PUB_SUB_PREFIX` | `'mex'` |  |
| ✅ | ✅ | ✅ |  | ✅ | .Oauth.ClientId | string |  |  `MEX_OAUTH_CLIENT_ID` | _none_ |  |
| ✅ | ✅ | ✅ |  | ✅ | .Oauth.ProducerGroupId | string |  |  `MEX_OAUTH_PRODUCER_GROUP_ID` | _none_ |  |
| ✅ | ✅ | ✅ |  | ✅ | .Oauth.ConsumerGroupId | string |  |  `MEX_OAUTH_CONSUMER_GROUP_ID` | _none_ |  |
| ✅ | ✅ | ✅ |  |  | .Oauth.InternalAuthServiceHostname | string |  |  `MEX_OAUTH_INTERNAL_AUTH_SERVICE_HOSTNAME` | _none_ |  |
|  |  |  |  | ✅ | .Oauth.Server.Enabled | bool |  |  `MEX_OAUTH_SERVER_ENABLED` | `'false'` |  |
|  |  |  |  | ✅ | .Oauth.Server.ClientSecrets | []string | 🔒 |  `MEX_OAUTH_SERVER_CLIENT_SECRETS` | _none_ |  |
|  |  |  |  | ✅ | .Oauth.Server.RedirectUris | []string |  |  `MEX_OAUTH_SERVER_REDIRECT_URIS` | _none_ |  |
|  |  |  |  | ✅ | .Oauth.Server.GrantFlows | []string |  |  `MEX_OAUTH_SERVER_GRANT_FLOWS` | `'client_credentials,authorization_code,refresh_token'` |  |
|  |  |  |  | ✅ | .Oauth.Server.SigningPrivateKeyFile | string |  |  `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_FILE` | _none_ |  |
|  |  |  |  | ✅ | .Oauth.Server.KeyFileAccessAttempts | uint32 |  |  `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_ATTEMPTS` | `'20'` |  |
|  |  |  |  | ✅ | .Oauth.Server.KeyFileAccessPause | message |  |  `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_PAUSE` | `'2s'` |  |
|  |  |  |  | ✅ | .Oauth.Server.SigningPrivateKeyPem | bytes | 🔒 | ❗ `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_PEM_B64` | _none_ |  |
|  |  |  |  | ✅ | .Oauth.Server.SignatureAlg | string |  |  `MEX_OAUTH_SERVER_SIGNATURE_ALG` | `'RS256'` |  |
|  |  |  |  | ✅ | .Oauth.Server.AuthCodeValidity | message |  |  `MEX_OAUTH_SERVER_AUTH_CODE_VALIDITY` | `'1m'` |  |
|  |  |  |  | ✅ | .Oauth.Server.AccessTokenValidity | message |  |  `MEX_OAUTH_SERVER_ACCESS_TOKEN_VALIDITY` | `'1h'` |  |
|  |  |  |  | ✅ | .Oauth.Server.RefreshTokenValidity | message |  |  `MEX_OAUTH_SERVER_REFRESH_TOKEN_VALIDITY` | `'12h'` |  |
|  |  |  |  |  | .Codings.BundleUri | string |  |  `MEX_CODINGS_BUNDLE_URI` | _none_ |  |
| ✅ | ✅ | ✅ |  |  | .EntityTypes.RepoType | enum |  |  `MEX_ENTITY_TYPES_REPO_TYPE` | `'CACHED'` |  |
| ✅ | ✅ | ✅ |  |  | .FieldDefs.RepoType | enum |  |  `MEX_FIELD_DEFS_REPO_TYPE` | `'CACHED'` |  |
|  |  | ✅ |  |  | .SearchConfig.RepoType | enum |  |  `MEX_SEARCH_CONFIG_REPO_TYPE` | `'CACHED'` |  |
| ✅ | ✅ | ✅ |  |  | .Jwks.RemoteKeysUri | string |  |  `MEX_JWKS_REMOTE_KEYS_URI` | _none_ |  |
| ✅ | ✅ | ✅ |  |  | .Jwks.ConnectionAttempts | uint32 |  |  `MEX_JWKS_CONNECTION_ATTEMPTS` | `'20'` |  |
| ✅ | ✅ | ✅ |  |  | .Jwks.ConnectionPause | message |  |  `MEX_JWKS_CONNECTION_PAUSE` | `'2s'` |  |
| ✅ | ✅ |  | ✅ |  | .Jobs.Expiration | message |  |  `MEX_JOBS_EXPIRATION` | `'5m'` |  |
| ✅ | ✅ |  |  |  | .AutoIndexer.SetExpiration | message |  |  `MEX_AUTO_INDEXER_SET_EXPIRATION` | `'5m'` |  |
| ✅ |  |  |  |  | .Indexing.DuplicationDetectionAlgorithm | enum |  | ❗ `MEX_SERVICES_DUPLICATE_DETECTION_ALGORITHM` | `'LATEST_ONLY'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Logging.LogLevelGrpc | string |  |  `MEX_LOGGING_LOG_LEVEL_GRPC` | `'warn'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Logging.RedactPersonalFields | bool |  |  `MEX_LOGGING_REDACT_PERSONAL_FIELDS` | `'true'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Logging.RedactQueryParams | []string |  |  `MEX_LOGGING_REDACT_QUERY_PARAMS` | `'code_challenge,state'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Logging.TraceEnabled | bool |  |  `MEX_LOGGING_TRACE_ENABLED` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Logging.TraceSecret | string | 🔒 |  `MEX_LOGGING_TRACE_SECRET` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Telemetry.PingerUpdateInterval | message |  |  `MEX_TELEMETRY_PINGER_UPDATE_INTERVAL` | `'15s'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Telemetry.StatusUpdateInterval | message |  |  `MEX_TELEMETRY_STATUS_UPDATE_INTERVAL` | `'3s'` |  |
|  |  |  | ✅ |  | .Auth.ApiKeysRoles | bytes | 🔒 | ❗ `MEX_AUTH_API_KEYS_ROLES_B64` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Services.BiEventsFilter.Origin | string |  |  `MEX_SERVICES_BI_EVENTS_FILTER_ORIGIN` | _none_ |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Services.BiEventsFilter.Path | string |  |  `MEX_SERVICES_BI_EVENTS_FILTER_PATH` | `'/api/v1/events'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Services.BiEventsFilter.Secret | string | 🔒 | ❗ `MEX_SERVICES_BI_EVENTS_FILTER_SECRET` | _none_ | BI events filter API secret |
| ✅ | ✅ |  |  |  | .Services.Blobs.MasterTableName | string |  |  `MEX_SERVICES_BLOBS_MASTER_TABLE_NAME` | `'blob_store'` |  |
| ✅ | ✅ | ✅ |  |  | .Services.Config.Origin | string |  |  `MEX_SERVICES_CONFIG_ORIGIN` | _none_ |  |
|  |  |  | ✅ |  | .Services.Config.EnvPath | string |  |  `MEX_SERVICES_CONFIG_ENV_PATH` | `'/'` |  |
|  |  |  | ✅ |  | .Services.Config.ApiKeys | []string | 🔒 |  `MEX_SERVICES_CONFIG_API_KEYS` | _none_ |  |
|  |  |  | ✅ |  | .Services.Config.Github.RepoName | string |  |  `MEX_SERVICES_CONFIG_GITHUB_REPO_NAME` | _none_ |  |
|  |  |  | ✅ |  | .Services.Config.Github.DefaultBranchName | string |  |  `MEX_SERVICES_CONFIG_GITHUB_DEFAULT_BRANCH_NAME` | `'main'` |  |
|  |  |  | ✅ |  | .Services.Config.Github.DeployKeyPem | bytes | 🔒 | ❗ `MEX_SERVICES_CONFIG_GITHUB_DEPLOY_KEY_PEM_B64` | _none_ |  |
|  |  |  | ✅ |  | .Services.Config.UpdateTimeout | message |  |  `MEX_SERVICES_CONFIG_UPDATE_TIMEOUT` | `'180s'` | Maximum duration a config update may take |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.Search.ToleratePartialFailures | bool |  |  `MEX_STRICTNESS_SEARCH_TOLERATE_PARTIAL_FAILURES` | `'true'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.StrictJsonParsing.Auth | bool |  |  `MEX_STRICTNESS_STRICT_JSON_PARSING_AUTH` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.StrictJsonParsing.Config | bool |  |  `MEX_STRICTNESS_STRICT_JSON_PARSING_CONFIG` | `'false'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.StrictJsonParsing.Index | bool |  |  `MEX_STRICTNESS_STRICT_JSON_PARSING_INDEX` | `'true'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.StrictJsonParsing.Metadata | bool |  |  `MEX_STRICTNESS_STRICT_JSON_PARSING_METADATA` | `'true'` |  |
| ✅ | ✅ | ✅ | ✅ | ✅ | .Strictness.StrictJsonParsing.Query | bool |  |  `MEX_STRICTNESS_STRICT_JSON_PARSING_QUERY` | `'true'` |  |
| ✅ |  |  |  |  | .Notify.EmailerType | enum |  |  `MEX_NOTIFY_EMAILER_TYPE` | `'MOCKMAILER'` |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.OriginOauth | string |  |  `MEX_NOTIFY_FLOWMAILER_ORIGIN_OAUTH` | `'https://login.flowmailer.net'` |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.OriginApi | string |  |  `MEX_NOTIFY_FLOWMAILER_ORIGIN_API` | `'https://api.flowmailer.net'` |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.ClientId | string |  |  `MEX_NOTIFY_FLOWMAILER_CLIENT_ID` | _none_ |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.ClientSecret | string | 🔒 |  `MEX_NOTIFY_FLOWMAILER_CLIENT_SECRET` | _none_ |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.AccountId | string |  |  `MEX_NOTIFY_FLOWMAILER_ACCOUNT_ID` | _none_ |  |
| ✅ |  |  |  |  | .Notify.Flowmailer.NoreplyEmailAddress | string |  |  `MEX_NOTIFY_FLOWMAILER_NOREPLY_EMAIL_ADDRESS` | `'noreply@data4life.care'` |  |
## Configuration details
### `MEX_TENANT_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.TenantId` |
| Environment variable: | `MEX_TENANT_ID`  |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_READ_TIMEOUT`: HTTP service read timeout
#### Summary

This value is the maximum duration for reading the entire request, including the body
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.ReadTimeout` |
| Environment variable: | `MEX_WEB_READ_TIMEOUT`  |
| Default value: | `'5s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_WRITE_TIMEOUT`: HTTP service write timeout
#### Summary

This value is the maximum duration before timing out writes of the response
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.WriteTimeout` |
| Environment variable: | `MEX_WEB_WRITE_TIMEOUT`  |
| Default value: | `'5s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_IDLE_TIMEOUT`: HTTP service idle timeout
#### Summary

This value is the maximum amount of time to wait for the next request when keep-alives are enabled
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.IdleTimeout` |
| Environment variable: | `MEX_WEB_IDLE_TIMEOUT`  |
| Default value: | `'5s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_MAX_HEADER_BYTES`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.MaxHeaderBytes` |
| Environment variable: | `MEX_WEB_MAX_HEADER_BYTES`  |
| Default value: | `'2097152'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_MAX_BODY_BYTES`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.MaxBodyBytes` |
| Environment variable: | `MEX_WEB_MAX_BODY_BYTES`  |
| Default value: | `'2097152'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_API_HOST`: 
#### Summary

Host and port for the exposed HTTP service
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.ApiHost` |
| Environment variable: | `MEX_WEB_API_HOST`  |
| Default value: | `'0.0.0.0:3000'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_GRPC_HOST`: 
#### Summary

Host and port for the gRPC service whose methods are exposed via an HTTP-gRPC gateway under `Web.APIHost`
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.GrpcHost` |
| Environment variable: | `MEX_WEB_GRPC_HOST`  |
| Default value: | `'0.0.0.0:9000'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_METRICS_PATH`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.MetricsPath` |
| Environment variable: | `MEX_WEB_METRICS_PATH`  |
| Default value: | `'/metrics'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_FILES`: 
#### Summary

Additional CA certificates to consider when making HTTPS or other TLS-protected requests
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.CaCerts.AdditionalCaCertsFiles` |
| Environment variable: | `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_FILES`  |
| Default value: | `'∅'` |
| Used by: | <ul><li>_all_</li></ul> |
#### Description

The files must be PEM files and the single string parameter must be a set of base64-encoded PEM blocks.
All certificates of all such PEM blocks are then added to the trusted certificates for TLS.

----
### `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_PEM_B64`: 
#### Summary

Additional certificates specified in PEM format
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.CaCerts.AdditionalCaCertsPem` |
| Environment variable: | `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_PEM_B64`  |
| Vault source variable: | `MEX_WEB_CA_CERTS_ADDITIONAL_CA_CERTS_PEM` |
| Secret: | **yes** |
| Used by: | <ul><li>_all_</li></ul> |
#### Description

This field can be set to a string that is the content of a PEM file.
Multiple PEM blocks are possible so that multiple certificates can be specified.

----
### `MEX_WEB_CA_CERTS_ACCESS_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.CaCerts.AccessAttempts` |
| Environment variable: | `MEX_WEB_CA_CERTS_ACCESS_ATTEMPTS`  |
| Default value: | `'20'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_CA_CERTS_ACCESS_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.CaCerts.AccessPause` |
| Environment variable: | `MEX_WEB_CA_CERTS_ACCESS_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_IP_FILTER_ENABLED`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.IpFilter.Enabled` |
| Environment variable: | `MEX_WEB_IP_FILTER_ENABLED`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_IP_FILTER_ALLOWED_IPS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.IpFilter.AllowedIps` |
| Environment variable: | `MEX_WEB_IP_FILTER_ALLOWED_IPS`  |
| Default value: | `'∅'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_RATE_LIMITING_ENABLED`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.RateLimiting.Enabled` |
| Environment variable: | `MEX_WEB_RATE_LIMITING_ENABLED`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_RATE_LIMITING_PERIOD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.RateLimiting.Period` |
| Environment variable: | `MEX_WEB_RATE_LIMITING_PERIOD`  |
| Default value: | `'1s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_RATE_LIMITING_LIMIT`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.RateLimiting.Limit` |
| Environment variable: | `MEX_WEB_RATE_LIMITING_LIMIT`  |
| Default value: | `'100'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_WEB_RATE_LIMITING_CLIENT_IP_HEADER`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Web.RateLimiting.ClientIpHeader` |
| Environment variable: | `MEX_WEB_RATE_LIMITING_CLIENT_IP_HEADER`  |
| Default value: | `'X-Real-Ip'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_DB_USER`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.User` |
| Environment variable: | `MEX_DB_USER`  |
| Default value: | `'postgres'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_PASSWORD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.Password` |
| Environment variable: | `MEX_DB_PASSWORD`  |
| Secret: | **yes** |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_HOSTNAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.Hostname` |
| Environment variable: | `MEX_DB_HOSTNAME`  |
| Default value: | `'localhost'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_PORT`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.Port` |
| Environment variable: | `MEX_DB_PORT`  |
| Default value: | `'5432'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_NAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.Name` |
| Environment variable: | `MEX_DB_NAME`  |
| Default value: | `'postgres'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_SEARCH_PATH`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.SearchPath` |
| Environment variable: | `MEX_DB_SEARCH_PATH`  |
| Default value: | `'mex,public'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_SSL_MODE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.SslMode` |
| Environment variable: | `MEX_DB_SSL_MODE`  |
| Default value: | `'verify-full'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_CONNECTION_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.ConnectionAttempts` |
| Environment variable: | `MEX_DB_CONNECTION_ATTEMPTS`  |
| Default value: | `'10'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_CONNECTION_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.ConnectionPause` |
| Environment variable: | `MEX_DB_CONNECTION_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_DB_SLOW_THRESHOLD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Db.SlowThreshold` |
| Environment variable: | `MEX_DB_SLOW_THRESHOLD`  |
| Default value: | `'200ms'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_ORIGIN`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.Origin` |
| Environment variable: | `MEX_SOLR_ORIGIN`  |
| Default value: | `'http://localhost:8983'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_COLLECTION`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.Collection` |
| Environment variable: | `MEX_SOLR_COLLECTION`  |
| Default value: | `'mex'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_CONFIGSET_NAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.ConfigsetName` |
| Environment variable: | `MEX_SOLR_CONFIGSET_NAME`  |
| Default value: | `'mex_rki'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_CONNECTION_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.ConnectionAttempts` |
| Environment variable: | `MEX_SOLR_CONNECTION_ATTEMPTS`  |
| Default value: | `'10'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_CONNECTION_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.ConnectionPause` |
| Environment variable: | `MEX_SOLR_CONNECTION_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_BASICAUTH_USER`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.BasicauthUser` |
| Environment variable: | `MEX_SOLR_BASICAUTH_USER`  |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_BASICAUTH_PASSWORD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.BasicauthPassword` |
| Environment variable: | `MEX_SOLR_BASICAUTH_PASSWORD`  |
| Secret: | **yes** |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_INDEX_BATCH_SIZE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.IndexBatchSize` |
| Environment variable: | `MEX_SOLR_INDEX_BATCH_SIZE`  |
| Default value: | `'100'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_COMMIT_WITHIN`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.CommitWithin` |
| Environment variable: | `MEX_SOLR_COMMIT_WITHIN`  |
| Default value: | `'1000ms'` |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_SOLR_REPLICATION_FACTOR`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Solr.ReplicationFactor` |
| Environment variable: | `MEX_SOLR_REPLICATION_FACTOR`  |
| Used by: | <ul><li>index</li><li>query</li></ul> |

----
### `MEX_REDIS_HOSTNAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.Hostname` |
| Environment variable: | `MEX_REDIS_HOSTNAME`  |
| Default value: | `'localhost'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_PORT`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.Port` |
| Environment variable: | `MEX_REDIS_PORT`  |
| Default value: | `'6379'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_PASSWORD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.Password` |
| Environment variable: | `MEX_REDIS_PASSWORD`  |
| Secret: | **yes** |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_DB`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.Db` |
| Environment variable: | `MEX_REDIS_DB`  |
| Default value: | `'1'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_CONNECTION_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.ConnectionAttempts` |
| Environment variable: | `MEX_REDIS_CONNECTION_ATTEMPTS`  |
| Default value: | `'10'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_CONNECTION_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.ConnectionPause` |
| Environment variable: | `MEX_REDIS_CONNECTION_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_SHUTDOWN_GRACE_PERIOD`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.ShutdownGracePeriod` |
| Environment variable: | `MEX_REDIS_SHUTDOWN_GRACE_PERIOD`  |
| Default value: | `'200ms'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_USE_TLS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.UseTls` |
| Environment variable: | `MEX_REDIS_USE_TLS`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_REDIS_PUB_SUB_PREFIX`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Redis.PubSubPrefix` |
| Environment variable: | `MEX_REDIS_PUB_SUB_PREFIX`  |
| Default value: | `'mex'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_OAUTH_CLIENT_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.ClientId` |
| Environment variable: | `MEX_OAUTH_CLIENT_ID`  |
| Used by: | <ul><li>auth</li><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_OAUTH_PRODUCER_GROUP_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.ProducerGroupId` |
| Environment variable: | `MEX_OAUTH_PRODUCER_GROUP_ID`  |
| Used by: | <ul><li>auth</li><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_OAUTH_CONSUMER_GROUP_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.ConsumerGroupId` |
| Environment variable: | `MEX_OAUTH_CONSUMER_GROUP_ID`  |
| Used by: | <ul><li>auth</li><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_OAUTH_INTERNAL_AUTH_SERVICE_HOSTNAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.InternalAuthServiceHostname` |
| Environment variable: | `MEX_OAUTH_INTERNAL_AUTH_SERVICE_HOSTNAME`  |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_OAUTH_SERVER_ENABLED`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.Enabled` |
| Environment variable: | `MEX_OAUTH_SERVER_ENABLED`  |
| Default value: | `'false'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_CLIENT_SECRETS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.ClientSecrets` |
| Environment variable: | `MEX_OAUTH_SERVER_CLIENT_SECRETS`  |
| Secret: | **yes** |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_REDIRECT_URIS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.RedirectUris` |
| Environment variable: | `MEX_OAUTH_SERVER_REDIRECT_URIS`  |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_GRANT_FLOWS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.GrantFlows` |
| Environment variable: | `MEX_OAUTH_SERVER_GRANT_FLOWS`  |
| Default value: | `'client_credentials,authorization_code,refresh_token'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_FILE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.SigningPrivateKeyFile` |
| Environment variable: | `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_FILE`  |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.KeyFileAccessAttempts` |
| Environment variable: | `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_ATTEMPTS`  |
| Default value: | `'20'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.KeyFileAccessPause` |
| Environment variable: | `MEX_OAUTH_SERVER_KEY_FILE_ACCESS_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_PEM_B64`: 
#### Summary

Private key in PEM format that will be used for signing JWTs
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.SigningPrivateKeyPem` |
| Environment variable: | `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_PEM_B64`  |
| Vault source variable: | `MEX_OAUTH_SERVER_SIGNING_PRIVATE_KEY_PEM` |
| Secret: | **yes** |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_SIGNATURE_ALG`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.SignatureAlg` |
| Environment variable: | `MEX_OAUTH_SERVER_SIGNATURE_ALG`  |
| Default value: | `'RS256'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_AUTH_CODE_VALIDITY`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.AuthCodeValidity` |
| Environment variable: | `MEX_OAUTH_SERVER_AUTH_CODE_VALIDITY`  |
| Default value: | `'1m'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_ACCESS_TOKEN_VALIDITY`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.AccessTokenValidity` |
| Environment variable: | `MEX_OAUTH_SERVER_ACCESS_TOKEN_VALIDITY`  |
| Default value: | `'1h'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_OAUTH_SERVER_REFRESH_TOKEN_VALIDITY`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Oauth.Server.RefreshTokenValidity` |
| Environment variable: | `MEX_OAUTH_SERVER_REFRESH_TOKEN_VALIDITY`  |
| Default value: | `'12h'` |
| Used by: | <ul><li>auth</li></ul> |

----
### `MEX_CODINGS_BUNDLE_URI`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Codings.BundleUri` |
| Environment variable: | `MEX_CODINGS_BUNDLE_URI`  |
| Used by: | <ul></ul> |

----
### `MEX_ENTITY_TYPES_REPO_TYPE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.EntityTypes.RepoType` |
| Environment variable: | `MEX_ENTITY_TYPES_REPO_TYPE`  |
| Default value: | `'CACHED'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_FIELD_DEFS_REPO_TYPE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.FieldDefs.RepoType` |
| Environment variable: | `MEX_FIELD_DEFS_REPO_TYPE`  |
| Default value: | `'CACHED'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_SEARCH_CONFIG_REPO_TYPE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.SearchConfig.RepoType` |
| Environment variable: | `MEX_SEARCH_CONFIG_REPO_TYPE`  |
| Default value: | `'CACHED'` |
| Used by: | <ul><li>query</li></ul> |

----
### `MEX_JWKS_REMOTE_KEYS_URI`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Jwks.RemoteKeysUri` |
| Environment variable: | `MEX_JWKS_REMOTE_KEYS_URI`  |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_JWKS_CONNECTION_ATTEMPTS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Jwks.ConnectionAttempts` |
| Environment variable: | `MEX_JWKS_CONNECTION_ATTEMPTS`  |
| Default value: | `'20'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_JWKS_CONNECTION_PAUSE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Jwks.ConnectionPause` |
| Environment variable: | `MEX_JWKS_CONNECTION_PAUSE`  |
| Default value: | `'2s'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>query</li></ul> |

----
### `MEX_JOBS_EXPIRATION`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Jobs.Expiration` |
| Environment variable: | `MEX_JOBS_EXPIRATION`  |
| Default value: | `'5m'` |
| Used by: | <ul><li>metadata</li><li>index</li><li>config</li></ul> |

----
### `MEX_AUTO_INDEXER_SET_EXPIRATION`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.AutoIndexer.SetExpiration` |
| Environment variable: | `MEX_AUTO_INDEXER_SET_EXPIRATION`  |
| Default value: | `'5m'` |
| Used by: | <ul><li>metadata</li><li>index</li></ul> |

----
### `MEX_SERVICES_DUPLICATE_DETECTION_ALGORITHM`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Indexing.DuplicationDetectionAlgorithm` |
| Environment variable: | `MEX_SERVICES_DUPLICATE_DETECTION_ALGORITHM` (Note the name deviation!) |
| Default value: | `'LATEST_ONLY'` |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_LOGGING_LOG_LEVEL_GRPC`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Logging.LogLevelGrpc` |
| Environment variable: | `MEX_LOGGING_LOG_LEVEL_GRPC`  |
| Default value: | `'warn'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_LOGGING_REDACT_PERSONAL_FIELDS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Logging.RedactPersonalFields` |
| Environment variable: | `MEX_LOGGING_REDACT_PERSONAL_FIELDS`  |
| Default value: | `'true'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_LOGGING_REDACT_QUERY_PARAMS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Logging.RedactQueryParams` |
| Environment variable: | `MEX_LOGGING_REDACT_QUERY_PARAMS`  |
| Default value: | `'code_challenge,state'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_LOGGING_TRACE_ENABLED`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Logging.TraceEnabled` |
| Environment variable: | `MEX_LOGGING_TRACE_ENABLED`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_LOGGING_TRACE_SECRET`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Logging.TraceSecret` |
| Environment variable: | `MEX_LOGGING_TRACE_SECRET`  |
| Secret: | **yes** |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_TELEMETRY_PINGER_UPDATE_INTERVAL`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Telemetry.PingerUpdateInterval` |
| Environment variable: | `MEX_TELEMETRY_PINGER_UPDATE_INTERVAL`  |
| Default value: | `'15s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_TELEMETRY_STATUS_UPDATE_INTERVAL`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Telemetry.StatusUpdateInterval` |
| Environment variable: | `MEX_TELEMETRY_STATUS_UPDATE_INTERVAL`  |
| Default value: | `'3s'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_AUTH_API_KEYS_ROLES_B64`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Auth.ApiKeysRoles` |
| Environment variable: | `MEX_AUTH_API_KEYS_ROLES_B64`  |
| Vault source variable: | `MEX_AUTH_API_KEYS_ROLES` |
| Secret: | **yes** |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_BI_EVENTS_FILTER_ORIGIN`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.BiEventsFilter.Origin` |
| Environment variable: | `MEX_SERVICES_BI_EVENTS_FILTER_ORIGIN`  |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_SERVICES_BI_EVENTS_FILTER_PATH`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.BiEventsFilter.Path` |
| Environment variable: | `MEX_SERVICES_BI_EVENTS_FILTER_PATH`  |
| Default value: | `'/api/v1/events'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_SERVICES_BI_EVENTS_FILTER_SECRET`: BI events filter API secret
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.BiEventsFilter.Secret` |
| Environment variable: | `MEX_SERVICES_BI_EVENTS_FILTER_SECRET`  |
| Vault source variable: | `BI_EVENTS_FILTER_SECRET` |
| Secret: | **yes** |
| Used by: | <ul><li>_all_</li></ul> |
#### Description

Note: The source is a value coming from a Vault common secret `apps/<ENV>/phdp/common`.

----
### `MEX_SERVICES_BLOBS_MASTER_TABLE_NAME`: 
#### Summary

The blob store uses the same database as configured under `DB` above
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Blobs.MasterTableName` |
| Environment variable: | `MEX_SERVICES_BLOBS_MASTER_TABLE_NAME`  |
| Default value: | `'blob_store'` |
| Used by: | <ul><li>metadata</li><li>index</li></ul> |

----
### `MEX_SERVICES_CONFIG_ORIGIN`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.Origin` |
| Environment variable: | `MEX_SERVICES_CONFIG_ORIGIN`  |
| Used by: | <ul><li>metadata</li><li>query</li><li>index</li></ul> |

----
### `MEX_SERVICES_CONFIG_ENV_PATH`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.EnvPath` |
| Environment variable: | `MEX_SERVICES_CONFIG_ENV_PATH`  |
| Default value: | `'/'` |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_CONFIG_API_KEYS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.ApiKeys` |
| Environment variable: | `MEX_SERVICES_CONFIG_API_KEYS`  |
| Secret: | **yes** |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_CONFIG_GITHUB_REPO_NAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.Github.RepoName` |
| Environment variable: | `MEX_SERVICES_CONFIG_GITHUB_REPO_NAME`  |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_CONFIG_GITHUB_DEFAULT_BRANCH_NAME`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.Github.DefaultBranchName` |
| Environment variable: | `MEX_SERVICES_CONFIG_GITHUB_DEFAULT_BRANCH_NAME`  |
| Default value: | `'main'` |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_CONFIG_GITHUB_DEPLOY_KEY_PEM_B64`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.Github.DeployKeyPem` |
| Environment variable: | `MEX_SERVICES_CONFIG_GITHUB_DEPLOY_KEY_PEM_B64`  |
| Vault source variable: | `MEX_SERVICES_CONFIG_GITHUB_DEPLOY_KEY_PEM` |
| Secret: | **yes** |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_SERVICES_CONFIG_UPDATE_TIMEOUT`: Maximum duration a config update may take
#### Summary

If not all services report GREEN with the corresponding config hash after this time, a config updae is considered failed.
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Services.Config.UpdateTimeout` |
| Environment variable: | `MEX_SERVICES_CONFIG_UPDATE_TIMEOUT`  |
| Default value: | `'180s'` |
| Used by: | <ul><li>config</li></ul> |

----
### `MEX_STRICTNESS_SEARCH_TOLERATE_PARTIAL_FAILURES`: 
#### Summary

If true, certain partial failures of Solr search do not cause a 500 response
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.Search.ToleratePartialFailures` |
| Environment variable: | `MEX_STRICTNESS_SEARCH_TOLERATE_PARTIAL_FAILURES`  |
| Default value: | `'true'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_STRICTNESS_STRICT_JSON_PARSING_AUTH`: 
#### Summary

If true, unknown properties in data handled by auth service will cause an error
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.StrictJsonParsing.Auth` |
| Environment variable: | `MEX_STRICTNESS_STRICT_JSON_PARSING_AUTH`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_STRICTNESS_STRICT_JSON_PARSING_CONFIG`: 
#### Summary

If true, unknown properties in data handled by config service will cause an error
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.StrictJsonParsing.Config` |
| Environment variable: | `MEX_STRICTNESS_STRICT_JSON_PARSING_CONFIG`  |
| Default value: | `'false'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_STRICTNESS_STRICT_JSON_PARSING_INDEX`: 
#### Summary

If true, unknown properties in data handled by index service will cause an error
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.StrictJsonParsing.Index` |
| Environment variable: | `MEX_STRICTNESS_STRICT_JSON_PARSING_INDEX`  |
| Default value: | `'true'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_STRICTNESS_STRICT_JSON_PARSING_METADATA`: 
#### Summary

If true, unknown properties in data handled by metadata service will cause an error
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.StrictJsonParsing.Metadata` |
| Environment variable: | `MEX_STRICTNESS_STRICT_JSON_PARSING_METADATA`  |
| Default value: | `'true'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_STRICTNESS_STRICT_JSON_PARSING_QUERY`: 
#### Summary

If true, unknown properties in data handled by query service will cause an error
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Strictness.StrictJsonParsing.Query` |
| Environment variable: | `MEX_STRICTNESS_STRICT_JSON_PARSING_QUERY`  |
| Default value: | `'true'` |
| Used by: | <ul><li>_all_</li></ul> |

----
### `MEX_NOTIFY_EMAILER_TYPE`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.EmailerType` |
| Environment variable: | `MEX_NOTIFY_EMAILER_TYPE`  |
| Default value: | `'MOCKMAILER'` |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_ORIGIN_OAUTH`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.OriginOauth` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_ORIGIN_OAUTH`  |
| Default value: | `'https://login.flowmailer.net'` |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_ORIGIN_API`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.OriginApi` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_ORIGIN_API`  |
| Default value: | `'https://api.flowmailer.net'` |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_CLIENT_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.ClientId` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_CLIENT_ID`  |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_CLIENT_SECRET`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.ClientSecret` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_CLIENT_SECRET`  |
| Secret: | **yes** |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_ACCOUNT_ID`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.AccountId` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_ACCOUNT_ID`  |
| Used by: | <ul><li>metadata</li></ul> |

----
### `MEX_NOTIFY_FLOWMAILER_NOREPLY_EMAIL_ADDRESS`: 
#### Info

| Key | Value |
| --- | ----- |
| Go struct field: | `.Notify.Flowmailer.NoreplyEmailAddress` |
| Environment variable: | `MEX_NOTIFY_FLOWMAILER_NOREPLY_EMAIL_ADDRESS`  |
| Default value: | `'noreply@data4life.care'` |
| Used by: | <ul><li>metadata</li></ul> |

----
