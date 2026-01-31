package e2e

// Test configuration templates
const (
	// BasicConfig has dev and prod environments
	BasicConfig = `
[dev]
API_URL = "http://localhost:3000"
DB_HOST = "localhost"
DEBUG = "true"

[prod]
API_URL = "https://api.example.com"
DB_HOST = "prod-db.example.com"
DEBUG = "false"
`

	// NamespaceConfig has default and db namespaces
	NamespaceConfig = `
[dev]
API_URL = "http://localhost:3000"
ENV = "development"

[prod]
API_URL = "https://api.example.com"
ENV = "production"

[db.local]
DB_HOST = "localhost"
DB_PORT = "5432"
DB_NAME = "myapp_dev"

[db.prod]
DB_HOST = "prod-db.example.com"
DB_PORT = "5432"
DB_NAME = "myapp_prod"
`

	// MetadataConfig includes _web_url metadata
	MetadataConfig = `
[dev]
API_URL = "http://localhost:3000"
_web_url = "http://localhost:3000/admin"

[prod]
API_URL = "https://api.example.com"
_web_url = "https://api.example.com/admin"

[staging]
API_URL = "https://staging.example.com"
`

	// MultiNamespaceConfig has multiple namespaces
	MultiNamespaceConfig = `
[dev]
ENV = "development"

[prod]
ENV = "production"

[db.local]
DB_HOST = "localhost"

[db.prod]
DB_HOST = "prod-db.example.com"

[deploy.aws]
CLOUD = "aws"
REGION = "us-east-1"

[deploy.gcp]
CLOUD = "gcp"
REGION = "us-central1"
`

	// InvalidConfig has syntax errors
	InvalidConfig = `
[dev
API_URL = "broken
`

	// EmptyConfig has no environments
	EmptyConfig = ``
)

// LegacyState represents old state format (single config)
const LegacyState = `current_config = "dev"`

// NewState represents new state format (namespace-aware)
const NewStateDefault = `
[current]
"" = "dev"
`

const NewStateMulti = `
[current]
"" = "prod"
db = "local"
`
