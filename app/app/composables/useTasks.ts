
import type {
  Task,
  TaskSummary,
  CreateTaskInput,
  UpdateTaskInput,
  TaskStatus,
} from '~/types/task'

export const useTasks = () => {
  const { apiBase } = useApi()

  // Obtener todas las tareas, con filtros opcionales.
  const getTasks = async (filters?: {
    status?: TaskStatus
    search?: string
  }): Promise<Task[]> => {
    return await $fetch<Task[]>(`${apiBase}/api/v1/tasks`, {
      query: filters,
    })
  }

  // Obtener una tarea por su identificador.
  const getTask = async (id: number): Promise<Task> => {
    return await $fetch<Task>(`${apiBase}/api/v1/tasks/${id}`)
  }

  // Obtener el resumen de las tareas.
  const getSummary = async (): Promise<TaskSummary> => {
    return await $fetch<TaskSummary>(
      `${apiBase}/api/v1/tasks/summary`,
    )
  }

  // Crear una nueva tarea.
  const createTask = async (
    input: CreateTaskInput,
  ): Promise<Task> => {
    return await $fetch<Task>(`${apiBase}/api/v1/tasks`, {
      method: 'POST',
      body: input,
    })
  }

  // Actualizar una tarea existente.
  const updateTask = async (
    id: number,
    input: UpdateTaskInput,
  ): Promise<Task> => {
    return await $fetch<Task>(`${apiBase}/api/v1/tasks/${id}`, {
      method: 'PUT',
      body: input,
    })
  }

  // Eliminar una tarea.
  const deleteTask = async (id: number): Promise<void> => {
    await $fetch<void>(`${apiBase}/api/v1/tasks/${id}`, {
      method: 'DELETE',
    })
  }

  return {
    getTasks,
    getTask,
    getSummary,
    createTask,
    updateTask,
    deleteTask,
  }
}
