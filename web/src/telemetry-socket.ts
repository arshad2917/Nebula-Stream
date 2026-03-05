'use client'

import { useEffect, useState } from 'react'

import type { EdgeMessage } from '@/message-flow'
import type { NodeLog, NodeLogLevel } from '@/node-logs'

type TelemetrySnapshot = {
  throughput: number
  activeNodes: number
  latencyMs: number
  logs: NodeLog[]
  messages: EdgeMessage[]
}

const pipelineEdges = [
  ['trigger', 'orchestrator'],
  ['orchestrator', 'nats'],
  ['nats', 'wasm'],
  ['nats', 'ai'],
  ['wasm', 'state'],
  ['ai', 'state'],
] as const

const logLevels: NodeLogLevel[] = ['info', 'warn', 'error']

function nextEdge(index: number): EdgeMessage {
  const [from, to] = pipelineEdges[index % pipelineEdges.length]
  const now = Date.now()
  return {
    id: `${from}-${to}-${now}`,
    from,
    to,
    at: now,
  }
}

function nextLog(index: number): NodeLog {
  const level = logLevels[index % logLevels.length]
  const now = Date.now()
  const nodeId = index % 2 === 0 ? 'wasm-node-a1' : 'ai-node-b2'
  const messageByLevel: Record<NodeLogLevel, string> = {
    info: 'step completed in expected latency budget',
    warn: 'queue depth increased, applying soft backpressure',
    error: 'transient transport error, retrying on secondary route',
  }

  return {
    id: `${nodeId}-${now}`,
    nodeId,
    level,
    message: messageByLevel[level],
    at: now,
  }
}

export function connectTelemetrySocket(onTick: (snapshot: TelemetrySnapshot) => void): () => void {
  let cursor = 0

  const interval = setInterval(() => {
    cursor += 1
    const snapshot: TelemetrySnapshot = {
      throughput: 46800 + (cursor % 12) * 310,
      activeNodes: 7 + (cursor % 3),
      latencyMs: 12 + (cursor % 8),
      logs: [nextLog(cursor)],
      messages: [nextEdge(cursor), nextEdge(cursor + 1)],
    }

    onTick(snapshot)
  }, 1250)

  return () => clearInterval(interval)
}

export function useTelemetryFeed() {
  const [data, setData] = useState<TelemetrySnapshot>({
    throughput: 47000,
    activeNodes: 8,
    latencyMs: 14,
    logs: [
      {
        id: 'boot-log-1',
        nodeId: 'orchestrator',
        level: 'info',
        message: 'telemetry stream initialized',
        at: Date.now(),
      },
    ],
    messages: [
      {
        id: 'boot-edge-1',
        from: 'trigger',
        to: 'orchestrator',
        at: Date.now(),
      },
      {
        id: 'boot-edge-2',
        from: 'orchestrator',
        to: 'nats',
        at: Date.now(),
      },
    ],
  })

  useEffect(() => {
    const unsubscribe = connectTelemetrySocket((snapshot) => {
      setData((prev) => ({
        throughput: snapshot.throughput,
        activeNodes: snapshot.activeNodes,
        latencyMs: snapshot.latencyMs,
        logs: [...snapshot.logs, ...prev.logs].slice(0, 16),
        messages: [...snapshot.messages, ...prev.messages].slice(0, 12),
      }))
    })

    return unsubscribe
  }, [])

  return data
}
