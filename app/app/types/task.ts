
export type TaskStatus =
  | 'pending'
  | 'in_progress'
  | 'completed'

export type TaskPriority =
  | 'low'
  | 'medium'
  | 'high'

export interface Task {
  id: number
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  created_at: string
  updated_at: string
}

export interface TaskSummary {
  total: number
  pending: number
  in_progress: number
  completed: number
}


export interface CreateTaskInput {
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
}


export interface UpdateTaskInput {
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
}
