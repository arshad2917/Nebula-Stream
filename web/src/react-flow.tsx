'use client'

import { memo, useMemo } from 'react'
import ReactFlow, { Background, Controls, type Edge, type Node } from 'reactflow'
import 'reactflow/dist/style.css'

import type { EdgeMessage } from '@/message-flow'

export const MAX_VISIBLE_NODES = 100

const baseNodes: Node[] = [
  {
    id: 'trigger',
    position: { x: 40, y: 120 },
    data: { label: 'External Trigger' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
  {
    id: 'orchestrator',
    position: { x: 280, y: 120 },
    data: { label: 'Orchestrator Core' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
  {
    id: 'nats',
    position: { x: 550, y: 120 },
    data: { label: 'NATS Bus' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
  {
    id: 'wasm',
    position: { x: 820, y: 30 },
    data: { label: 'WASM Runtime' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
  {
    id: 'ai',
    position: { x: 820, y: 200 },
    data: { label: 'AI Worker' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
  {
    id: 'state',
    position: { x: 1080, y: 120 },
    data: { label: 'State / Callback' },
    style: { background: '#1f3a54', color: '#ecf4fd', border: '1px solid #4c7da8' },
  },
]

function mapMessagesToEdges(messages: EdgeMessage[]): Edge[] {
  return messages.map((message) => ({
    id: message.id,
    source: message.from,
    target: message.to,
    animated: true,
    style: { stroke: '#5bc0eb', strokeWidth: 2.6 },
  }))
}

export const WorkflowCanvas = memo(function WorkflowCanvas({ messages }: { messages: EdgeMessage[] }) {
  const edges = useMemo(() => mapMessagesToEdges(messages), [messages])

  return (
    <div style={{ width: '100%', height: 380 }}>
      <ReactFlow fitView nodes={baseNodes} edges={edges} nodesConnectable={false}>
        <Background color="#33556f" gap={30} />
        <Controls />
      </ReactFlow>
    </div>
  )
})
