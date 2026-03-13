<script setup lang="ts">
const { items, ensureProducts, fetchProducts } = useProducts()
const { cart, total, totalItems, addProduct, updateQuantity, removeProduct, clearCart } = usePosCart()

const query = ref('')
const message = ref('')
const messageType = ref<'success' | 'error'>('success')
const selling = ref(false)

const visibleProducts = computed(() => {
  const normalized = query.value.trim().toLowerCase()
  const source = items.value.filter(product => product.is_active)
  if (!normalized) {
    return source.slice(0, 30)
  }

  return source
    .filter(product =>
      product.name.toLowerCase().includes(normalized)
      || (product.sku || '').toLowerCase().includes(normalized)
      || (product.barcode || '').toLowerCase().includes(normalized),
    )
    .slice(0, 30)
})

function currency(value: number) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(value)
}

function onQuantityInput(productId: number, event: Event) {
  const target = event.target as HTMLInputElement
  updateQuantity(productId, Number(target.value))
}

function getErrorMessage(error: unknown) {
  if (error && typeof error === 'object') {
    const maybeError = error as {
      data?: { error?: string }
      message?: string
    }

    return maybeError.data?.error || maybeError.message || 'Failed to complete sale'
  }

  return 'Failed to complete sale'
}

async function refreshProducts() {
  await fetchProducts()
}

async function sell() {
  if (!cart.value.length)
    return

  selling.value = true
  message.value = ''

  try {
    await useApi('/orders', {
      method: 'POST',
      body: {
        items: cart.value.map(item => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      },
    })

    clearCart()
    await refreshProducts()
    message.value = 'Sale completed successfully.'
    messageType.value = 'success'
  }
  catch (error: unknown) {
    message.value = getErrorMessage(error)
    messageType.value = 'error'
    await refreshProducts()
  }
  finally {
    selling.value = false
  }
}

await ensureProducts()
</script>

<template>
  <section class="grid">
    <div class="page-header">
      <div>
        <p class="eyebrow">
          Checkout
        </p>
        <h2>POS</h2>
        <p>Search products on the left, review cart on the right, then complete the sale.</p>
      </div>
      <button class="btn-secondary" @click="refreshProducts">
        Refresh stock
      </button>
    </div>

    <div v-if="message" class="alert" :class="[messageType === 'error' ? 'alert-error' : 'alert-success']">
      {{ message }}
    </div>

    <div class="pos-grid">
      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Find items
            </p>
            <h2>Product search</h2>
          </div>
        </div>

        <input
          v-model="query"
          class="search-input"
          placeholder="Search by name, SKU, or barcode"
        >

        <div class="product-list">
          <article v-for="product in visibleProducts" :key="product.id" class="product-item">
            <div>
              <h4>{{ product.name }}</h4>
              <p>{{ currency(product.price) }}</p>
              <p class="muted">
                Stock: {{ product.stock }} | SKU: {{ product.sku || "-" }}
              </p>
            </div>
            <button class="btn" :disabled="product.stock <= 0" @click="addProduct(product)">
              {{ product.stock > 0 ? "Add" : "Out of stock" }}
            </button>
          </article>
          <p v-if="!visibleProducts.length" class="muted">
            No products match your search.
          </p>
        </div>
      </section>

      <section class="panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">
              Current sale
            </p>
            <h2>Cart</h2>
          </div>
          <span class="badge badge-success">{{ totalItems }} items</span>
        </div>

        <div class="cart-list">
          <article v-for="item in cart" :key="item.product_id" class="cart-item">
            <div class="cart-row">
              <div>
                <h4>{{ item.name }}</h4>
                <p class="muted">
                  {{ currency(item.price) }} each
                </p>
              </div>
              <button class="btn-danger" @click="removeProduct(item.product_id)">
                Remove
              </button>
            </div>

            <div class="cart-actions" style="margin-top: 0.8rem;">
              <label style="flex: 1;">
                Quantity
                <input
                  :value="item.quantity"
                  class="qty-input"
                  type="number"
                  min="1"
                  :max="item.stock"
                  @input="onQuantityInput(item.product_id, $event)"
                >
              </label>
              <strong>{{ currency(item.quantity * item.price) }}</strong>
            </div>
          </article>
          <p v-if="!cart.length" class="muted">
            Cart is empty.
          </p>
        </div>

        <div class="cart-footer">
          <div class="summary-row">
            <span>Total</span>
            <span class="total-value">{{ currency(total) }}</span>
          </div>
          <button class="btn" style="width: 100%; margin-top: 1rem;" :disabled="selling || !cart.length" @click="sell">
            {{ selling ? "Processing..." : "SELL" }}
          </button>
        </div>
      </section>
    </div>
  </section>
</template>
