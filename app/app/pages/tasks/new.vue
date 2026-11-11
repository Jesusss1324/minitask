
<script setup lang="ts">
import type { CreateTaskInput } from '~/types/task'

const { createTask } = useTasks()

const saving = ref(false)
const errorMessage = ref('')

const handleCreate = async (values: CreateTaskInput) => {
  if (saving.value) return

  saving.value = true
  errorMessage.value = ''

  try {
    await createTask(values)

    await navigateTo('/')
  } catch (error: unknown) {
    console.error('Error al crear la tarea:', error)

    const fetchError = error as {
      data?: {
        error?: {
          message?: string
        }
      }
    }

    errorMessage.value =
      fetchError.data?.error?.message
      ?? 'No se pudo crear la tarea. Inténtalo nuevamente.'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-highlighted">
        Crear tarea
      </h2>

      <p class="mt-2 text-muted">
        Completa los campos para registrar una nueva tarea.
      </p>
    </div>

    <UAlert
      v-if="errorMessage"
      color="error"
      variant="soft"
      title="No se pudo guardar"
      :description="errorMessage"
    />

    <TasksTaskForm
      :saving="saving"
      submit-label="Crear tarea"
      @submit="handleCreate"
    />
  </section>
</template>
