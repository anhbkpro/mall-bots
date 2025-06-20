# Generate proto files

## Configure buf

- Define `buf.yaml` and `buf.gen.yaml` files for each package

### `buf.yaml`

#### What?

- `buf.yaml` is a configuration file for buf, a Protocol Buffers (protobuf) linting and code generation tool.

#### Code explanation

```yaml
version: v1 # Specifies this is using version 1 of buf's configuration format.
lint: # linting rules to ensure code quality and consistency.
  enum_zero_value_suffix: _UNKNOWN # Enums should have a suffix of _UNKNOWN if they have a zero value.
  except: # Exclude certain rules from being applied to specific packages.
    - PACKAGE_VERSION_SUFFIX # Exclude the rule for package version suffix.
    - PACKAGE_DIRECTORY_MATCH # Exclude the rule for package directory match.
breaking: # Configures buf's ability to detect API-breaking changes between protobuf versions.
  use: # Specifies the strategy for breaking change detection.
    - FILE # `FILE` means it compares against files in the local workspace rather than a remote repository or registry.
```

#### Purpose

- This configuration ensures consistent protobuf code style while allowing some flexibility by exempting certain naming rules, and enables detection of breaking changes in your API definitions during development.

### `buf.gen.yaml`

#### What?

- `buf.gen.yaml` is a configuration file that defines how buf generates code from Protocol Buffer definitions

#### Code explanation

```yaml
version: v1 # Uses buf's v1 configuration format.
managed:
  enabled: true # Enables buf's managed mode, which automatically handles Go package paths and imports.
  go_package_prefix: # Configures Go package naming:
    default: eda-in-golang/baskets/basketspb # Sets the base Go module path for generated code
    except:
      - buf.build/googleapis/googleapis # Excludes Google APIs from this prefix (they use their own standard paths)
plugins: # Defines the plugins to use for code generation.
  - name: go # Generates Go structs and message types in the current directory with relative paths.
    out: . # Specifies the output directory for the generated code.
    opt:
      - paths=source_relative
  - name: go-grpc # Generates Go gRPC server and client code for service definitions.
    out: .
    opt:
      - paths=source_relative
  - name: grpc-gateway # Generates HTTP/REST gateway code that translates REST calls to gRPC, using API annotations from the specified YAML file.
    out: .
    opt:
      - paths=source_relative
      - grpc_api_configuration=internal/rest/api.annotations.yaml
  - name: openapiv2 # Generates OpenAPI/Swagger documentation in internal/rest/, merging all specs into a single api file.
    out: internal/rest
    opt:
      - grpc_api_configuration=internal/rest/api.annotations.yaml
      - openapi_configuration=internal/rest/api.openapi.yaml
      - allow_merge=true
      - merge_file_name=api
```

#### Purpose

This setup creates a complete Go microservice with:

- gRPC services for high-performance internal communication
- REST API gateway for external HTTP clients
- OpenAPI documentation for API consumers
- Proper Go module structure with consistent package naming

## Generate proto files

```bash
// we are in the baskets directory
cd baskets

// generate the proto files
buf generate
```
