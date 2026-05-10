/**
 * AG-UI protocol helpers for AgentField TypeScript reasoners.
 *
 * Mirrors `sdk/python/agentfield/agui.py` 1:1 so a Node-side reasoner
 * has the same authoring surface as a Python one.
 *
 * Two ways to use this module:
 *
 *   1. **Buffered mode** — return a normal JSON response from a reasoner
 *      with the optional `toolCalls` / `state` / `stateDelta` /
 *      `reasoning` fields. The control plane translates those into
 *      AG-UI `TOOL_CALL_*` / `STATE_*` / `REASONING_*` events.
 *
 *   2. **Streaming mode** — return `Content-Type: application/x-ndjson`
 *      and stream chunks built with `textChunk()`, `reasoningChunk()`,
 *      `toolCallStartChunk()`, etc. Each chunk becomes one live AG-UI
 *      event (see `internal/handlers/agui_runs_streaming.go`).
 *
 * Reasoners reach the AG-UI / CopilotKit frontend via the control
 * plane's `POST /api/v1/agui/runs/<node>/<reasoner>` adapter.
 */
import type { ToolCallTrace, ToolCallRecord } from '../ai/ToolCalling.js';

export const STREAMING_CONTENT_TYPE = 'application/x-ndjson';

/** A single AG-UI tool-call entry (buffered-mode `toolCalls` array). */
export interface ToolCallEntry {
  id?: string;
  name: string;
  arguments: Record<string, unknown>;
  result?: unknown;
}

/** RFC 6902 patch op. */
export interface JsonPatchOp {
  op: 'replace' | 'add' | 'remove';
  path: string;
  value?: unknown;
}

/** A reasoning segment for buffered REASONING_* emission. */
export interface ReasoningSegment {
  content: string;
  id?: string;
}

// ---------------------------------------------------------------------------
// Buffered-mode helpers
// ---------------------------------------------------------------------------

/**
 * Build a single AG-UI tool-call entry. The control plane translates each
 * entry into a `TOOL_CALL_START` / `TOOL_CALL_ARGS` / `TOOL_CALL_END`
 * triad. When `result` is set (or `hasResult` is true), it also emits
 * `TOOL_CALL_RESULT` so a server-side trace renders in the UI.
 *
 * @param name Tool name. CopilotKit pattern-matches this against
 *   `useCopilotAction({name, render})` registrations.
 * @param args JSON-serializable arguments mapping.
 * @param opts.id Optional stable ID. If omitted, the control plane
 *   synthesizes one (works for one-shots; breaks correlation with
 *   follow-up tool messages).
 * @param opts.result Optional pre-executed result.
 * @param opts.hasResult Force `result: undefined` to be treated as an
 *   explicit "executed and returned null" instead of "not executed yet".
 */
export function toolCall(
  name: string,
  args?: Record<string, unknown>,
  opts: { id?: string; result?: unknown; hasResult?: boolean } = {},
): ToolCallEntry {
  const entry: ToolCallEntry = { name, arguments: { ...(args ?? {}) } };
  if (opts.id !== undefined) entry.id = opts.id;
  if (opts.result !== undefined || opts.hasResult) entry.result = opts.result;
  return entry;
}

/**
 * Convert a `ToolCallTrace` from `ctx.aiWithTools(...)` into the AG-UI
 * `toolCalls` list shape.
 */
export function toolCallsFromTrace(trace: ToolCallTrace | null | undefined): ToolCallEntry[] {
  if (!trace || !trace.calls?.length) return [];
  return trace.calls.map((rec, i) => recordToEntry(rec, i));
}

function recordToEntry(rec: ToolCallRecord, index: number): ToolCallEntry {
  const entry: ToolCallEntry = {
    id: `tc-trace-${index}`,
    name: rec.toolName,
    arguments: { ...(rec.arguments ?? {}) },
  };
  if (rec.error !== undefined && rec.error !== null) {
    entry.result = { error: rec.error };
  } else if (rec.result !== undefined && rec.result !== null) {
    entry.result = rec.result;
  }
  return entry;
}

/** Build a single RFC 6902 `replace` patch op for a `stateDelta` array. */
export function stateDeltaReplace(path: string, value: unknown): JsonPatchOp {
  if (!path.startsWith('/')) {
    throw new Error("RFC 6902 paths must start with '/'");
  }
  return { op: 'replace', path, value };
}

/**
 * Compute a minimal RFC 6902 patch list for top-level keys that differ
 * between `before` and `after`. Shallow only.
 */
export function stateDeltaFromDiff(
  before: Record<string, unknown>,
  after: Record<string, unknown>,
): JsonPatchOp[] {
  const ops: JsonPatchOp[] = [];
  const keys = new Set<string>([...Object.keys(before), ...Object.keys(after)]);
  for (const k of [...keys].sort()) {
    const path = `/${k}`;
    const inBefore = k in before;
    const inAfter = k in after;
    if (inBefore && inAfter) {
      if (!deepEqual(before[k], after[k])) ops.push({ op: 'replace', path, value: after[k] });
    } else if (inAfter) {
      ops.push({ op: 'add', path, value: after[k] });
    } else {
      ops.push({ op: 'remove', path });
    }
  }
  return ops;
}

function deepEqual(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  if (typeof a !== typeof b) return false;
  if (a === null || b === null) return a === b;
  if (typeof a !== 'object') return false;
  // JSON-serializable comparison is sufficient for shallow patch ops.
  try {
    return JSON.stringify(a) === JSON.stringify(b);
  } catch {
    return false;
  }
}

/**
 * Build a single REASONING_MESSAGE segment. Each segment becomes a
 * `REASONING_MESSAGE_START` / `_CONTENT` / `_END` triad inside a
 * `REASONING_START` / `_END` boundary.
 */
export function reasoningSegment(content: string, opts: { id?: string } = {}): ReasoningSegment {
  const out: ReasoningSegment = { content };
  if (opts.id !== undefined) out.id = opts.id;
  return out;
}

/**
 * Build a `reasoning` field value from a mix of strings and segments.
 *
 * @example
 *   return { result: text, reasoning: agui.reasoning('step 1', 'step 2') };
 */
export function reasoning(...segments: Array<string | ReasoningSegment>): Array<string | ReasoningSegment> {
  const out: Array<string | ReasoningSegment> = [];
  for (const s of segments) {
    if (typeof s === 'string') {
      if (s) out.push(s);
    } else if (s && typeof s === 'object' && typeof s.content === 'string') {
      out.push({ ...s });
    } else {
      throw new TypeError(`reasoning() segments must be string or {content,id?}; got ${typeof s}`);
    }
  }
  return out;
}

// ---------------------------------------------------------------------------
// Streaming chunk builders
//
// Each function returns a small object in the wire shape the control plane's
// streaming dispatcher consumes (see internal/handlers/agui_runs_streaming.go).
// The reasoner author yields these from an async generator; serializeStream
// turns each yield into one NDJSON line for the streaming response.
// ---------------------------------------------------------------------------

export type StreamingChunk = Record<string, unknown> & { type: string };

/** One chunk of assistant text. Concatenated client-side. */
export function textChunk(delta: string): StreamingChunk {
  return { type: 'text', delta };
}

/** One chunk of chain-of-thought, rendered in CopilotKit's "Thinking…" pane. */
export function reasoningChunk(delta: string): StreamingChunk {
  return { type: 'reasoning', delta };
}

/** Closes the current reasoning segment so the next reasoningChunk opens a fresh one. */
export function reasoningEndChunk(): StreamingChunk {
  return { type: 'reasoning_end' };
}

/**
 * Open a tool call. If you already have the full `arguments`, pass them
 * here and the dispatcher emits one `TOOL_CALL_ARGS` frame immediately;
 * otherwise stream them with `toolCallArgsChunk`.
 */
export function toolCallStartChunk(
  id: string,
  name: string,
  opts: { arguments?: Record<string, unknown>; parentMessageId?: string } = {},
): StreamingChunk {
  const out: StreamingChunk = { type: 'tool_call_start', id, name };
  if (opts.arguments !== undefined) out.arguments = { ...opts.arguments };
  if (opts.parentMessageId !== undefined) out.parentMessageId = opts.parentMessageId;
  return out;
}

/** One chunk of streaming tool-call arguments JSON. */
export function toolCallArgsChunk(id: string, delta: string): StreamingChunk {
  return { type: 'tool_call_args', id, delta };
}

/** Close a tool call. */
export function toolCallEndChunk(id: string): StreamingChunk {
  return { type: 'tool_call_end', id };
}

/** Server-side tool result — use after pre-executing the tool. */
export function toolCallResultChunk(
  id: string,
  content: string,
  opts: { role?: string } = {},
): StreamingChunk {
  return { type: 'tool_call_result', id, content, role: opts.role ?? 'tool' };
}

/** Full agent state snapshot (the value `useCoAgent({state})` reads). */
export function stateChunk(snapshot: unknown): StreamingChunk {
  return { type: 'state', snapshot };
}

/** RFC 6902 patch ops applied incrementally on top of the last snapshot. */
export function stateDeltaChunk(ops: JsonPatchOp[]): StreamingChunk {
  return { type: 'state_delta', ops: [...ops] };
}

/** Mark the start of a named step inside the run. */
export function stepStartedChunk(name: string): StreamingChunk {
  return { type: 'step_started', name };
}

/** Mark a step finished. */
export function stepFinishedChunk(name: string): StreamingChunk {
  return { type: 'step_finished', name };
}

/** Pass a foreign-system event through verbatim. */
export function rawChunk(event: unknown, opts: { source?: string } = {}): StreamingChunk {
  const out: StreamingChunk = { type: 'raw', event };
  if (opts.source !== undefined) out.source = opts.source;
  return out;
}

/** Application-defined event. Frontends subscribe by `name`. */
export function customChunk(name: string, value?: unknown): StreamingChunk {
  const out: StreamingChunk = { type: 'custom', name };
  if (value !== undefined) out.value = value;
  return out;
}

/**
 * Trailing buffered envelope — the dispatcher applies any
 * `toolCalls` / `state` / `stateDelta` / `reasoning` / `result` fields
 * here as if from a non-streaming reasoner.
 */
export function finalChunk(data: Record<string, unknown>): StreamingChunk {
  return { type: 'final', data: { ...data } };
}

/** Terminal error. The dispatcher emits RUN_ERROR and stops the run. */
export function errorChunk(message: string, opts: { code?: string } = {}): StreamingChunk {
  const out: StreamingChunk = { type: 'error', message };
  if (opts.code !== undefined) out.code = opts.code;
  return out;
}

// ---------------------------------------------------------------------------
// Streaming serialization
// ---------------------------------------------------------------------------

/**
 * Serialize an async iterable of chunk objects (or strings — strings are
 * wrapped as text chunks) into an async iterable of NDJSON-encoded
 * `Uint8Array`, suitable for any Node streaming response (Express,
 * Fastify, Hono, the built-in `http` module, or a Web `Response`
 * built from a `ReadableStream`).
 *
 * Express:
 *
 *   res.setHeader('Content-Type', agui.STREAMING_CONTENT_TYPE);
 *   for await (const buf of agui.serializeStream(chunks)) res.write(buf);
 *   res.end();
 *
 * Web `Response` (works in Node 20+, Hono, edge runtimes):
 *
 *   const body = new ReadableStream({
 *     async start(controller) {
 *       for await (const buf of agui.serializeStream(chunks)) controller.enqueue(buf);
 *       controller.close();
 *     }
 *   });
 *   return new Response(body, { headers: { 'Content-Type': agui.STREAMING_CONTENT_TYPE }});
 *
 * Bare strings yielded by the generator are auto-wrapped as text chunks
 * for ergonomics. Anything else must be a chunk object produced by one
 * of the chunk builders above (or a hand-rolled equivalent).
 */
export async function* serializeStream(
  source: AsyncIterable<StreamingChunk | string> | Iterable<StreamingChunk | string>,
): AsyncIterable<Uint8Array> {
  const encoder = new TextEncoder();
  for await (const item of source as AsyncIterable<StreamingChunk | string>) {
    let payload: StreamingChunk;
    if (typeof item === 'string') {
      payload = textChunk(item);
    } else if (item && typeof item === 'object') {
      payload = item;
    } else {
      throw new TypeError(
        `streaming reasoner yielded non-string/non-object value of type ${typeof item}; ` +
          'use one of the agui chunk builders',
      );
    }
    yield encoder.encode(JSON.stringify(payload) + '\n');
  }
}

// ---------------------------------------------------------------------------
// Harness relay
// ---------------------------------------------------------------------------

/**
 * Relay a `@anthropic-ai/claude-agent-sdk` async-iterable of messages
 * into AG-UI streaming chunks, message-by-message.
 *
 * Mirrors `relay_harness_stream` in the Python SDK. Recognized message
 * shapes (the dict form `HarnessResult.messages` records):
 *
 *   - `{ type:'assistant', message:{ content:[{type:'text', text:'...'}, ...] }}`
 *     → one `text` chunk per text block
 *   - `{ type:'assistant', message:{ content:[{type:'thinking', thinking:'...'}, ...] }}`
 *     → one `reasoning` chunk per thinking block
 *   - `{ type:'assistant', message:{ content:[{type:'tool_use', id:'...', name:'...', input:{...}}, ...] }}`
 *     → `tool_call_start` + `tool_call_end` per tool_use block
 *   - `{ type:'user', message:{ content:[{type:'tool_result', tool_use_id:'...', content:'...'}, ...] }}`
 *     → `tool_call_result` per tool_result block
 *   - `{ type:'result', subtype:'success', result:'...' }` →
 *     terminal — yields nothing (the dispatcher's stream-end logic wraps
 *     the run with MESSAGES_SNAPSHOT + RUN_FINISHED).
 *   - Anything unrecognized is wrapped as a `raw` chunk so the trace is
 *     preserved without inventing ad-hoc event types.
 *
 * Note: the Claude Agent SDK buffers per-message rather than per-token,
 * so this path streams at message granularity. True per-token streaming
 * requires the raw Anthropic streaming API, not the harness.
 */
export async function* relayHarnessStream(
  harnessIter: AsyncIterable<unknown> | Iterable<unknown>,
): AsyncIterable<StreamingChunk> {
  for await (const raw of harnessIter as AsyncIterable<unknown>) {
    let msg: Record<string, unknown>;
    if (raw && typeof raw === 'object' && !Array.isArray(raw)) {
      msg = raw as Record<string, unknown>;
    } else {
      yield rawChunk({ raw: String(raw) }, { source: 'harness' });
      continue;
    }

    const msgType = String(msg.type ?? '');
    if (msgType === 'result') {
      // Final aggregated text — dispatcher's stream-end synthesizes
      // MESSAGES_SNAPSHOT / RUN_FINISHED, so emit nothing here.
      continue;
    }
    if (msgType === 'system') {
      yield rawChunk(msg, { source: 'harness' });
      continue;
    }

    if (msgType === 'assistant' || msgType === 'user') {
      const content = harnessMessageContent(msg);
      if (content === undefined || content === null) {
        yield rawChunk(msg, { source: 'harness' });
        continue;
      }
      if (typeof content === 'string') {
        if (msgType === 'assistant' && content) yield textChunk(content);
        continue;
      }
      if (Array.isArray(content)) {
        for (const block of content) {
          if (!block || typeof block !== 'object') continue;
          const b = block as Record<string, unknown>;
          const btype = b.type;
          if (btype === 'text') {
            const text = String(b.text ?? '');
            if (text) yield textChunk(text);
          } else if (btype === 'thinking') {
            const thinking = String(b.thinking ?? '');
            if (thinking) yield reasoningChunk(thinking);
          } else if (btype === 'tool_use') {
            const tcid = String(b.id ?? '');
            const name = String(b.name ?? '');
            if (tcid && name) {
              const inp =
                b.input && typeof b.input === 'object' && !Array.isArray(b.input)
                  ? (b.input as Record<string, unknown>)
                  : {};
              yield toolCallStartChunk(tcid, name, { arguments: inp });
              yield toolCallEndChunk(tcid);
            }
          } else if (btype === 'tool_result') {
            const tcid = String(b.tool_use_id ?? '');
            if (tcid) {
              let inner = b.content;
              if (Array.isArray(inner)) {
                inner = (inner as unknown[])
                  .filter((x): x is Record<string, unknown> => !!x && typeof x === 'object')
                  .map((x) => String(x.text ?? ''))
                  .join('');
              } else if (typeof inner !== 'string') {
                inner = String(inner ?? '');
              }
              yield toolCallResultChunk(tcid, inner as string, { role: 'tool' });
            }
          } else {
            yield rawChunk(b, { source: 'harness' });
          }
        }
      }
      continue;
    }

    yield rawChunk(msg, { source: 'harness' });
  }
}

function harnessMessageContent(msg: Record<string, unknown>): unknown {
  if ('content' in msg) return msg.content;
  const inner = msg.message;
  if (inner && typeof inner === 'object' && !Array.isArray(inner)) {
    return (inner as Record<string, unknown>).content;
  }
  return undefined;
}
