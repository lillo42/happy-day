<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-2">
      <Toast />
      <ConfrimDialog />

      <DataTable :value="data.items" :loading="loading" :filters="filters" filter-display="row">
        <template #header>
          <div class="flex flex-wrap align-items-center justify-content-between gap-2">
            <span class="text-xl text-900 font-bold">Orders</span>
            <Button icon="pi pi-plus" rounded raised @click="addOrEdit('new')" />
            <Button icon="pi pi-refresh" rounded raised @click="fetchData()" />
          </div>
        </template>

        <template #empty> No orders found. </template>
        <template #loading> Loading orders data. Please wait. </template>

        <Column field="id" header="Id" />

        <Column filed="comment" header="Comment">
          <template #filter>
            <InputText
              v-model="comment"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by comment"
            />
          </template>
        </Column>

        <Column field="address" header="Address">
          <template #filter>
            <InputText
              v-model="address"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by address"
            />
          </template>
        </Column>

        <Column field="customerName" header="Customer Name">
          <template #filter>
            <InputText
              v-model="customerName"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by customer name"
            />
          </template>
        </Column>

        <Column field="customerPhone" header="Customer Phone">
          <template #filter>
            <InputText
              v-model="customerPhone"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by customer phone"
            />
          </template>
        </Column>

        <Column field="finalPrice" header="Final Price">
          <template #body="slotProps">
            {{ formatCurrency(slotProps.data.finalPrice) }}
          </template>
        </Column>

        <Column field="deliveryAt" header="Delivery At">
          <template #body="slotProps">
            {{ formatDate(slotProps.data.deliveryAt) }}
          </template>
        </Column>

        <Column field="pickUpAt" header="Picku At">
          <template #body="slotProps">
            {{ formatDate(slotProps.data.pickUpAt) }}
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
                @click="remove(slotProps.data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>
      <Paginator
        v-model:first="page"
        v-model:rows="pageSize"
        :total-records="data.totalItems"
        :rows-per-page-options="[10, 20, 25, 50]"
      ></Paginator>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { FilterMatchMode } from 'primevue/api'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'

import type { OrderService } from '@/services/order.service'

const router = useRouter()
const toast = useToast()
const confirm = useConfirm()
const service = inject<OrderService>('OrderService')
if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const filters = ref({
  comment: { value: null, matchMode: FilterMatchMode.CONTAINS },
  address: { value: null, matchMode: FilterMatchMode.CONTAINS },
  customerName: { value: null, matchMode: FilterMatchMode.CONTAINS },
  customerPhone: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const comment = ref('')
const address = ref('')
const customerName = ref('')
const customerPhone = ref('')
const page = ref(0)
const pageSize = ref(10)

const loading = ref(true)
const data = ref({
  items: [] as OrderView[],
  totalItems: 0
})

onMounted(async () => await fetchData())
watch(
  [page, pageSize, comment, address, customerName, customerPhone],
  async () => await fetchData()
)

const addOrEdit = (id: string) => {
  router.push({ name: 'order', params: { id } })
}

const remove = (order: OrderView) => {
  confirm.require({
    message: `Are you sure you want to delete order with id ${order.id}?`,
    header: 'Order delete confirmation',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    rejectClass: 'p-button-secondary p-button-outlined',
    acceptClass: 'p-button-danger',
    reject: () => {},
    accept: async () => {
      try {
        const response = await service.delete(order.id)
        if (response.success) {
          toast.add({
            severity: 'success',
            summary: 'Order deleted',
            detail: `Order with id ${order.id} has been deleted`
          })
          await fetchData()
        } else {
          toast.add({
            severity: 'error',
            summary: 'Error during delete order',
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
    }
  })
}

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value)
}

const formatDate = (value: Date) => {
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'numeric',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
    hour12: false,
    timeZone: 'Europe/London'
  }

  return new Intl.DateTimeFormat('en-US', options).format(value)
}

const formatPhone = (value: string, pattern: string): string => {
  let i = 0
  let v = value.toString()
  return pattern.replace(/#/g, () => {
    if (i >= v.length) {
      return ''
    }
    return v[i++] || ''
  })
}

const truncate = (message: string | null): string => {
  if (message === null) {
    return ''
  }

  if (message.length > 20) {
    return message.substring(0, 20) + '...'
  }

  return message
}

const fetchData = async () => {
  loading.value = true
  data.value = { items: [], totalItems: 0 }

  try {
    const response = await service.fetchAll(
      comment.value,
      address.value,
      customerName.value,
      customerPhone.value,
      page.value,
      pageSize.value
    )

    if (response.success) {
      const items = response.data!.items ?? []
      data.value = {
        totalItems: response.data!.totalItems,
        items: items.map((item) => ({
          id: item.id,
          address: truncate(item.address),
          comment: truncate(item.comment),
          customerName: item.customer.name,
          finalPrice: item.finalPrice,
          deliveryAt: item.deliveryAt,
          pickUpAt: item.pickUpAt,
          customerPhone: truncate(
            item.customer.phones.map((phone) => formatPhone(phone, '(##) ####-#####')).join(', ')
          )
        }))
      }
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error during fetch data',
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

interface OrderView {
  id: string
  address: string
  comment: string | null

  customerName: string
  customerPhone: string

  finalPrice: number

  deliveryAt: Date
  pickUpAt: Date
}
</script>
