
<script setup lang="ts">
import type { TaskStatus } from '~/types/task'

type StatusFilter = TaskStatus | 'all'

const props = defineProps<{
  status: StatusFilter
  search: string
}>()

const emit = defineEmits<{
  'update:status': [value: StatusFilter]
  'update:search': [value: string]
}>()

const statusOptions = [
  { label: 'Todos los estados', value: 'all' },
  { label: 'Pendientes', value: 'pending' },
  { label: 'En progreso', value: 'in_progress' },
  { label: 'Completadas', value: 'completed' },
]

const selectedStatus = computed({
  get: () => props.status,
  set: (value: StatusFilter) => {
    emit('update:status', value)
  },
})

const searchText = computed({
  get: () => props.search,
  set: (value: string) => {
    emit('update:search', value)
  },
})
</script>

<template>
  <UCard>
    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
      <UFormField label="Buscar tareas">
        <UInput
          v-model="searchText"
          placeholder="Buscar por título o descripción..."
          icon="i-lucide-search"
          class="w-full"
        />
      </UFormField>

      <UFormField label="Filtrar por estado">
        <USelect
          v-model="selectedStatus"
          :items="statusOptions"
          class="w-full"
        />
      </UFormField>
    </div>
  </UCard>
</template>
