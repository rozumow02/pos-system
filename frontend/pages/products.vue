<template>
  <section class="grid">
    <div class="page-header">
      <div>
        <p class="eyebrow">Catalog</p>
        <h2>Products</h2>
        <p>Manage the product list, prices, and current stock levels.</p>
      </div>
      <button class="btn" @click="openCreate">Add product</button>
    </div>

    <div v-if="message" :class="['alert', messageType === 'error' ? 'alert-error' : 'alert-success']">
      {{ message }}
    </div>

    <section class="panel">
      <div class="toolbar">
        <div class="field" style="flex: 1;">
          <input
            v-model="search"
            class="search-input"
            placeholder="Search by name, SKU, or barcode"
            @input="applySearch"
          />
        </div>
        <button class="btn-secondary" @click="reload">Refresh</button>
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>SKU</th>
              <th>Barcode</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in filteredProducts" :key="product.id">
              <td>{{ product.name }}</td>
              <td>{{ product.sku || "-" }}</td>
              <td>{{ product.barcode || "-" }}</td>
              <td>{{ currency(product.price) }}</td>
              <td>{{ product.stock }}</td>
              <td>
                <span :class="['badge', product.is_active ? 'badge-success' : 'badge-warn']">
                  {{ product.is_active ? "Active" : "Inactive" }}
                </span>
              </td>
              <td>
                <button class="btn-secondary" @click="openEdit(product)">Edit</button>
              </td>
            </tr>
            <tr v-if="!filteredProducts.length">
              <td colspan="7" class="muted">No products found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <ProductFormModal
      :open="modalOpen"
      :product="selectedProduct"
      @close="closeModal"
      @save="saveProduct"
    />
  </section>
</template>

<script setup lang="ts">
import type { Product } from "~/types"

const { items, fetchProducts, ensureProducts } = useProducts()
const search = ref("")
const filteredProducts = ref<Product[]>([])
const modalOpen = ref(false)
const selectedProduct = ref<Product | null>(null)
const message = ref("")
const messageType = ref<"success" | "error">("success")

function currency(value: number) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD"
  }).format(value)
}

function applySearch() {
  const normalized = search.value.trim().toLowerCase()
  filteredProducts.value = items.value.filter((product) =>
    product.name.toLowerCase().includes(normalized) ||
    (product.sku || "").toLowerCase().includes(normalized) ||
    (product.barcode || "").toLowerCase().includes(normalized)
  )
}

async function reload() {
  await fetchProducts(search.value)
  applySearch()
}

function openCreate() {
  selectedProduct.value = null
  modalOpen.value = true
}

function openEdit(product: Product) {
  selectedProduct.value = product
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function saveProduct(payload: {
  name: string
  sku: string | null
  barcode: string | null
  price: number
  stock: number
  is_active: boolean
}) {
  message.value = ""

  try {
    if (selectedProduct.value) {
      await useApi(`/products/${selectedProduct.value.id}`, {
        method: "PATCH",
        body: payload
      })
      message.value = "Product updated successfully."
    } else {
      await useApi("/products", {
        method: "POST",
        body: payload
      })
      message.value = "Product created successfully."
    }

    messageType.value = "success"
    closeModal()
    await reload()
  } catch (error: any) {
    message.value = error?.data?.error || error?.message || "Failed to save product"
    messageType.value = "error"
  }
}

await ensureProducts()
applySearch()
</script>
