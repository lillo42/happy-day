<template>
  <div class="product-list-container flex flex-column gap-2">
    <Toast />
    <div class="flex flex-row justify-content-center gap-1">
      <Dropdown v-model="selectedFilter" :options="filterOptions" optionLabel="label" />
      <InputText v-model="productName" placeholder="Search" />
    </div>

    <DataTable :value="products">
      <template #header>
        <div class="flex flex-wrap align-items-center justify-content-between gap-2">
          <span class="text-xl text-900 font-bold">Products</span>
          <Button icon="pi pi-plus" rounded raised @click="addOrEditProduct('new')" />
          <Button icon="pi pi-refresh" rounded raised @click="refreshProducts()" />
        </div>
      </template>
      <Column field="id" header="Id"></Column>
      <Column field="name" header="Name"></Column>
      <Column field="price" header="Price">
        <template #body="slotProps">
          {{ formatCurrency(slotProps.data.price) }}
        </template>
      </Column>
      <Column header="Actions">
        <template #body="slotProps">
          <Button icon="pi pi-pencil" rounded raised @click="addOrEditProduct(slotProps.data.id)" />
          <Button icon="pi pi-trash" severity="danger" rounded raised />
        </template>
      </Column>
    </DataTable>
    <Paginator
      v-model:first="page"
      v-model:rows="pageSize"
      :totalRecords="totatProducts"
      :rowsPerPageOptions="[10, 20, 25, 50]"
    ></Paginator>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'

import type { DropdownOption } from '@/models/common'
import type { ProductService } from '@/services/product.service'

const toast = useToast()
const router = useRouter()
const service = inject<ProductService>('ProductService')
if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const filterOptions = ref<DropdownOption[]>([{ label: 'Name', value: 'name' }])
const selectedFilter = ref<DropdownOption>(filterOptions.value[0])

const productName = ref('')
const page = ref(0)
const pageSize = ref(10)

const data = ref({
  products: [] as ProductView[],
  totatProducts: 0
})

const products = computed(() => data.value.products)
const totatProducts = computed(() => data.value.totatProducts)

watch([productName, page, pageSize], async () => await refreshProducts())

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value)
}

onMounted(async () => await refreshProducts())

const addOrEditProduct = (id: string | null = null) => {
  id = id ?? 'new'
  router.push(`/products/${id}`)
}

const refreshProducts = async () => {
  let tmp = {
    totatProducts: 0,
    products: [] as ProductView[]
  }

  try {
    const response = await service.getAll(productName.value, page.value, pageSize.value)
    if (response.success) {
      const items = response.data!.items ?? []
      tmp = {
        totatProducts: response.data!.totalItems,
        products: items.map((item) => ({
          id: item.id,
          name: item.name,
          price: item.price
        }))
      }
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error to connect API',
        detail: response.error?.message ?? 'unexpected error'
      })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: error.toString() ?? 'unexpected error'
    })
  }

  data.value = tmp
}

interface ProductView {
  id: string
  name: string
  price: number
}
</script>

<style scoped>
.product-list-container {
  justify-content: center;
  padding: 0 10%;
}

.products-filter-container {
  max-width: 20%;
}
</style>
