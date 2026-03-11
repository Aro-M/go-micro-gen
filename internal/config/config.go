package config

// DBType represents the type of database to use.
type DBType string

const (
	DBPostgres DBType = "postgres"
	DBMongo    DBType = "mongo"
	DBNone     DBType = "none"
)

// BrokerType represents the message broker to use.
type BrokerType string

const (
	BrokerKafka    BrokerType = "kafka"
	BrokerRabbitMQ BrokerType = "rabbitmq"
	BrokerNATS     BrokerType = "nats"
	BrokerNone     BrokerType = "none"
)

// TransportType represents the primary service transport.
type TransportType string

const (
	TransportHTTP TransportType = "http"
	TransportGRPC TransportType = "grpc"
	TransportBoth TransportType = "both"
	TransportNone TransportType = "none"
)

// CloudProvider represents the cloud provider.
type CloudProvider string

const (
	CloudAWS  CloudProvider = "aws"
	CloudGCP  CloudProvider = "gcp"
	CloudNone CloudProvider = "none"
)

// ArchType represents the architecture pattern.
type ArchType string

const (
	ArchClean     ArchType = "clean"
	ArchHexagonal ArchType = "hexagonal"
	ArchDDD       ArchType = "ddd"
	ArchVertical  ArchType = "vertical"
	ArchStandard  ArchType = "standard"
)

// ArchFolders maps standard internal paths to architectural-specific paths.
type ArchFolders struct {
	Service    string
	Repository string
	Transport  string
	Domain     string
}

// GetArchFolders returns the directory mapping for the specified architecture.
func GetArchFolders(arch ArchType) ArchFolders {
	switch arch {
	case ArchDDD:
		return ArchFolders{
			Service:    "internal/application",
			Repository: "internal/infrastructure/persistence",
			Transport:  "internal/interfaces",
			Domain:     "internal/domain",
		}
	case ArchVertical:
		return ArchFolders{
			Service:    "internal/features/core",
			Repository: "internal/features/data",
			Transport:  "internal/features/web",
			Domain:     "internal/features/domain",
		}
	case ArchStandard:
		return ArchFolders{
			Service:    "internal/app",
			Repository: "internal/store",
			Transport:  "internal/api",
			Domain:     "internal/models",
		}
	default: // clean & hexagonal
		return ArchFolders{
			Service:    "internal/service",
			Repository: "internal/repository",
			Transport:  "internal/transport",
			Domain:     "internal/domain",
		}
	}
}

// CIType represents the CI/CD provider.
type CIType string

const (
	CIGitHub CIType = "github"
	CIGitLab CIType = "gitlab"
	CINone   CIType = "none"
)

// ServiceConfig holds all configuration for generating a new microservice.
type ServiceConfig struct {
	// ServiceName is the human-readable service name (e.g. "order-service")
	ServiceName string
	// ModulePath is the Go module path (e.g. "github.com/acme/order-service")
	ModulePath string
	// Architecture pattern to use
	Architecture ArchType
	// Architectural Folder overrides
	Arch ArchFolders
	// Database type
	Database DBType
	// Message broker type
	Broker BrokerType
	// Transport type (http, grpc, both)
	Transport TransportType
	// Whether to include Redis
	IncludeRedis bool
	// Whether to include GraphQL endpoint
	IncludeGraphQL bool
	// Whether to include JWT Auth Middleware
	IncludeJWT bool
	// Whether to include Serverless wrappers
	IncludeServerless bool
	// Whether to include DB mock seeder
	IncludeSeeding bool
	// Whether to include Grafana
	IncludeGrafana bool
	// Whether to include Docker assets
	IncludeDocker bool
	// Whether to include K8s
	IncludeK8s bool
	// Whether to include Helm
	IncludeHelm bool
	// Cloud provider
	Cloud CloudProvider
	// CI/CD provider
	CI CIType
	// OutputDir is the directory where the service will be generated
	OutputDir string
	// GoVersion to use in generated go.mod
	GoVersion string
}

// PackageName returns a safe Go package name derived from ServiceName.
func (c *ServiceConfig) PackageName() string {
	name := ""
	capitalize := false
	for i, ch := range c.ServiceName {
		if ch == '-' || ch == '_' {
			capitalize = true
			continue
		}
		if i == 0 || capitalize {
			name += string([]rune{ch - 32})
			capitalize = false
		} else {
			name += string(ch)
		}
	}
	return name
}
