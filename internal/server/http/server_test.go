package http

import (
	"testing"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"

	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

// TestNewServer tests the creation of a new gRPC server with the provided dependencies
// This test checks if the server is not nil after creation
func TestNewServer(t *testing.T) {
	// Initialize mock database and Redis instances
	// These mocks are used to isolate dependencies and make the test more efficient
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	// Create a new gRPC server with the provided dependencies
	// The validation module is also provided to enable request validation
	server := NewServer(validation.New(), mockDB, mockRedis)

	// Assert that the server is not nil
	// This check ensures that the server was created successfully
	assert.NotNil(t, server)
}

// TestServer_GetEngine tests the GetEngine method of the gRPC server
// This test checks if the engine is not nil after calling the GetEngine method
func TestServer_GetEngine(t *testing.T) {
	// Initialize mock database and Redis instances
	// These mocks are used to isolate dependencies and make the test more efficient
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	// Create a new gRPC server with the provided dependencies
	// The validation module is also provided to enable request validation
	server := NewServer(validation.New(), mockDB, mockRedis)

	// Assert that the server is not nil
	// This check ensures that the server was created successfully
	assert.NotNil(t, server)

	// Call the GetEngine method of the gRPC server
	// This method returns the underlying engine used by the server
	engine := server.GetEngine()

	// Assert that the engine is not nil
	// This check ensures that the engine was retrieved successfully
	assert.NotNil(t, engine)
}

// TestServer_MapRoutes tests the MapRoutes method of the gRPC server
// This test checks if the error is nil after calling the MapRoutes method
func TestServer_MapRoutes(t *testing.T) {
	// Initialize mock database and Redis instances
	// These mocks are used to isolate dependencies and make the test more efficient
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	// Create a new gRPC server with the provided dependencies
	// The validation module is also provided to enable request validation
	server := NewServer(validation.New(), mockDB, mockRedis)

	// Assert that the server is not nil
	// This check ensures that the server was created successfully
	assert.NotNil(t, server)

	// Call the MapRoutes method of the gRPC server
	// This method maps the routes for the server
	err := server.MapRoutes()

	// Assert that the error is nil
	// This check ensures that the routes were mapped successfully
	assert.Nil(t, err)
}
