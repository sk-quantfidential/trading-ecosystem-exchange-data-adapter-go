package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/pkg/models"
)

// ServiceDiscoveryBehaviorTestSuite tests the behavior of service discovery operations
type ServiceDiscoveryBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestServiceDiscoveryBehaviorSuite runs the service discovery behavior test suite
func TestServiceDiscoveryBehaviorSuite(t *testing.T) {
	suite.Run(t, new(ServiceDiscoveryBehaviorTestSuite))
}

// TestServiceRegistration tests service registration and discovery
func (suite *ServiceDiscoveryBehaviorTestSuite) TestServiceRegistration() {
	var serviceID = GenerateTestID("service")

	suite.Given("a service to register", func() {
		// Service defined below
	}).When("registering the service", func() {
		service := suite.CreateTestServiceRegistration(serviceID, func(s *models.ServiceRegistration) {
			s.Name = "test-exchange-service"
			s.Version = "1.0.0"
			s.Status = "healthy"
		})

		err := suite.adapter.RegisterService(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).Then("the service should be discoverable", func() {
		services, err := suite.adapter.GetServicesByName(suite.ctx, "test-exchange-service")
		suite.Require().NoError(err)
		suite.GreaterOrEqual(len(services), 1)

		var found bool
		for _, svc := range services {
			if svc.ID == serviceID {
				found = true
				suite.Equal("1.0.0", svc.Version)
				suite.Equal("healthy", svc.Status)
				break
			}
		}
		suite.True(found, "Should find registered service")
	})
}

// TestServiceHeartbeat tests service heartbeat updates
func (suite *ServiceDiscoveryBehaviorTestSuite) TestServiceHeartbeat() {
	var serviceID = GenerateTestID("heartbeat-service")

	suite.Given("a registered service", func() {
		service := suite.CreateTestServiceRegistration(serviceID, func(s *models.ServiceRegistration) {
			s.Name = "heartbeat-test-service"
		})
		err := suite.adapter.RegisterService(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).When("sending a heartbeat", func() {
		err := suite.adapter.SendHeartbeat(suite.ctx, serviceID)
		suite.Require().NoError(err)
	}).Then("the service last heartbeat should be updated", func() {
		service, err := suite.adapter.GetService(suite.ctx, serviceID)
		suite.Require().NoError(err)
		suite.NotNil(service)
		suite.WithinDuration(time.Now(), service.LastHeartbeat, 5*time.Second)
	})
}
