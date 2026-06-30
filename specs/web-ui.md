# Web UI

The Silmari admin interface — a **React + TypeScript** single-page application embedded in the Go control plane binary. Provides dashboards for monitoring workflows, managing agents, configuring triggers, and inspecting cryptographic audit trails.

**Root:** `control-plane/web/client/`

## Technology Stack

| Technology | Purpose |
|------------|---------|
| **React 18+** | UI framework (functional components + hooks) |
| **TypeScript** | Type-safe development |
| **Vite** | Build tool and dev server (HMR) |
| **Tailwind CSS** | Utility-first styling |
| **Radix UI** | Accessible headless UI primitives |
| **React Router** | Client-side routing |

**Code reference:** `control-plane/web/client/package.json` — full dependency list

## Architecture

```
┌─────────────────────────────────────────┐
│                 Pages                    │
│        src/pages/ (route components)     │
├─────────────────────────────────────────┤
│              Components                  │
│  ┌────────┐ ┌────────┐ ┌──────────────┐ │
│  │Dashboard│ │Workflow│ │Execution/    │ │
│  │        │ │  DAG   │ │Runs          │ │
│  └────────┘ └────────┘ └──────────────┘ │
│  ┌────────┐ ┌────────┐ ┌──────────────┐ │
│  │ Nodes  │ │Triggers│ │DID/VC        │ │
│  └────────┘ └────────┘ └──────────────┘ │
│  ┌────────┐ ┌────────┐ ┌──────────────┐ │
│  │Reasoner│ │ Forms  │ │Layout/Nav    │ │
│  │Catalog │ │        │ │              │ │
│  └────────┘ └────────┘ └──────────────┘ │
├─────────────────────────────────────────┤
│         State Management                 │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │ Contexts │ │  Hooks   │ │Queries  │  │
│  │(React    │ │(Custom   │ │(Data    │  │
│  │ Context) │ │ hooks)   │ │fetching)│  │
│  └──────────┘ └──────────┘ └─────────┘  │
├─────────────────────────────────────────┤
│         Services / Types / Utils         │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │ Services │ │  Types   │ │  Utils  │  │
│  │(API calls│ │(TS types)│ │(Helpers)│  │
│  └──────────┘ └──────────┘ └─────────┘  │
└─────────────────────────────────────────┘
```

## Directory Structure

```
src/
├── assets/           # Static assets (logos, images)
├── components/       # Reusable UI components by domain
│   ├── authorization/ # Auth and permission components
│   ├── dashboard/    # Dashboard widgets and cards
│   ├── did/          # DID identity management UI
│   ├── execution/    # Execution monitoring and details
│   ├── forms/        # Form components and validations
│   ├── layout/       # Page layout components
│   ├── Navigation/   # Navigation sidebar and menus
│   ├── nodes/        # Agent node management
│   ├── notes/        # Execution notes and annotations
│   ├── reasoners/    # Reasoner/skill catalog UI
│   ├── runs/         # Execution run history
│   ├── status/       # Status indicators and badges
│   ├── triggers/     # Trigger configuration UI
│   ├── ui/           # Generic UI primitives (Radix wrappers)
│   ├── vc/           # Verifiable Credential inspection
│   ├── workflow/     # Workflow management
│   ├── workflows/    # Workflow list and overview
│   └── WorkflowDAG/  # DAG visualization engine
│       ├── hooks/    # DAG-specific React hooks
│       ├── layouts/  # Graph layout algorithms
│       └── sections/ # DAG UI sections
├── config/           # Frontend configuration
├── contexts/         # React context providers
├── hooks/            # Custom React hooks
│   └── queries/      # Data fetching hooks (TanStack Query)
├── lib/              # Library utilities
├── pages/            # Route-level page components
├── services/         # API client services
├── styles/           # Global styles
├── types/            # TypeScript type definitions
├── utils/            # Utility functions
└── test/             # Test files mirroring src structure
    ├── components/   # Component tests
    ├── contexts/     # Context tests
    ├── hooks/        # Hook tests
    ├── lib/          # Utility tests
    ├── pages/        # Page tests
    ├── services/     # Service tests
    └── utils/        # Test utilities
```

## Component Architecture

### Workflow DAG (`components/WorkflowDAG/`)

The most complex component — renders interactive directed acyclic graphs of workflow executions:

- **hooks/** — Custom hooks for DAG state, selection, and interaction
- **layouts/** — Graph layout algorithms (Dagre, custom layouts)
- **sections/** — DAG sub-sections (node details, edge details, minimap)

The DAG visualizes parent-child relationships between agent executions, showing parallel branches, dependencies, and execution status (running, completed, failed, pending).

**Code reference:** `control-plane/web/client/src/components/WorkflowDAG/`

### Dashboard (`components/dashboard/`)

Overview widgets showing system health, active agents, recent workflows, and execution metrics.

**Code reference:** `control-plane/web/client/src/components/dashboard/`

### Nodes (`components/nodes/`)

Agent node lifecycle management — register, start, stop, view status, inspect logs.

**Code reference:** `control-plane/web/client/src/components/nodes/`

### Execution (`components/execution/`)

Execution monitoring — view active and historical executions, inspect inputs/outputs, pause/resume/cancel.

**Code reference:** `control-plane/web/client/src/components/execution/`

### Triggers (`components/triggers/`)

Configure external triggers — GitHub webhooks, Stripe events, Slack integrations, cron schedules.

**Code reference:** `control-plane/web/client/src/components/triggers/`

### DID/VC (`components/did/`, `components/vc/`)

Cryptographic identity inspection — view DIDs, verify credentials, inspect audit chains.

**Code reference:** `control-plane/web/client/src/components/did/`, `control-plane/web/client/src/components/vc/`

## State Management

State flows through three layers:

1. **React Context** (`src/contexts/`) — Global state providers (auth, theme, configuration)
2. **Custom Hooks** (`src/hooks/`) — Encapsulated state logic, shared across components
3. **Query Hooks** (`src/hooks/queries/`) — Server state via TanStack Query for data fetching, caching, and synchronization

API calls go through the services layer (`src/services/`) which wraps the control plane REST API.

**Code reference:** `control-plane/web/client/src/contexts/`, `control-plane/web/client/src/hooks/queries/`

## Routing

Client-side routing via React Router. Pages map to routes:

- `/` — Dashboard
- `/workflows` — Workflow list
- `/workflows/:id` — Workflow detail with DAG
- `/nodes` — Agent node management
- `/triggers` — Trigger configuration
- `/did` — DID/VC inspection

**Code reference:** `control-plane/web/client/src/pages/`

## Embedded Build Pipeline

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│ npm run     │────▶│ dist/        │────▶│ Go embed     │
│ build       │     │ (static      │     │ (embedded    │
│ (Vite)      │     │  HTML/JS/CSS)│     │  in binary)  │
└─────────────┘     └──────────────┘     └──────────────┘
```

In production, the Vite build output (`dist/`) is embedded into the Go binary using Go's `embed` package. In development, Vite runs a separate dev server at `http://localhost:5173` with hot module replacement, proxying API requests to the control plane.

**Code reference:** `control-plane/internal/embedded/` — Go embed configuration, `control-plane/web/client/vite.config.ts` — Vite configuration

## Development

```bash
cd control-plane/web/client
npm install
npm run dev    # Vite dev server at http://localhost:5173

# In parallel, run control plane:
cd control-plane
go run ./cmd/agentfield-server
```

The Vite dev server proxies `/api` requests to the control plane (default `http://localhost:8080`).

## Testing

```bash
npm run lint    # ESLint
npm test        # Vitest test runner
```

Tests mirror the source structure under `src/test/`. Component tests use React Testing Library.

**Code reference:** `control-plane/web/client/src/test/`
