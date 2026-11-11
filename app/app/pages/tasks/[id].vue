
<script setup lang="ts">
import type {
  CreateTaskInput,
  UpdateTaskInput,
} from '~/types/task'

const route = useRoute()
const { getTask, updateTask } = useTasks()

const taskId = Number(route.params.id)

const validId =
  Number.isSafeInteger(taskId) && taskId > 0

const saving = ref(false)
const saveError = ref('')

const {
  data: task,
  status,
  error,
  refresh,
} = await useAsyncData(
  `minitask-task-${route.params.id}`,
  () => getTask(taskId),
  {
    immediate: validId,
  },
)

const loading = computed(() =>
  status.value === 'pending'
)

const notFound = computed(() => {
  if (!validId) return true

  const fetchError = error.value as {
    statusCode?: number
    status?: number
  } | null

  return fetchError?.statusCode === 404
    || fetchError?.status === 404
})

const loadError = computed(() =>
  !loading.value
  && !notFound.value
  && Boolean(error.value)
)

const initialValues = computed<CreateTaskInput | undefined>(() => {
  if (!task.value) return undefined

  return {
    title: task.value.title,
    description: task.value.description,
    status: task.value.status,
    priority: task.value.priority,
  }
})

const formatDate = (value: string) => {
  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return 'Fecha no disponible'
  }

  return date.toLocaleString('es-DO', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

const handleUpdate = async (values: UpdateTaskInput) => {
  if (saving.value || !validId) return

  saving.value = true
  saveError.value = ''

  try {
    await updateTask(taskId, values)

    await refreshNuxtData([
      'minitask-tasks',
      'minitask-summary',
    ])

    await navigateTo('/')
  } catch (err: unknown) {
    console.error('Error al actualizar la tarea:', err)

    const fetchError = err as {
      data?: {
        error?: {
          message?: string
        }
      }
    }

    saveError.value =
      fetchError.data?.error?.message
      ?? 'No se pudo actualizar la tarea.'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-highlighted">
        Detalles de la tarea
      </h2>

      <p class="mt-2 text-muted">
        Consulta y modifica la información de la tarea.
      </p>
    </div>

    <div
      v-if="loading"
      class="py-8 text-center text-muted"
    >
      Cargando tarea...
    </div>

    <UCard v-else-if="notFound">
      <div class="space-y-4 py-6 text-center">
        <h3 class="text-lg font-semibold text-highlighted">
          Tarea no encontrada
        </h3>

        <p class="text-muted">
          La tarea solicitada no existe o el identificador
          no es válido.
        </p>

        <UButton
          to="/"
          label="Volver al listado"
        />
      </div>
    </UCard>

    <UCard v-else-if="loadError">
      <div class="space-y-4 py-6 text-center">
        <p class="text-error">
          No se pudo consultar la tarea.
        </p>

        <UButton
          label="Reintentar"
          @click="refresh"
        />
      </div>
    </UCard>

    <template v-else-if="task && initialValues">
      <UCard>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div>
            <p class="text-sm text-muted">
              ID de la tarea
            </p>

            <p class="mt-1 font-semibold text-highlighted">
              {{ task.id }}
            </p>
          </div>

          <div>
            <p class="text-sm text-muted">
              Fecha de creación
            </p>

            <p class="mt-1 font-semibold text-highlighted">
              {{ formatDate(task.created_at) }}
            </p>
          </div>

          <div>
            <p class="text-sm text-muted">
              Última actualización
            </p>

            <p class="mt-1 font-semibold text-highlighted">
              {{ formatDate(task.updated_at) }}
            </p>
          </div>
        </div>
      </UCard>

      <UAlert
        v-if="saveError"
        color="error"
        variant="soft"
        title="No se pudieron guardar los cambios"
        :description="saveError"
      />

      <TasksTaskForm
        :initial-values="initialValues"
        :saving="saving"
        submit-label="Guardar cambios"
        @submit="handleUpdate"
      />
    </template>
  </section>
</template>
