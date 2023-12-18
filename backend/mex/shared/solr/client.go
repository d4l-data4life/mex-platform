package solr

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

// ClientAPI specifies the interface for a Solr client
type ClientAPI interface {
	DoRequest(ctx context.Context, method string, relativePath string, body []byte) (int, []byte, error)
	DoJSONQuery(ctx context.Context, urlParams url.Values, body *QueryBody) (*QueryResponse, int, error)

	GetSchemaUniqueKey(ctx context.Context) (string, error)

	GetSchemaFields(ctx context.Context) ([]FieldDef, error)
	AddSchemaFields(ctx context.Context, fields []FieldDef) error
	RemoveSchemaFields(ctx context.Context, fieldNames []string) error

	GetSchemaCopyFields(ctx context.Context) ([]CopyFieldResponse, error)
	AddSchemaCopyFields(ctx context.Context, fields []CopyFieldDef) error
	RemoveSchemaCopyFields(ctx context.Context, fields []RemoveCopyFieldSubBody) error

	GetSchemaDynamicFields(ctx context.Context) ([]DynamicFieldDef, error)
	AddSchemaDynamicFields(ctx context.Context, fieldDef []DynamicFieldDef) error
	RemoveSchemaDynamicFields(ctx context.Context, fieldNamePatterns []string) error

	GetCollections(ctx context.Context) ([]string, error)
	DeleteCollection(ctx context.Context, collectionName string) error
	CreateCollection(ctx context.Context, collectionName string, configsetName string, replicationFactor uint32) error

	DropIndex(ctx context.Context) error

	AddDocuments(ctx context.Context, docs []string, num int) error
	RemoveDocuments(ctx context.Context, docIDs []string) error

	GetClusterStatus(ctx context.Context) (*ClusterStatus, int, error)
	Ping(ctx context.Context) error
}

type SchemaFieldListResponse struct {
	Fields []FieldDef `json:"fields"`
}

type SchemaCopyFieldListResponse struct {
	CopyFields []CopyFieldResponse `json:"copyFields"`
}
type SchemaDynamicFieldListResponse struct {
	DynamicFields []DynamicFieldDef `json:"dynamicFields"`
}

type UniqueKeyResponse struct {
	ResponseHeader map[string]interface{} `json:"responseHeader"`
	UniqueKey      string                 `json:"uniqueKey"`
}

type FieldName struct {
	Name string `json:"name"`
}

type AddFieldBody struct {
	AddField []FieldDef `json:"add-field"`
}

type RemoveFieldBody struct {
	DeleteField []FieldName `json:"delete-field"`
}

type AddDynamicFieldBody struct {
	AddDynamicField []DynamicFieldDef `json:"add-dynamic-field"`
}

type RemoveDynamicFieldBody struct {
	DeleteDynamicField []FieldName `json:"delete-dynamic-field"`
}

type AddCopyFieldSubBody struct {
	Source string   `json:"source"`
	Dest   []string `json:"dest"`
}

type AddCopyFieldBody struct {
	AddCopyField []AddCopyFieldSubBody `json:"add-copy-field"`
}

type RemoveCopyFieldSubBody struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

type RemoveCopyFieldBody struct {
	DeleteCopyField []RemoveCopyFieldSubBody `json:"delete-copy-field"`
}

type AdminCollections struct {
	Collections []string `json:"collections"`
}

// Client represents the API of a specific Solr instance & collection
type solrClient struct {
	log L.Logger

	origin     string
	collection string

	basicAuth    string
	rootsCAs     *x509.CertPool
	batchSize    int
	commitWithin time.Duration
}

type ClientOption func(c *solrClient)

func NewClient(origin string, collection string, options ...ClientOption) ClientAPI {
	client := &solrClient{
		origin:       origin,
		collection:   collection,
		batchSize:    DefaultSolrBatchSize,
		commitWithin: DefaultSolrCommitTime,
	}

	for _, option := range options {
		option(client)
	}

	// If no logger was set, use a dummy one to avoid nil checks further down
	if client.log == nil {
		client.log = &L.NullLogger{}
	}

	return client
}

//revive:disable-next-line:unexported-return
func WithLogger(logger L.Logger) ClientOption {
	return func(c *solrClient) {
		c.log = logger
	}
}

//revive:disable-next-line:unexported-return
func WithBasicAuth(user string, password string) ClientOption {
	return func(c *solrClient) {
		if user != "" && password != "" {
			c.basicAuth = b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))
		}
	}
}

//revive:disable-next-line:unexported-return
func WithCertificates(certs *x509.CertPool) ClientOption {
	return func(c *solrClient) {
		c.rootsCAs = certs
	}
}

//revive:disable-next-line:unexported-return
func WithBatchSize(batchSize uint32) ClientOption {
	return func(c *solrClient) {
		c.batchSize = int(batchSize)
	}
}

//revive:disable-next-line:unexported-return
func WithCommitWithin(commitWithin time.Duration) ClientOption {
	return func(c *solrClient) {
		c.commitWithin = commitWithin
	}
}

func createError(errCode codes.Code, methodName string, failureDesc string, err error) error {
	errText := ""
	if err != nil {
		errText = err.Error()
	}
	return status.Error(errCode, fmt.Sprintf("%s - %s: %s", methodName, failureDesc, errText))
}

func (c *solrClient) DoRequest(ctx context.Context, method string, relativePath string, body []byte) (int, []byte, error) {
	if len(relativePath) == 0 {
		relativePath = "/"
	}

	if relativePath[0:1] != "/" {
		relativePath = "/" + relativePath
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.origin, relativePath), bytes.NewReader(body))
	if err != nil {
		return -1, nil, err
	}
	defer req.Body.Close()

	if c.basicAuth != "" {
		req.Header.Set("authorization", "Basic "+c.basicAuth)
	} else {
		c.log.Warn(ctx, L.Message("performing Solr request without authorization header"))
	}

	if method != "GET" {
		req.Header.Set("content-type", guessMIMEType(body))
	}

	client := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Second,
			//nolint:gosec // gosec complains about this issue: https://github.com/go-redis/redis/issues/1553
			TLSClientConfig: &tls.Config{
				RootCAs: c.rootsCAs,
			},
		},
	}
	defer client.CloseIdleConnections()

	start := time.Now()
	c.log.Trace(ctx, L.Messagef("%s %s", method, req.URL))
	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	c.log.Trace(ctx, L.Messagef("duration: %dms", time.Since(start).Milliseconds()))

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}

	return resp.StatusCode, data, nil
}

const (
	mimeApplicationOctetStream = "application/octet-stream"
	mimeApplicationJSON        = "application/json"
	mimeTextXML                = "text/xml"
)

// Quick and dirty MIME guesser.
func guessMIMEType(data []byte) string {
	if data == nil {
		return mimeApplicationOctetStream
	}

	if len(data) < 2 {
		return mimeApplicationOctetStream
	}

	prefix := string(data[0:1])
	switch {
	case prefix == "{":
		return mimeApplicationJSON
	case prefix == "[":
		return mimeApplicationJSON
	case prefix == "<":
		return mimeTextXML
	}

	return mimeApplicationOctetStream
}

// DoJSONQuery directly uses the JSON API of a Solr instance to carry out a query
func (c *solrClient) DoJSONQuery(ctx context.Context, urlParams url.Values, body *QueryBody) (*QueryResponse, int, error) {
	methodName := "DoJSONQuery"

	// Prepare query
	urlPath := fmt.Sprintf("/solr/%s/query", c.collection)
	if len(urlParams) > 0 {
		urlPath += fmt.Sprintf("?%v", urlParams.Encode())
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, createError(codes.InvalidArgument, methodName, "could not create JSON", err)
	}

	doStatus, responseBody, err := c.DoRequest(ctx, "POST", urlPath, marshalledBody)
	if err != nil {
		return nil, 0, err
	}

	// Parse and return result
	var queryResponse QueryResponse
	err = json.Unmarshal(responseBody, &queryResponse)
	if err != nil {
		return nil, 0, createError(codes.Internal, methodName, fmt.Sprintf("status: %d, failed to parse Solr response", doStatus), err)
	}

	return &queryResponse, doStatus, nil
}

// GetSchemaUniqueKey retrieves the name of the field set as unique ID in the Solr schema
func (c *solrClient) GetSchemaUniqueKey(ctx context.Context) (string, error) {
	methodName := "GetSchemaUniqueKey"

	statusCode, responseBody, err := c.DoRequest(ctx, "GET", fmt.Sprintf("/solr/%s/schema/uniquekey", c.collection), nil)
	if err != nil {
		return "", createError(codes.Internal, methodName, "could not perform GET", err)
	}

	if statusCode != http.StatusOK {
		return "", createError(codes.Internal, methodName, fmt.Sprintf("request to Solr did not succeed - status code: %d", statusCode), nil)
	}

	// Parse and return result
	var response UniqueKeyResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", createError(codes.Internal, methodName, "failed to parse Solr response", err)
	}

	return response.UniqueKey, nil
}

// GetSchemaFields retrieves all fields in the current Solr schema
func (c *solrClient) GetSchemaFields(ctx context.Context) ([]FieldDef, error) {
	methodName := "GetSchemaFields"

	statusCode, responseBody, err := c.DoRequest(ctx, "GET", fmt.Sprintf("/solr/%s/schema/fields", c.collection), nil)
	if err != nil {
		return nil, createError(codes.Internal, methodName, "could not perform GET", err)
	}

	if statusCode != http.StatusOK {
		return nil, createError(codes.Internal, methodName, fmt.Sprintf("request to Solr did not succeed - status code: %d", statusCode), nil)
	}

	// Parse and return result
	var response SchemaFieldListResponse
	decodeErr := json.Unmarshal(responseBody, &response)
	if decodeErr != nil {
		return nil, createError(codes.Internal, methodName, "failed to parse Solr response", decodeErr)
	}

	return response.Fields, nil
}

// AddSchemaFields adds fields to a Solr schema using the SchemaUpdates API
func (c *solrClient) AddSchemaFields(ctx context.Context, fields []FieldDef) error {
	c.log.Trace(ctx, L.Messagef("AddSchemaFields: %v", fields), L.Phase("solr-client"))

	methodName := "AddSchemaFields"
	marshalledBody, err := json.Marshal(AddFieldBody{AddField: fields})
	if err != nil {
		return createError(codes.Internal, methodName, "could not create JSON", err)
	}

	_, _, err = c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	if err != nil {
		return createError(codes.Internal, methodName, "could not add schema field", err)
	}

	return nil
}

// RemoveSchemaFields removes fields from a Solr schema using the SchemaUpdates API
func (c *solrClient) RemoveSchemaFields(ctx context.Context, fieldNames []string) error {
	c.log.Trace(ctx, L.Messagef("RemoveSchemaFields: %v", fieldNames), L.Phase("solr-client"))
	if len(fieldNames) == 0 {
		return nil
	}
	methodName := "RemoveSchemaFields"
	body := RemoveFieldBody{
		DeleteField: make([]FieldName, len(fieldNames)),
	}

	for i, fieldName := range fieldNames {
		body.DeleteField[i].Name = fieldName
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not create JSON", err)
	}

	_, _, err = c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not remove schema fields", err)
	}

	return nil
}

// GetSchemaCopyFields retrieves all fields in the current Solr schema
func (c *solrClient) GetSchemaCopyFields(ctx context.Context) ([]CopyFieldResponse, error) {
	methodName := "GetSchemaCopyFields"

	statusCode, responseBody, err := c.DoRequest(ctx, "GET", fmt.Sprintf("/solr/%s/schema/copyfields", c.collection), nil)
	if err != nil {
		return nil, createError(codes.Internal, methodName, "could not perform GET", err)
	}

	if statusCode != http.StatusOK {
		return nil, createError(codes.Internal, methodName, fmt.Sprintf("request to Solr did not succeed - status code: %d", statusCode), nil)
	}

	// Parse and return result
	var response SchemaCopyFieldListResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, createError(codes.Internal, methodName, "failed to parse Solr response", err)
	}

	return response.CopyFields, nil
}

// AddSchemaCopyFields adds copy fields to a Solr schema using the SchemaUpdates API
func (c *solrClient) AddSchemaCopyFields(ctx context.Context, fields []CopyFieldDef) error {
	c.log.Trace(ctx, L.Messagef("AddSchemaCopyFields: %v", fields), L.Phase("solr-client"))

	methodName := "AddSchemaCopyFields"
	body := AddCopyFieldBody{
		AddCopyField: make([]AddCopyFieldSubBody, len(fields)),
	}

	for i, f := range fields {
		body.AddCopyField[i] = AddCopyFieldSubBody{
			Source: f.Source,
			Dest:   f.Destination,
		}
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not create JSON", err)
	}

	_, _, err = c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not add schema copy field", err)
	}

	return nil
}

// RemoveSchemaCopyFields removes copy fields from a Solr schema using the SchemaUpdates API
func (c *solrClient) RemoveSchemaCopyFields(ctx context.Context, fields []RemoveCopyFieldSubBody) error {
	c.log.Trace(ctx, L.Messagef("RemoveSchemaCopyFields: %v", fields), L.Phase("solr-client"))
	if len(fields) == 0 {
		return nil
	}
	methodName := "RemoveSchemaCopyFields"
	body := RemoveCopyFieldBody{DeleteCopyField: fields}
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not create JSON", err)
	}

	_, _, err = c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not remove schema copy field", err)
	}

	return nil
}

// GetSchemaDynamicFields retrieves all fields in the current Solr schema
func (c *solrClient) GetSchemaDynamicFields(ctx context.Context) ([]DynamicFieldDef, error) {
	methodName := "GetSchemaDynamicFields"

	statusCode, responseBody, err := c.DoRequest(ctx, "GET", fmt.Sprintf("/solr/%s/schema/dynamicfields", c.collection), nil)
	if err != nil {
		return nil, createError(codes.Internal, methodName, "could not perform GET", err)
	}

	if statusCode != http.StatusOK {
		return nil, createError(codes.Internal, methodName, fmt.Sprintf("request to Solr did not succeed - status code: %d", statusCode), nil)
	}

	// Parse and return result
	var response SchemaDynamicFieldListResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, createError(codes.Internal, methodName, "failed to parse Solr response", err)
	}

	return response.DynamicFields, nil
}

// AddSchemaDynamicFields adds dynamic fields to a Solr schema using the SchemaUpdates API
func (c *solrClient) AddSchemaDynamicFields(ctx context.Context, dynamicFields []DynamicFieldDef) error {
	c.log.Trace(ctx, L.Messagef("AddSchemaDynamicFields: %v", dynamicFields), L.Phase("solr-client"))

	methodName := "AddSchemaDynamicFields"
	marshalledBody, marshallErr := json.Marshal(AddDynamicFieldBody{dynamicFields})
	if marshallErr != nil {
		return createError(codes.InvalidArgument, methodName, "could not create JSON", marshallErr)
	}

	_, _, err := c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	return err
}

// RemoveSchemaDynamicFields removes dynamic fields from a Solr schema using the SchemaUpdates API
func (c *solrClient) RemoveSchemaDynamicFields(ctx context.Context, fieldNamePatterns []string) error {
	c.log.Trace(ctx, L.Messagef("RemoveSchemaDynamicField: %v", fieldNamePatterns), L.Phase("solr-client"))
	if len(fieldNamePatterns) == 0 {
		return nil
	}
	methodName := "RemoveSchemaDynamicFields"
	body := RemoveDynamicFieldBody{
		DeleteDynamicField: make([]FieldName, len(fieldNamePatterns)),
	}

	for i, fieldName := range fieldNamePatterns {
		body.DeleteDynamicField[i].Name = fieldName
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not create JSON", err)
	}

	_, _, err = c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/schema", c.collection), marshalledBody)
	if err != nil {
		return createError(codes.InvalidArgument, methodName, "could not remove schema dynamic field", err)
	}

	return nil
}

func (c *solrClient) GetCollections(ctx context.Context) ([]string, error) {
	statusCode, responseBody, err := c.DoRequest(ctx, "GET", "/solr/admin/collections?action=LIST", nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("get admin collections status code: %d", statusCode)
	}

	var adminCollections AdminCollections
	err = json.Unmarshal(responseBody, &adminCollections)
	if err != nil {
		return nil, err
	}

	return adminCollections.Collections, nil
}

func (c *solrClient) CreateCollection(ctx context.Context, collectionName string, configsetName string, replicationFactor uint32) error {
	c.log.Trace(ctx, L.Messagef("CreateCollection: collection: '%s', configset: '%s'", collectionName, configsetName), L.Phase("solr-client"))

	//nolint:lll
	statusCode, body, err := c.DoRequest(ctx, "GET",
		fmt.Sprintf(`/solr/admin/collections?action=CREATE&autoAddReplicas=false&collection.configName=%s&maxShardsPerNode=1&name=%s&numShards=1&replicationFactor=%d&router.name=compositeId&wt=json`, configsetName, collectionName, replicationFactor), nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("collection creation failed, status code: %d (%s)", statusCode, string(body))
	}

	return nil
}

func (c *solrClient) DeleteCollection(ctx context.Context, collectionName string) error {
	c.log.Trace(ctx, L.Messagef("DeleteCollection: collection: '%s'", collectionName), L.Phase("solr-client"))

	statusCode, body, err := c.DoRequest(ctx, "GET",
		fmt.Sprintf(`/solr/admin/collections?action=DELETE&name=%s&wt=json`, collectionName), nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("collection deletion failed, status code: %d (%s)", statusCode, string(body))
	}

	return nil
}

func (c *solrClient) Ping(ctx context.Context) error {
	statusCode, _, err := c.DoRequest(ctx, "GET", "/solr/admin/info/system", nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("solr ping error: %d", statusCode)
	}

	return nil
}

// GetClusterStatus retrieves the SOlr cluster status and does a first parsing of it
func (c *solrClient) GetClusterStatus(ctx context.Context) (*ClusterStatus, int, error) {
	urlWIthCollection := fmt.Sprintf("/solr/admin/collections?action=CLUSTERSTATUS&collection=%s", c.collection)
	doStatus, responseBody, err := c.DoRequest(ctx, "GET", urlWIthCollection, nil)
	if err != nil {
		return nil, 0, err
	}

	// Parse and return result
	var clusterStatus ClusterStatus
	err = json.Unmarshal(responseBody, &clusterStatus)
	if err != nil {
		return nil, 0, createError(codes.Internal, "GetClusterStatus", fmt.Sprintf("status: %d, failed to parse Solr response", doStatus), err)
	}

	return &clusterStatus, doStatus, nil
}

// DropIndex deletes the current index completely
func (c *solrClient) DropIndex(ctx context.Context) error {
	c.log.Trace(ctx, L.Message("DropIndex"), L.Phase("solr-client"))

	deleteBody := `<delete><query>*:*</query></delete>`
	_, _, err := c.DoRequest(ctx, "POST", fmt.Sprintf("/solr/%s/update?commit=true", c.collection), []byte(deleteBody))
	return err
}

func (c *solrClient) AddDocuments(ctx context.Context, docs []string, num int) error {
	docsCount := num
	fullBatches := docsCount / c.batchSize
	lastBatchSize := docsCount % c.batchSize
	c.log.Info(ctx, L.Messagef("AddDocuments: %d (no. of full batches of size %d: %d - size of last non-full batch: %d)", docsCount, c.batchSize, fullBatches, lastBatchSize))

	// Load full batches
	for b := 0; b < fullBatches; b++ {
		batch := make([]string, c.batchSize)
		for i := 0; i < c.batchSize; i++ {
			batch[i] = docs[b*c.batchSize+i]
		}

		_, _, err := c.DoRequest(
			ctx, "POST",
			fmt.Sprintf("/solr/%s/update", c.collection),
			[]byte(fmt.Sprintf(`<add commitWithin="%d" overwrite="true">%s</add>`, c.commitWithin/time.Millisecond, strings.Join(batch, ""))),
		)
		if err != nil {
			c.log.Error(ctx, L.Messagef("Failed to index full document batch no. %d in Solr", b+1))
			return err
		}
	}

	if lastBatchSize > 0 {
		batch := make([]string, lastBatchSize)
		for i := 0; i < lastBatchSize; i++ {
			batch[i] = docs[docsCount-i-1]
		}
		_, _, err := c.DoRequest(
			ctx, "POST",
			fmt.Sprintf("/solr/%s/update", c.collection),
			[]byte(fmt.Sprintf(`<add commitWithin="%d" overwrite="true">%s</add>`, c.commitWithin/time.Millisecond, strings.Join(batch, ""))),
		)
		if err != nil {
			c.log.Error(ctx, L.Message("Failed to index leftover document batch in Solr"))
			return err
		}
	}

	return nil
}

func (c *solrClient) RemoveDocuments(ctx context.Context, docIDs []string) error {
	c.log.Trace(ctx, L.Messagef("RemoveDocuments: #: %d", len(docIDs)), L.Phase("solr-client"))
	if len(docIDs) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, id := range docIDs {
		sb.WriteString(fmt.Sprintf("<id>%s</id>", id))
	}

	_, _, err := c.DoRequest(
		ctx, "POST",
		fmt.Sprintf("/solr/%s/update", c.collection),
		[]byte(fmt.Sprintf(`<delete commitWithin="%d">%s</delete>`, c.commitWithin/time.Millisecond, sb.String())),
	)
	if err != nil {
		return err
	}
	return nil
}
