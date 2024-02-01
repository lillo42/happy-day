<template>
  <div class="flex flex-column gap-2 padding-x-10">
    <Toast />
    <ConfirmDialog />

    <DataTable :value="data.items" :loading="loading" :filters="filters" filterDisplay="row">
      <template #header>
        <div class="flex flex-wrap align-items-center justify-content-between gap-2">
          <span class="text-xl text-900 font-bold">Discounts</span>
          <Button icon="pi pi-plus" rounded raised @click="addOrEdit('new')" />
          <Button icon="pi pi-refresh" rounded raised @click="fetchData()" />
        </div>
      </template>

      <template #empty> No discount found. </template>
      <template #loading> Loading discount data. Please wait. </template>

      <Column field="id" header="Id" />
      <Column field="name" header="Name">
        <template #filter>
          <InputText
            v-model="name"
            @input="fetchData()"
            class="p-column-filter"
            placeholder="Search by name"
          />
        </template>
      </Column>
      <Column field="price" header="Price">
        <template #body="slotProps">
          {{ formatCurrency(slotProps.data.price) }}
        </template>
      </Column>

      <Column header="Actions">
        <template #body="slotProps">
          <div class="flex flex-row gap-1">
            <Button icon="pi pi-pencil" rounded raised @click="addOrEdit(slotProps.data.id)" />
            <Button
              icon="pi pi-trash"
              severity="danger"
              rounded
              raised
              @click="deleteDiscount(slotProps.data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>
    <Paginator
      v-model:first="page"
      v-model:rows="pageSize"
      :totalRecords="data.totalItems"
      :rowsPerPageOptions="[10, 20, 25, 50]"
    ></Paginator>
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { FilterMatchMode } from 'primevue/api'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'

import type { DiscountService } from '@/services/discount.service'

const toast = useToast()
const confirm = useConfirm()
const router = useRouter()
const service = inject<DiscountService>('DiscountService')
if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const name = ref('')
const page = ref(0)
const pageSize = ref(10)

const loading = ref(true)
const data = ref({
  items: [] as DiscountView[],
  totalItems: 0
})

onMounted(async () => await fetchData())
watch([name, page, pageSize], async () => await fetchData())

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value)
}

const addOrEdit = (id: string | null = null) => {
  id = id ?? 'new'
  router.push(`/discounts/${id}`)
}

const deleteDiscount = (item: DiscountView) => {
  confirm.require({
    message: `Are you sure you want to delete this discount(name: ${item.name})?`,
    header: 'Product delete confirmation',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    rejectClass: 'p-button-secondary p-button-outlined',
    acceptClass: 'p-button-danger',
    reject: () => {},
    accept: async () => {
      try {
        const response = await service.delete(item.id)
        if (response.success) {
          toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Discount ${item.name} deleted`
          })
          await fetchData()
        } else {
          toast.add({
            severity: 'error',
            summary: 'Error during delete discount',
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
    }
  })
}

const fetchData = async () => {
  loading.value = true
  data.value = {
    totalItems: 0,
    items: []
  }

  try {
    const response = await service.fetchAll(name.value, page.value, pageSize.value)
    if (response.success) {
      const items = response.data!.items ?? []
      data.value = {
        totalItems: response.data!.totalItems,
        items: items.map((item) => ({
          id: item.id,
          name: item.name,
          price: item.price
        }))
      }
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error during fetching data',
        detail: response.error?.message ?? 'unexpected error'
      })
    }
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error'
    })
  }

  loading.value = false
}

interface DiscountView {
  id: string
  name: string
  price: number
}
</script>

<style scoped>
.padding-x-10 {
  padding: 0 10%;
}
</style>
