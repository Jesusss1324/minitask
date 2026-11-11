
<script setup lang="ts">
import type {
  CreateTaskInput,
  TaskStatus,
  TaskPriority,
} from '~/types/task'

const props = withDefaults(defineProps<{
  initialValues?: CreateTaskInput
  saving?: boolean
  submitLabel?: string
}>(), {
  saving: false,
  submitLabel: 'Guardar tarea',
})

const emit = defineEmits<{
  submit: [values: CreateTaskInput]
}>()

const form = reactive<CreateTaskInput>({
  title: props.initialValues?.title ?? '',
  description: props.initialValues?.description ?? '',
  status: props.initialValues?.status ?? 'pending',
  priority: props.initialValues?.priority ?? 'medium',
})

const statusOptions = [
  { label: 'Pendiente', value: 'pending' },
  { label: 'En progreso', value: 'in_progress' },
  { label: 'Completada', value: 'completed' },
]

const priorityOptions = [
  { label: 'Baja', value: 'low' },
  { label: 'Media', value: 'medium' },
  { label: 'Alta', value: 'high' },
]

const validStatuses: TaskStatus[] = [
  'pending',
  'in_progress',
  'completed',
]

const validPriorities: TaskPriority[] = [
  'low',
  'medium',
  'high',
]

const status = computed({
  get: () => form.status,
  set: (value: TaskStatus) => {
    form.status = value
  },
})

const priority = computed({
  get: () => form.priority,
  set: (value: TaskPriority) => {
    form.priority = value
  },
})

const validationError = ref('')
const isDirty = ref(false)

watch(
  form,
  () => {
    isDirty.value = true
  },
  { deep: true },
)

watch(
  () => props.initialValues,
  (values) => {
    if (!values || isDirty.value) return

    Object.assign(form, values)
    isDirty.value = false
  },
)

const handleSubmit = () => {
  if (props.saving) return

  validationError.value = ''

  const title = form.title.trim()

  if ([...title].length < 3 || [...title].length > 100) {
    validationError.value =
      'El título debe contener entre 3 y 100 caracteres.'
    return
  }

  if ([...form.description].length > 1000) {
    validationError.value =
      'La descripción no puede superar los 1000 caracteres.'
    return
  }

  if (!validStatuses.includes(form.status)) {
    validationError.value = 'Selecciona un estado válido.'
    return
  }

  if (!validPriorities.includes(form.priority)) {
    validationError.value = 'Selecciona una prioridad válida.'
    return
  }

  emit('submit', {
    title,
    description: form.description,
    status: form.status,
    priority: form.priority,
  })
}
</script>

<template>
  <UCard>
    <form class="space-y-6" @submit.prevent="handleSubmit">
      <UAlert
        v-if="validationError"
        color="error"
        variant="soft"
        :description="validationError"
      />

      <UFormField label="Título" required>
        <UInput
          v-model="form.title"
          placeholder="Escribe el título de la tarea"
          :disabled="saving"
          class="w-full"
        />
      </UFormField>

      <UFormField label="Descripción">
        <UTextarea
          v-model="form.description"
          placeholder="Describe los detalles de la tarea"
          :rows="4"
          :disabled="saving"
          class="w-full"
        />
      </UFormField>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <UFormField label="Estado" required>
          <USelect
            v-model="status"
            :items="statusOptions"
            :disabled="saving"
            class="w-full"
          />
        </UFormField>

        <UFormField label="Prioridad" required>
          <USelect
            v-model="priority"
            :items="priorityOptions"
            :disabled="saving"
            class="w-full"
          />
        </UFormField>
      </div>

      <div class="flex justify-end gap-3">
        <UButton
          to="/"
          label="Cancelar"
          color="neutral"
          variant="outline"
          :disabled="saving"
        />

        <UButton
          type="submit"
          :label="submitLabel"
          icon="i-lucide-save"
          :loading="saving"
          :disabled="saving"
        />
      </div>
    </form>
  </UCard>
</template>
