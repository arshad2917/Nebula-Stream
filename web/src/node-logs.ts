export type NodeLogLevel = 'info' | 'warn' | 'error'

export type NodeLog = {
  id: string
  nodeId: string
  level: NodeLogLevel
  message: string
  at: number
}
