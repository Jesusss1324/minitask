
<script setup lang="ts">
import type { Task, TaskStatus } from '~/types/task'

type StatusFilter = TaskStatus | 'all'

const { getTasks, getSummary } = useTasks()

const selectedStatus = ref<StatusFilter>('all')
const searchText = ref('')
const debouncedSearch = ref('')

let searchTimeout: ReturnType<typeof setTimeout> | undefined

watch(searchText, (value) => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }

  searchTimeout = setTimeout(() => {
    debouncedSearch.value = value.trim()
  }, 350)
})

onBeforeUnmount(() => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
})

const filters = computed(() => ({
  status: selectedStatus.value === 'all'
    ? undefined
    : selectedStatus.value,
  search: debouncedSearch.value || undefined,
}))

const {
  data: tasks,
  status: tasksStatus,
  error: tasksError,
  refresh: refreshTasks,
} = await useAsyncData(
  'minitask-tasks',
  () => getTasks(filters.value),
  {
    default: () => [] as Task[],
    watch: [filters],
  },
)

const {
  data: summary,
  status: summaryStatus,
  error: summaryError,
  refresh: refreshSummary,
} = await useAsyncData(
  'minitask-summary',
  () => getSummary(),
)

const loading = computed(() =>
  tasksStatus.value === 'pending'
  || summaryStatus.value === 'pending'
)

const hasError = computed(() =>
  Boolean(tasksError.value || summaryError.value)
)

const retry = async () => {
  await Promise.all([
    refreshTasks(),
    refreshSummary(),
  ])
}

const hasActiveFilters = computed(() =>
  selectedStatus.value !== 'all'
  || searchText.value.trim() !== ''
)

const statusLabels: Record<TaskStatus, string> = {
  pending: 'Pendiente',
  in_progress: 'En progreso',
  completed: 'Completada',
}

const priorityLabels: Record<Task['priority'], string> = {
  low: 'Baja',
  medium: 'Media',
  high: 'Alta',
}
</script>

<template>
  <section>
    <div class="mb-6 flex items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-highlighted">
          Mis tareas
        </h2>

        <p class="mt-2 text-muted">
          Administra tus tareas desde un solo lugar.
        </p>
      </div>

      <UButton
        to="/tasks/new"
        label="Crear tarea"
        icon="i-lucide-plus"
      />
    </div>

    <div
      v-if="hasError"
      class="py-8 text-center"
    >
      <p class="mb-4 text-error">
        No se pudieron cargar los datos.
      </p>

      <UButton
        label="Reintentar"
        @click="retry"
      />
    </div>

    <template v-else>
      <TasksTaskCounters
        v-if="summary"
        :summary="summary"
        class="mb-6"
      />

      <TasksTaskFilters
        v-model:status="selectedStatus"
        v-model:search="searchText"
        class="mb-6"
      />

      <UCard>
        <div
          v-if="loading"
          class="py-8 text-center text-muted"
        >
          Cargando tareas...
        </div>

        <div
          v-else-if="tasks?.length === 0"
          class="py-8 text-center"
        >
          <template v-if="hasActiveFilters">
            <h3 class="text-lg font-semibold text-highlighted">
              No se encontraron resultados
            </h3>

            <p class="mt-2 text-muted">
              Prueba con otros criterios de búsqueda.
            </p>
          </template>

          <template v-else>
            <h3 class="text-lg font-semibold text-highlighted">
              No tienes tareas registradas
            </h3>

            <p class="mt-2 text-muted">
              Crea tu primera tarea para comenzar.
            </p>

            <UButton
              to="/tasks/new"
              class="mt-4"
              label="Crear tarea"
            />
          </template>
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="task in tasks"
            :key="task.id"
            class="rounded-lg border border-default p-4"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h3 class="text-lg font-semibold text-highlighted">
                  {{ task.title }}
                </h3>

                <p class="mt-2 text-muted">
                  {{ task.description || 'Sin descripción' }}
                </p>
              </div>

              <UBadge color="neutral" variant="soft">
                {{ statusLabels[task.status] }}
              </UBadge>
            </div>

            <div class="mt-4 flex items-center justify-between">
              <span class="text-sm text-muted">
                Prioridad: {{ priorityLabels[task.priority] }}
              </span>

              <UButton
                :to="`/tasks/${task.id}`"
                label="Ver detalles"
                variant="outline"
                size="sm"
              />
            </div>
          </div>
        </div>
      </UCard>
    </template>
  </section>
</template>
