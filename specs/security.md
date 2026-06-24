# Security

AgentField's security model provides cryptographic identity, verifiable audit trails, and access control for multi-agent systems. The implementation uses **W3C Decentralized Identifiers (DIDs)** and **Verifiable Credentials (VCs)** under the hood.

## Architecture

```
┌─────────────────────────────────────────┐
│           Identity Layer                 │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │   DIDs   │ │   VCs    │ │  Keys   │  │
│  │(Decentral│ │(Verifiab.│ │(Ed25519 │  │
│  │Identifiers│ │Credent. │ │  + secp)│  │
│  └──────────┘ └──────────┘ └─────────┘  │
├─────────────────────────────────────────┤
│           Access Control                 │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │  Agent   │ │  Tool    │ │ Policy  │  │
│  │  IAM     │ │  Access  │ │ Engine  │  │
│  └──────────┘ └──────────┘ └─────────┘  │
├─────────────────────────────────────────┤
│           Audit Layer                    │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │Execution │ │  Signed  │ │  Chain  │  │
│  │ Records  │ │ Receipts │ │Verify   │  │
│  └──────────┘ └──────────┘ └─────────┘  │
├─────────────────────────────────────────┤
│           Cryptographic Primitives       │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐  │
│  │Encryption│ │ Signing  │ │ Hashing │  │
│  │(internal │ │(Ed25519) │ │(SHA-256│  │
│  │/encrypt) │ │          │ │ + BLAKE)│  │
│  └──────────┘ └──────────┘ └─────────┘  │
└─────────────────────────────────────────┘
```

## Decentralized Identity (DID)

AgentField implements the W3C DID standard for agent identity. Each agent can have a unique DID that serves as its cryptographic identity in the system.

### DID System Components

| Component | Control Plane | Python SDK | Go SDK |
|-----------|--------------|------------|--------|
| DID Creation | `internal/encryption/encryption.go` | `did_manager.py` | `did/did_manager.go` |
| Key Management | `internal/encryption/encryption.go` | `did_manager.py` | `did/did_manager.go` |
| DID Resolution | Services layer | `did_manager.py` | `did/did_client.go` |

### Configuration

DID is opt-in per agent:

**Python:**
```python
app.vc_generator.set_enabled(True)
```

**Go:**
```go
cfg.EnableDID = true   // or cfg.VCEnabled = true
```

When enabled, the agent automatically initializes its DID during `Initialize()` (Go) or startup (Python).

**Code reference:** `control-plane/internal/encryption/encryption.go` — crypto primitives, `sdk/go/agent/agent_did.go` — Go DID auto-init, `sdk/python/agentfield/did_manager.py` — Python DID management

## Verifiable Credentials (VCs)

AgentField generates W3C Verifiable Credentials for agent executions, creating a cryptographically verifiable record of what ran, who ran it, and how a workflow evolved.

### VC Generation

1. Agent execution completes → control plane generates VC
2. VC includes: agent DID, execution input hash, output hash, timestamp, workflow context
3. VC is signed with agent's private key
4. VC is stored in execution record
5. VCs form a chain across workflow steps

### VC Chain Export

```bash
# Export audit trail for a workflow
GET /api/v1/did/workflow/{workflow_id}/vc-chain

# Verify offline
af verify audit.json
```

**Code reference:** `sdk/python/agentfield/vc_generator.py` — Python VC generation, `sdk/go/did/vc_generator.go` — Go VC generation, `control-plane/internal/handlers/verify_audit.go` — audit verification handler

## Cryptographic Audit Trails

The biggest difference between an agent app and an AI backend: the backend needs to prove what happened.

### What's Recorded

- **Who ran:** Agent DID and node identity
- **What ran:** Reasoner/skill name, input hash, output hash
- **How it evolved:** Full workflow DAG with parent-child execution links
- **When:** Cryptographic timestamps
- **Decisions:** Human approvals, model decisions, routing choices

### Audit Chain Structure

```
┌─────────────────┐
│ Workflow Root   │
│ DID: did:af:wf1 │
│ Timestamp: t0   │
└────────┬────────┘
         │
┌────────▼────────┐
│ Execution A     │
│ DID: did:af:agA │
│ Input hash: h1  │
│ Output hash: h2 │── signed by agA's key
│ Parent: wf1     │
└────────┬────────┘
         │
┌────────▼────────┐
│ Execution B     │
│ DID: did:af:agB │
│ Input hash: h3  │
│ Output hash: h4 │── signed by agB's key
│ Parent: execA   │
└─────────────────┘
```

Each execution links to its parent, forming an append-only cryptographic chain. Any tampering breaks the chain.

**Code reference:** `control-plane/internal/handlers/verify_audit.go` — audit verification, `control-plane/internal/encryption/encryption.go` — signing primitives

## Access Control & IAM

AgentField implements first-class IAM for AI backends:

### Model

- **Callers** may be another agent, a human user, or an external system
- **Targets** may be an agent reasoner, a tool, or a memory scope
- **Decisions** may need to be proven later (hence DID/VC integration)

### Access Policies

Access policies control which agents can call which reasoners and access which memory scopes:

- Agent-level: "Agent A can call Agent B"
- Reasoner-level: "Agent A can call Agent B.score_claim but not Agent B.delete_data"
- Memory-level: "Agent A can read from global scope but not write"
- Tool-level: "Agent A can use tool X but not tool Y"

### Implementation Layers

| Layer | Purpose | Code Reference |
|-------|---------|---------------|
| Handler middleware | Auth token validation, agent identity extraction | `internal/server/middleware/` |
| Service authorization | Policy evaluation for agent-to-agent calls | `internal/services/` |
| DID-based auth | Cryptographic proof of caller identity | `internal/encryption/encryption.go` |

**Code reference:** `control-plane/internal/server/middleware/` — auth middleware, `control-plane/web/client/src/components/authorization/` — UI for policy management

## Outbound API Identity

Agents calling external APIs need identity too. AgentField provides outbound API identity so external services can verify which agent made a call:

- Agents can present DIDs when calling external services
- External services can verify agent identity via control plane's DID resolution endpoint
- Signed requests provide non-repudiation for agent actions

## Cryptographic Primitives

`control-plane/internal/encryption/encryption.go` provides:

- **Key generation:** Ed25519 and secp256k1 key pairs
- **Signing:** Ed25519 digital signatures for VCs and audit records
- **Verification:** Signature verification for audit chain validation
- **Hashing:** SHA-256 and BLAKE3 for content addressing

## Security Considerations

1. **Private keys** are held by agents, not the control plane — the control plane cannot forge agent signatures
2. **Audit chains** are append-only — execution records are immutable once signed
3. **DID resolution** is centralized through the control plane (not a public ledger), optimized for operational use
4. **VC verification** can be done offline with exported audit chains
5. **Access policies** are evaluated at the control plane, which is the trust boundary
