import { describe, it, expect } from 'vitest';
import { agui } from '../src/index.js';
import type { ToolCallTrace } from '../src/ai/ToolCalling.js';

describe('agui — buffered helpers', () => {
  it('toolCall builds the canonical entry', () => {
    expect(agui.toolCall('showFlightCard', { from: 'SFO', to: 'JFK' })).toEqual({
      name: 'showFlightCard',
      arguments: { from: 'SFO', to: 'JFK' },
    });
  });

  it('toolCall handles empty arguments and explicit id', () => {
    expect(agui.toolCall('ping', undefined, { id: 'tc-1' })).toEqual({
      id: 'tc-1',
      name: 'ping',
      arguments: {},
    });
  });

  it('toolCall surfaces a result when provided', () => {
    expect(agui.toolCall('lookup', { q: 'x' }, { result: { ok: true } })).toEqual({
      name: 'lookup',
      arguments: { q: 'x' },
      result: { ok: true },
    });
  });

  it('toolCall hasResult forces a null result through', () => {
    expect(agui.toolCall('noop', undefined, { hasResult: true })).toEqual({
      name: 'noop',
      arguments: {},
      result: undefined,
    });
  });

  it('toolCallsFromTrace returns [] for empty / null traces', () => {
    expect(agui.toolCallsFromTrace(null)).toEqual([]);
    expect(agui.toolCallsFromTrace(undefined)).toEqual([]);
    const empty: ToolCallTrace = { calls: [], totalTurns: 0, totalToolCalls: 0 };
    expect(agui.toolCallsFromTrace(empty)).toEqual([]);
  });

  it('toolCallsFromTrace converts records, surfaces result and error', () => {
    const trace: ToolCallTrace = {
      totalTurns: 1,
      totalToolCalls: 2,
      calls: [
        { toolName: 'a', arguments: { x: 1 }, result: { ok: true }, latencyMs: 5, turn: 0 },
        { toolName: 'b', arguments: {}, error: 'boom', latencyMs: 5, turn: 0 },
      ],
    };
    expect(agui.toolCallsFromTrace(trace)).toEqual([
      { id: 'tc-trace-0', name: 'a', arguments: { x: 1 }, result: { ok: true } },
      { id: 'tc-trace-1', name: 'b', arguments: {}, result: { error: 'boom' } },
    ]);
  });

  it('stateDeltaReplace emits a JSON Patch op', () => {
    expect(agui.stateDeltaReplace('/counter', 2)).toEqual({
      op: 'replace',
      path: '/counter',
      value: 2,
    });
  });

  it('stateDeltaReplace rejects paths missing a leading slash', () => {
    expect(() => agui.stateDeltaReplace('counter', 1)).toThrow(/RFC 6902/);
  });

  it('stateDeltaFromDiff emits a minimal shallow patch', () => {
    const before = { a: 1, b: 2, c: 3 };
    const after = { a: 1, b: 99, d: 4 };
    expect(agui.stateDeltaFromDiff(before, after)).toEqual([
      { op: 'replace', path: '/b', value: 99 },
      { op: 'remove', path: '/c' },
      { op: 'add', path: '/d', value: 4 },
    ]);
  });

  it('reasoningSegment + reasoning() build the segment list', () => {
    const seg = agui.reasoningSegment('thinking', { id: 'r1' });
    expect(seg).toEqual({ content: 'thinking', id: 'r1' });
    expect(agui.reasoning('a', '', seg, 'b')).toEqual([
      'a',
      { content: 'thinking', id: 'r1' },
      'b',
    ]);
  });

  it('reasoning() rejects garbage segments', () => {
    expect(() => agui.reasoning(42 as unknown as string)).toThrow(/segments must be string/);
  });
});

describe('agui — streaming chunk builders', () => {
  it('text/reasoning/reasoning_end', () => {
    expect(agui.textChunk('hi')).toEqual({ type: 'text', delta: 'hi' });
    expect(agui.reasoningChunk('thinking')).toEqual({ type: 'reasoning', delta: 'thinking' });
    expect(agui.reasoningEndChunk()).toEqual({ type: 'reasoning_end' });
  });

  it('toolCallStart with and without args/parent', () => {
    expect(agui.toolCallStartChunk('tc1', 'foo')).toEqual({
      type: 'tool_call_start',
      id: 'tc1',
      name: 'foo',
    });
    expect(
      agui.toolCallStartChunk('tc2', 'bar', { arguments: { x: 1 }, parentMessageId: 'm1' }),
    ).toEqual({
      type: 'tool_call_start',
      id: 'tc2',
      name: 'bar',
      arguments: { x: 1 },
      parentMessageId: 'm1',
    });
  });

  it('toolCallArgs / toolCallEnd / toolCallResult', () => {
    expect(agui.toolCallArgsChunk('tc1', '{"x":')).toEqual({
      type: 'tool_call_args',
      id: 'tc1',
      delta: '{"x":',
    });
    expect(agui.toolCallEndChunk('tc1')).toEqual({ type: 'tool_call_end', id: 'tc1' });
    expect(agui.toolCallResultChunk('tc1', 'done')).toEqual({
      type: 'tool_call_result',
      id: 'tc1',
      content: 'done',
      role: 'tool',
    });
    expect(agui.toolCallResultChunk('tc1', 'done', { role: 'system' })).toMatchObject({
      role: 'system',
    });
  });

  it('state / state_delta', () => {
    expect(agui.stateChunk({ k: 1 })).toEqual({ type: 'state', snapshot: { k: 1 } });
    const ops = [agui.stateDeltaReplace('/k', 2)];
    expect(agui.stateDeltaChunk(ops)).toEqual({ type: 'state_delta', ops });
  });

  it('step_started / step_finished', () => {
    expect(agui.stepStartedChunk('plan')).toEqual({ type: 'step_started', name: 'plan' });
    expect(agui.stepFinishedChunk('plan')).toEqual({ type: 'step_finished', name: 'plan' });
  });

  it('raw / custom / final / error chunk shapes', () => {
    expect(agui.rawChunk({ k: 1 })).toEqual({ type: 'raw', event: { k: 1 } });
    expect(agui.rawChunk({ k: 1 }, { source: 'harness' })).toEqual({
      type: 'raw',
      event: { k: 1 },
      source: 'harness',
    });
    expect(agui.customChunk('progress', 0.5)).toEqual({
      type: 'custom',
      name: 'progress',
      value: 0.5,
    });
    expect(agui.customChunk('ping')).toEqual({ type: 'custom', name: 'ping' });
    expect(agui.finalChunk({ result: 'done' })).toEqual({
      type: 'final',
      data: { result: 'done' },
    });
    expect(agui.errorChunk('boom')).toEqual({ type: 'error', message: 'boom' });
    expect(agui.errorChunk('boom', { code: 'E_BOOM' })).toEqual({
      type: 'error',
      message: 'boom',
      code: 'E_BOOM',
    });
  });
});

describe('agui — serializeStream', () => {
  it('emits one NDJSON line per chunk', async () => {
    async function* chunks() {
      yield agui.textChunk('a');
      yield agui.textChunk('b');
      yield agui.toolCallEndChunk('tc1');
    }
    const decoder = new TextDecoder();
    const lines: string[] = [];
    for await (const buf of agui.serializeStream(chunks())) {
      lines.push(decoder.decode(buf));
    }
    expect(lines).toEqual([
      JSON.stringify({ type: 'text', delta: 'a' }) + '\n',
      JSON.stringify({ type: 'text', delta: 'b' }) + '\n',
      JSON.stringify({ type: 'tool_call_end', id: 'tc1' }) + '\n',
    ]);
  });

  it('auto-wraps bare strings as text chunks', async () => {
    async function* gen() {
      yield 'hello';
      yield agui.textChunk(' world');
    }
    const decoder = new TextDecoder();
    let combined = '';
    for await (const buf of agui.serializeStream(gen())) combined += decoder.decode(buf);
    expect(combined.trim().split('\n').map((l) => JSON.parse(l))).toEqual([
      { type: 'text', delta: 'hello' },
      { type: 'text', delta: ' world' },
    ]);
  });

  it('rejects non-string non-object values', async () => {
    async function* gen() {
      yield 42 as unknown as string;
    }
    await expect(async () => {
      for await (const _ of agui.serializeStream(gen())) {
        /* drain */
      }
    }).rejects.toThrow(/non-string\/non-object/);
  });

  it('accepts a synchronous iterable too', async () => {
    const chunks = [agui.textChunk('x'), agui.textChunk('y')];
    const decoder = new TextDecoder();
    let n = 0;
    for await (const buf of agui.serializeStream(chunks)) {
      const obj = JSON.parse(decoder.decode(buf));
      expect(obj.type).toBe('text');
      n++;
    }
    expect(n).toBe(2);
  });
});

describe('agui — relayHarnessStream', () => {
  async function* fromArray(items: unknown[]) {
    for (const x of items) yield x;
  }

  it('translates assistant text blocks into text chunks', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        {
          type: 'assistant',
          message: {
            content: [
              { type: 'text', text: 'hello ' },
              { type: 'text', text: 'world' },
            ],
          },
        },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'text', delta: 'hello ' },
      { type: 'text', delta: 'world' },
    ]);
  });

  it('translates assistant thinking blocks into reasoning chunks', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        {
          type: 'assistant',
          message: { content: [{ type: 'thinking', thinking: 'hmm' }] },
        },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([{ type: 'reasoning', delta: 'hmm' }]);
  });

  it('translates tool_use blocks into start+end pairs', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        {
          type: 'assistant',
          message: {
            content: [{ type: 'tool_use', id: 'tc1', name: 'lookup', input: { q: 'x' } }],
          },
        },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'tool_call_start', id: 'tc1', name: 'lookup', arguments: { q: 'x' } },
      { type: 'tool_call_end', id: 'tc1' },
    ]);
  });

  it('translates tool_result string content', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        {
          type: 'user',
          message: {
            content: [{ type: 'tool_result', tool_use_id: 'tc1', content: 'ok' }],
          },
        },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'tool_call_result', id: 'tc1', content: 'ok', role: 'tool' },
    ]);
  });

  it('translates tool_result list content by stitching text blocks', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        {
          type: 'user',
          message: {
            content: [
              {
                type: 'tool_result',
                tool_use_id: 'tc1',
                content: [
                  { type: 'text', text: 'a' },
                  { type: 'text', text: 'b' },
                ],
              },
            ],
          },
        },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'tool_call_result', id: 'tc1', content: 'ab', role: 'tool' },
    ]);
  });

  it('skips terminal result envelope and surfaces system as raw', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        { type: 'system', subtype: 'init' },
        { type: 'result', subtype: 'success', result: 'done' },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'raw', event: { type: 'system', subtype: 'init' }, source: 'harness' },
    ]);
  });

  it('preserves unknown blocks and unknown top-level messages as raw', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        { type: 'assistant', message: { content: [{ type: 'mystery', payload: 1 }] } },
        { type: 'no-such-thing' },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks[0]).toMatchObject({ type: 'raw', source: 'harness' });
    expect((chunks[0] as { event: { type: string } }).event.type).toBe('mystery');
    expect(chunks[1]).toMatchObject({ type: 'raw', source: 'harness' });
  });

  it('handles bare content and string content shapes', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(
      fromArray([
        { type: 'assistant', content: 'inline-string' },
        { type: 'assistant', content: [{ type: 'text', text: 'inline-list' }] },
      ]),
    )) {
      chunks.push(ch);
    }
    expect(chunks).toEqual([
      { type: 'text', delta: 'inline-string' },
      { type: 'text', delta: 'inline-list' },
    ]);
  });

  it('wraps non-object iterates as raw', async () => {
    const chunks = [];
    for await (const ch of agui.relayHarnessStream(fromArray(['scalar', 7, null]))) {
      chunks.push(ch);
    }
    expect(chunks.every((c) => (c as { type: string }).type === 'raw')).toBe(true);
    expect(chunks).toHaveLength(3);
  });
});

describe('agui — STREAMING_CONTENT_TYPE', () => {
  it('matches the wire constant from the Python and Go SDKs', () => {
    expect(agui.STREAMING_CONTENT_TYPE).toBe('application/x-ndjson');
  });
});
