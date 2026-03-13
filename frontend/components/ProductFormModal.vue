<template>
  <div v-if="open" class="modal-backdrop" @click.self="emit('close')">
    <section class="modal">
      <div class="split-header">
        <div>
          <p class="eyebrow">Products</p>
          <h3>{{ title }}</h3>
        </div>
        <button class="btn-secondary" type="button" @click="emit('close')">Close</button>
      </div>

      <form @submit.prevent="submit">
        <div class="form-grid">
          <label class="full">
            Product name
            <input v-model="form.name" required placeholder="Example: HDMI Cable 2m" />
          </label>

          <label>
            SKU
            <input v-model="form.sku" placeholder="Optional SKU" />
          </label>

          <label>
            Barcode
            <input v-model="form.barcode" placeholder="Optional barcode" />
          </label>

          <label>
            Price
            <input v-model.number="form.price" type="number" min="0" step="0.01" required />
          </label>

          <label>
            Stock
            <input v-model.number="form.stock" type="number" min="0" step="1" required />
          </label>

          <label class="full">
            Active
            <select v-model="form.is_active">
              <option :value="true">Active</option>
              <option :value="false">Inactive</option>
            </select>
          </label>
        </div>

        <div class="toolbar">
          <span class="muted">Changes are saved immediately to PostgreSQL.</span>
          <button class="btn" type="submit">{{ submitLabel }}</button>
        </div>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { Product } from "~/types"

const props = defineProps<{
  open: boolean
  product?: Product | null
}>()

const emit = defineEmits<{
  close: []
  save: [payload: {
    name: string
    sku: string | null
    barcode: string | null
    price: number
    stock: number
    is_active: boolean
  }]
}>()

const form = reactive({
  name: "",
  sku: "",
  barcode: "",
  price: 0,
  stock: 0,
  is_active: true
})

const title = computed(() => (props.product ? "Edit product" : "Add product"))
const submitLabel = computed(() => (props.product ? "Save changes" : "Create product"))

watch(
  () => props.product,
  (product) => {
    form.name = product?.name || ""
    form.sku = product?.sku || ""
    form.barcode = product?.barcode || ""
    form.price = product?.price || 0
    form.stock = product?.stock || 0
    form.is_active = product?.is_active ?? true
  },
  { immediate: true }
)

function submit() {
  emit("save", {
    name: form.name.trim(),
    sku: form.sku.trim() || null,
    barcode: form.barcode.trim() || null,
    price: Number(form.price),
    stock: Number(form.stock),
    is_active: form.is_active
  })
}
</script>
