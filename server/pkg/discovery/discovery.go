package discovery

import (
  "context"
  "fmt"
  "errors"
  "time"
  "math/rand"
)

type Registry interface {
  Register(ctx context.Context, instanceID, serviceName, hostPort string) error
  Deregister(ctx context.Context, instanceID, serviceName string) error
  ServiceAddresses(ctx context.Context, serviceName string) ([]string, error)
  ReportHealthyState(instanceID, serviceName string) error
}

var ErrNotFound = errors.New("no service addresses found")

func GenerateInstanceID(serviceName string) string {
  return fmt.Sprintf("%s-%d", serviceName,
    rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
  )
}