# o8n
A terminal UI for Operaton

```
         ____      
  ____  ( __ )____ 
 / __ \/ __  / __ \
/ /_/ / /_/ / / / /
\____/\____/_/ /_/ 
```

**o8n** is a powerful terminal UI for managing Operaton workflow engines, inspired by k9s.

## Quick Start

### 1. Configuration

There are three config files:

- `o8n-env.yaml` — Environment credentials and UI colors (keep secret, git-ignored)
- `o8n-cfg.yaml` — UI table definitions and app settings (version-controlled)
- `o8n-stat.yml` — Runtime state: active environment, skin, latency toggle, last view (auto-generated, git-ignored)

Create your environment configuration:

```bash
cp o8n-env.yaml.example o8n-env.yaml
# Edit o8n-env.yaml to add your Operaton environments
```

Example `o8n-env.yaml`:
```yaml
environments:
  local:
    url: http://localhost:8080/engine-rest
    username: demo
    password: demo
    ui_color: "#00A8E1"
    default_timeout: 10s
  production:
    url: https://operaton.example.com/engine-rest
    username: admin
    password: secret
    ui_color: "#FF5733"
```

> Note: `active` environment and `skin` are no longer stored in `o8n-env.yaml`. They are persisted in `o8n-stat.yml`.

### 2. Building

```bash
go build -o o8n .
```

### 3. Running

```bash
./o8n
```

## Usage

### Keyboard Shortcuts

**Global Actions:**
- `?` — Show help screen with all shortcuts
- `:` — Open context switcher (process-definition, process-instance, task, job, etc.)
- `<ctrl>+e` — Switch environment
- `<ctrl>+c` — Quit application

**Navigation:**
- `↑/↓` or `j/k` — Move selection up/down
- `Page Up/Down` — Jump through table
- `Enter` — Drill down (definitions → instances → variables)
- `Esc` — Go back one level

**View Actions:**
- `r` — Toggle auto-refresh (5s interval)
- `<ctrl>-r` — Manual refresh
- `L` — Toggle request latency display in footer (default: off)
- `/` — Filter/search (if implemented)

**Instance Actions:**
- `<ctrl>+d` — Delete selected instance (press twice to confirm)
- `<ctrl>+t` — Terminate instance

**Variable Actions:**
- `e` — Edit selected value (when column is editable)

### Features

**🎨 Theming & Skins**
- 20+ built-in color schemes (dracula, nord, gruvbox, solarized, etc.)
- Runtime skin switching
- Custom skin support in `/skins` folder

**📊 Dynamic Tables**
- Responsive column sizing based on terminal width
- Auto-hide low-priority columns on narrow terminals
- Customizable column visibility and widths in `o8n-cfg.yaml`
- Editable columns with type-aware input (opt-in)

**🔍 Context Switching**
- Fast context switching with `:` key
- Inline completion as you type
- Access all Operaton resources (process definitions, instances, tasks, jobs, incidents, etc.)

**⚡ Real-Time Updates**
- Auto-refresh mode with configurable intervals
- Visual indicator for API activity
- Error messages in footer with auto-clear

**🔒 Multi-Environment Support**
- Switch between environments with `<ctrl>+e`
- Environment-specific UI colors
- Secure credential management

**🎯 Drill-Down Navigation**
- Process Definition → Process Instances → Variables
- Breadcrumb navigation in footer
- Intuitive back navigation with `Esc`
- **View state restored on restart** — the app reopens at the last resource/drilldown level

**⚙️ Persistent State**
- Active environment, skin, and latency toggle are saved in `o8n-stat.yml`
- Last navigation position (resource type and drilldown path) is restored on startup
- Credentials stay stable in `o8n-env.yaml` (no runtime modifications)

### Debug Mode

Run the application with `--debug` to enable debug diagnostics and API access logging:

```bash
./o8n --debug
```

When enabled the application creates a `./debug` directory and writes two files:

- `./debug/last-screen.txt` — a dump of the last rendered TUI frame (useful for reproducing layout issues)
- `./debug/access.log` — an append-only log of API calls. Each API call is logged as two lines:

Example:

```
2026-02-17T02:11:25+01:00 API: FetchVariables instanceID="076899d8-0b54-11f1-b360-0242ac110002"
2026-02-17T02:11:25+01:00 API: GET /process-instance/{id}/variables
```

Query parameters are shown as a URL query string when present, e.g.:

```
2026-02-17T02:24:34+01:00 API: GET /process-instance?processDefinitionKey=invoice
```

The debug files are intended for troubleshooting and are safe to remove after use.

### Configuration Files

**o8n-env.yaml** (Environment Configuration):
```yaml
environments:
  <env-name>:
    url: <operaton-rest-api-url>
    username: <user>
    password: <password>
    ui_color: <hex-color>  # e.g., "#00A8E1"
    default_timeout: <duration> # e.g., "10s", "1m"
```

**o8n-stat.yml** (Runtime State — auto-generated, git-ignored):
```yaml
active_env: local
skin: dracula
show_latency: false
navigation:
  root: process-instance
  breadcrumb:
    - process-definitions
    - process-instances
  selected_definition_key: my-process
```

**o8n-cfg.yaml** (UI Configuration):
```yaml
tables:
  - name: process-definition
    columns:
      - name: key
        visible: true
        width: 20%
        align: left
      - name: name
        visible: true
        width: 40%
        align: left
      - name: version
        visible: true
        width: 15%
        align: center
      - name: resource
        visible: true
        width: 25%
        align: left
  - name: process-variables
    columns:
      - name: name
        visible: true
        width: 30%
        align: left
      - name: value
        visible: true
        width: 70%
        align: left
        editable: true
        input_type: auto
```

### Security Note

⚠️ **Important**: The environment file contains sensitive credentials.

- Add `o8n-env.yaml` to `.gitignore` (already configured)
- Never commit your actual `o8n-env.yaml` to version control
- Use appropriate file permissions: `chmod 600 o8n-env.yaml`
- Consider using environment variables for production deployments

## Development

### Prerequisites
- Go 1.24 or higher
- Docker (for API client generation)

### API Client Generation

Regenerate the Operaton REST API client:

```bash
./.devenv/scripts/generate-api-client.sh
```

### Testing

```bash
go test ./... -v
```

### Project Structure

```
o8n/
├── main.go              # Entry point (calls internal/app)
├── internal/
│   ├── app/             # TUI application logic (model, update, view)
│   ├── client/          # Operaton REST API client
│   ├── config/          # Config structs and loaders
│   └── ...
├── o8n-env.yaml         # Environment credentials (git-ignored)
├── o8n-cfg.yaml         # UI table definitions
├── o8n-stat.yml         # Runtime state (git-ignored, auto-generated)
├── resources/           # OpenAPI spec
├── skins/               # Color schemes
└── _bmad/core/prds/     # Design specifications
```

## Documentation

- [specification.md](specification.md) — Complete technical specification
- [Splash Screen Design](_bmad/core/prds/splash-screen-design.md)
- [Compact Header Design](_bmad/core/prds/compact-header-design.md)
- [Layout Design](_bmad/core/prds/layout-design-optimized.md)
- [Modal Confirmation Design](_bmad/core/prds/modal-confirmation-design.md)
- [Help Screen Design](_bmad/core/prds/help-screen-design.md)

## License

See [LICENSE](LICENSE) file.

## Contributing

Contributions welcome! Please read the specification.md for architecture details before submitting PRs.
