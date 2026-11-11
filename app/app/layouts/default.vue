
<template>
  <div class="min-h-screen flex flex-col">
    <!-- Header -->
    <header class="border-b border-gray-200 bg-white">
      <div class="mx-auto flex h-16 max-w-7xl items-center px-6">
        <h1 class="text-xl font-bold text-gray-900">
          MiniTask
        </h1>
      </div>
    </header>

    <!-- Main content -->
    <div class="flex flex-1">
      <!-- Sidebar -->
      <aside class="hidden w-64 border-r border-gray-200 bg-gray-50 p-6 md:block">
        <nav class="space-y-2">
          <NuxtLink
            to="/"
            class="block rounded-lg px-4 py-2 text-gray-700 hover:bg-gray-200"
          >
            Tareas
          </NuxtLink>

          <NuxtLink
            to="/tasks/new"
            class="block rounded-lg px-4 py-2 text-gray-700 hover:bg-gray-200"
          >
            Crear tarea
          </NuxtLink>
        </nav>
      </aside>

      <!-- Body -->
      <main class="flex-1 p-6">
        <slot />
      </main>
    </div>

    <!-- Footer -->
    <footer class="border-t border-gray-200 bg-white px-6 py-4">
      <p class="text-center text-sm text-gray-500">
        MiniTask - Gestor de tareas
      </p>
    </footer>
  </div>
</template>
