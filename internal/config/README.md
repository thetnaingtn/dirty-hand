# Atlas Configuration

Atlas uses [Viper](https://github.com/spf13/viper) for configuration management, providing flexible configuration options suitable for various deployment environments including Docker and Kubernetes.

## Configuration Sources

Configuration values are loaded in the following order of precedence (highest to lowest):

1. **Command Line Flags** (`--config` flag)
2. **Environment Variables** (without prefix)
3. **Configuration File** (YAML format)
4. **Default Values**

## Configuration File

### Default Location
By default, Atlas looks for `config.yaml` in the following locations:
- Current directory (`.`)
- `./internal/config/`
- `/etc/atlas/`
- `$HOME/.atlas/`

### Custom Configuration File
You can specify a custom configuration file using:

1. **Command Line Flag**: Use the `--config` flag (highest priority)
   ```bash
   ./atlas --config /path/to/your/config.yaml
   ./atlas --config /path/to/config/directory/
   ```

2. **Environment Variable**: Set `CONFIG_FILE` to the full path of your config file
   ```bash
   export CONFIG_FILE=/path/to/your/config.yaml
   ```

3. **Programmatically**: Pass the config file path when initializing the config
   ```go
   config, err := config.NewConfigWithPath("/path/to/config.yaml")
   ```

## Configuration Structure

```yaml
# Server configuration
server:
  addr: "localhost"  # Server address
  port: "8080"       # Server port

# Database configuration
database:
  dsn: ""            # Database data source name (DSN)

# Environment
environment: "development"  # development, staging, production
```

## Environment Variables

All configuration values can be overridden using environment variables:

### Server Configuration
- `SERVER_ADDR` or `ADDR`: Server address
- `SERVER_PORT` or `PORT`: Server port

### Database Configuration
- `DATABASE_DSN`: Database data source name

### General
- `ENVIRONMENT`: Application environment

## Container Environments

### Docker

For Docker deployments, use `internal/config/config.docker.yaml` as a template:

```yaml
server:
  addr: "0.0.0.0"  # Listen on all interfaces
  port: "8080"

database:
  dsn: "sqlite3://./data/atlas.db"

environment: "development"
```

**Docker Run Example:**
```bash
# Using environment variables
docker run -p 8080:8080 \
  -e DATABASE_DSN="postgresql://user:pass@host:5432/db" \
  atlas:latest

# Using mounted config file with --config flag
docker run -p 8080:8080 \
  -v /host/config.yaml:/etc/atlas/config.yaml \
  atlas:latest --config /etc/atlas/config.yaml

# Using CONFIG_FILE environment variable
docker run -p 8080:8080 \
  -v /host/config:/etc/atlas \
  -e CONFIG_FILE="/etc/atlas/production.yaml" \
  atlas:latest
```

### Kubernetes

For Kubernetes deployments, use ConfigMaps and Secrets:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: atlas-config
data:
  config.yaml: |
    server:
      addr: "0.0.0.0"
      port: "8080"
    environment: "production"
---
apiVersion: v1
kind: Secret
metadata:
  name: atlas-secrets
type: Opaque
stringData:
  DATABASE_DSN: "postgresql://user:password@postgres:5432/atlas"
```

**Deployment Example:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: atlas
spec:
  template:
    spec:
      containers:
      - name: atlas
        image: atlas:latest
        # Option 1: Use CONFIG_FILE environment variable
        env:
        - name: CONFIG_FILE
          value: "/etc/config/config.yaml"
        - name: DATABASE_DSN
          valueFrom:
            secretKeyRef:
              name: atlas-secrets
              key: DATABASE_DSN
        # Option 2: Use --config flag (modify command)
        command: ["./atlas"]
        args: ["--config", "/etc/config/config.yaml"]
        volumeMounts:
        - name: config
          mountPath: /etc/config
        ports:
        - containerPort: 8080
      volumes:
      - name: config
        configMap:
          name: atlas-config
```

## Configuration Methods

### Using Default Configuration (with --config flag support)
```go
// This automatically parses the --config flag if present
config, err := config.NewConfig()
if err != nil {
    log.Fatal(err)
}
```

### Using Custom Configuration File Path
```go
// Directly specify a config file path (bypasses flag parsing)
config, err := config.NewConfigWithPath("/path/to/config.yaml")
if err != nil {
    log.Fatal(err)
}
```

### Using External Flag Management
```go
// When you want to handle flags in your main application
import "flag"

var configPath = flag.String("config", "", "Path to configuration file")
flag.Parse()

config, err := config.NewConfigWithFlags(*configPath)
if err != nil {
    log.Fatal(err)
}
```

### Accessing Configuration Values
```go
// Get server address and port
addr := config.Server.Addr
port := config.Server.Port
serverAddr := fmt.Sprintf("%s:%s", config.Server.Addr, config.Server.Port)

// Get database DSN
dbDSN := config.Database.DSN

// Check environment
if config.IsProduction() {
    // Production-specific logic
}

if config.IsDevelopment() {
    // Development-specific logic
}
```

## Best Practices

1. **Use Environment Variables for Secrets**: Never store sensitive data like passwords or API keys in configuration files
2. **Container-Friendly Defaults**: Use `0.0.0.0` as the default address for containers
3. **Configuration Validation**: Validate required configuration values at startup
4. **Environment-Specific Configs**: Use different configuration files for different environments
5. **Twelve-Factor App**: Follow the [Twelve-Factor App](https://12factor.net/config) methodology for configuration
6. **Command-Line Flags**: Use `--config` flag for easy configuration file specification during development and deployment
7. **Flag Parsing**: Call `config.NewConfig()` after any other flag parsing in your application, or use `config.NewConfigWithFlags()` for external flag management
