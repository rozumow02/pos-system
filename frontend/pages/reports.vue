<script setup lang="ts">
const { items: products, ensureProducts, fetchProducts } = useProducts()
const { data, refresh } = useDashboard()
const { topProducts, lowStock, refreshTopProducts, refreshLowStock } = useReports()

const dateFrom = ref('')
const dateTo = ref('')
const errorMessage = ref('')

const dashboard = computed(() => data.value)

function currency(value: number) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(value)
}

async function load() {
  errorMessage.value = ''
  try {
    await Promise.all([
      refresh(),
      refreshTopProducts(dateFrom.value, dateTo.value, 10),
      refreshLowStock(),
      fetchProducts(),
    ])
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Failed to load reports'
  }
}

await ensureProducts()
await load()
</script>

<template>
  <section class="grid">
    <div class="page-header">
      <div>
        <p class="eyebrow">
          Analytics
        </p>
        <h2>Reports</h2>
        <p>Review daily revenue, best sellers, and full inventory status.</p>
      </div>
      <button class="btn-secondary" @click="load">
        Refresh
      </button>
    </div>

    <div v-if="errorMessage" class="alert alert-error">
      {{ errorMessage }}
    </div>

    <section class="panel">
      <div class="filters">
        <label>
          Date from
          <input v-model="dateFrom" type="date">
        </label>
        <label>
          Date to
          <input v-model="dateTo" type="date">
        </label>
        <button class="btn" @click="load">
          Apply filters
        </button>
      </div>
    </section>

    <div class="grid stats">
      <StatCard label="Revenue today" :value="currency(dashboard?.revenue_today || 0)" />
      <StatCard label="Orders today" :value="dashboard?.orders_today || 0" />
      <StatCard label="Items sold today" :value="dashboard?.items_sold_today || 0" />
    </div>

    <div class="grid" style="grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));">
      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Performance
            </p>
            <h2>Top products</h2>
          </div>
        </div>
        <div class="list-plain">
          <div v-for="item in topProducts" :key="item.product_id" class="list-row">
            <strong>{{ item.name }}</strong>
            <p>{{ item.quantity_sold }} units sold</p>
            <p class="muted">
              {{ currency(item.revenue) }}
            </p>
          </div>
          <p v-if="!topProducts.length" class="muted">
            No product sales found for this range.
          </p>
        </div>
      </section>

      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Inventory risk
            </p>
            <h2>Low stock</h2>
          </div>
        </div>
        <div class="list-plain">
          <div v-for="item in lowStock" :key="item.id" class="list-row">
            <strong>{{ item.name }}</strong>
            <p>{{ item.stock }} units remaining</p>
            <p class="muted">
              {{ currency(item.price) }}
            </p>
          </div>
          <p v-if="!lowStock.length" class="muted">
            No low stock alerts right now.
          </p>
        </div>
      </section>
    </div>

    <section class="panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">
            Inventory
          </p>
          <h2>Inventory status</h2>
        </div>
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>SKU</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in products" :key="product.id">
              <td>{{ product.name }}</td>
              <td>{{ product.sku || "-" }}</td>
              <td>{{ currency(product.price) }}</td>
              <td>{{ product.stock }}</td>
              <td>{{ product.is_active ? "Active" : "Inactive" }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </section>
</template>
