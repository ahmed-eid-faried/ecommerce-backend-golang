package grpc

import (
	"testing"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"

	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

// TestNewServer tests the creation of a new gRPC server with the provided dependencies
func TestNewServer(t *testing.T) {
	// Initialize mock database and Redis instances
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	// Create a new gRPC server with the provided dependencies
	server := NewServer(validation.New(), mockDB, mockRedis)

	// Assert that the server is not nil
	assert.NotNil(t, server)
}
