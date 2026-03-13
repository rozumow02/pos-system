<script setup lang="ts">
const { data, refresh } = useDashboard()
const metrics = computed(() => data.value)
const errorMessage = ref('')

function currency(value: number) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(value)
}

async function load() {
  errorMessage.value = ''
  try {
    await refresh()
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Failed to load dashboard'
  }
}

await load()
</script>

<template>
  <section class="grid">
    <div class="page-header">
      <div>
        <p class="eyebrow">
          Overview
        </p>
        <h2>Dashboard</h2>
        <p>Quick view of today's business activity and inventory risk.</p>
      </div>
      <button class="btn-secondary" @click="load">
        Refresh
      </button>
    </div>

    <div v-if="errorMessage" class="alert alert-error">
      {{ errorMessage }}
    </div>

    <div class="grid stats">
      <StatCard label="Revenue today" :value="currency(metrics?.revenue_today || 0)" />
      <StatCard label="Orders today" :value="metrics?.orders_today || 0" />
      <StatCard label="Items sold" :value="metrics?.items_sold_today || 0" />
      <StatCard label="Low stock alerts" :value="metrics?.low_stock_count || 0" />
    </div>

    <div class="grid" style="grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));">
      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Best sellers
            </p>
            <h2>Top products</h2>
          </div>
        </div>

        <div class="list-plain">
          <div v-for="item in metrics?.top_products || []" :key="item.product_id" class="list-row">
            <strong>{{ item.name }}</strong>
            <p>{{ item.quantity_sold }} units sold</p>
            <p class="muted">
              {{ currency(item.revenue) }}
            </p>
          </div>
          <p v-if="!(metrics?.top_products || []).length" class="muted">
            No sales yet today.
          </p>
        </div>
      </section>

      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Inventory
            </p>
            <h2>Low stock</h2>
          </div>
        </div>

        <div class="list-plain">
          <div v-for="product in metrics?.low_stock || []" :key="product.id" class="list-row">
            <strong>{{ product.name }}</strong>
            <p>{{ product.stock }} in stock</p>
            <p class="muted">
              {{ currency(product.price) }}
            </p>
          </div>
          <p v-if="!(metrics?.low_stock || []).length" class="muted">
            Stock levels are healthy.
          </p>
        </div>
      </section>
    </div>
  </section>
</template>
