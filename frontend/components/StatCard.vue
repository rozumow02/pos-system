<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
}>()
</script>

<template>
  <article class="stat-card">
    <p class="stat-label">
      {{ label }}
    </p>
    <p class="stat-value">
      {{ value }}
    </p>
  </article>
</template>
