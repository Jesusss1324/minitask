
<script setup lang="ts">
import type { TaskSummary } from '~/types/task'

defineProps<{
  summary: TaskSummary
}>()

const counters = [
  {
    key: 'total',
    label: 'Total de tareas',
    icon: 'i-lucide-list-todo',
  },
  {
    key: 'pending',
    label: 'Pendientes',
    icon: 'i-lucide-clock',
  },
  {
    key: 'in_progress',
    label: 'En progreso',
    icon: 'i-lucide-loader-circle',
  },
  {
    key: 'completed',
    label: 'Completadas',
    icon: 'i-lucide-check-circle',
  },
] as const
</script>

<template>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
    <UCard
      v-for="counter in counters"
      :key="counter.key"
    >
      <div class="flex items-center justify-between">
        <p class="text-sm text-muted">
          {{ counter.label }}
        </p>

        <UIcon
          :name="counter.icon"
          class="size-5 text-muted"
        />
      </div>

      <p class="mt-3 text-3xl font-bold text-highlighted">
        {{ summary[counter.key] }}
      </p>
    </UCard>
  </div>
</template>
