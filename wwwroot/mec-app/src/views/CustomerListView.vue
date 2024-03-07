<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-2">
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

        <template #empty> No customer found. </template>
        <template #loading> Loading customer data. Please wait. </template>

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

        <Column field="comment" header="Comment">
          <template #filter>
            <InputText
              v-model="comment"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by comment"
            />
          </template>
        </Column>

        <Column field="phone" header="Phone">
          <template #filter>
            <InputText
              v-model="phone"
              @input="fetchData()"
              class="p-column-filter"
              placeholder="Search by name"
            />
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
                @click="deleteCustomer(slotProps.data)"
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
      >
      </Paginator>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { FilterMatchMode } from 'primevue/api'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'

import type { CustomerService } from '@/services/customer.service'

const toast = useToast()
const confirm = useConfirm()
const router = useRouter()
const service = inject<CustomerService>('CustomerService')
if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.CONTAINS },
  comment: { value: null, matchMode: FilterMatchMode.CONTAINS },
  phone: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const name = ref('')
const comment = ref('')
const phone = ref('')
const page = ref(0)
const pageSize = ref(10)

const loading = ref(true)
const data = ref({
  items: [] as CustomerView[],
  totalItems: 0
})

onMounted(async () => await fetchData())

watch([page, pageSize, name, comment, phone], async () => await fetchData())

const addOrEdit = (id: string) => {
  router.push({ name: 'customer', params: { id } })
}

const deleteCustomer = (item: CustomerView) => {
  confirm.require({
    message: `Are you sure you want to delete this customer(name: ${item.name})?`,
    header: 'Customer delete confirmation',
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
            detail: `Customer ${item.name} deleted`
          })
          await fetchData()
        } else {
          toast.add({
            severity: 'error',
            summary: 'Error during delete customer',
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

const format = (value: string, pattern: string): string => {
  let i = 0
  let v = value.toString()
  return pattern.replace(/#/g, () => {
    if (i >= v.length) {
      return ''
    }
    return v[i++] || ''
  })
}

const fetchData = async () => {
  loading.value = true
  data.value = {
    totalItems: 0,
    items: []
  }

  try {
    const response = await service.fetchAll(
      name.value,
      comment.value,
      phone.value,
      page.value,
      pageSize.value
    )
    if (response.success) {
      const items = response.data!.items ?? []
      data.value = {
        totalItems: response.data!.totalItems,
        items: items.map((item) => ({
          id: item.id,
          name: item.name,
          comment: truncate(item.comment),
          phone: truncate(
            item.phones?.map((phone: string) => format(phone, '(##) ####-#####')).join(', ') ?? ''
          )
        }))
      }
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error during fetch customer',
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

  loading.value = false
}

function truncate(message: string | null): string {
  if (message === null) {
    return ''
  }

  if (message.length > 20) {
    return message.substring(0, 20) + '...'
  }

  return message
}

interface CustomerView {
  id: string
  name: string
  comment: string
  phone: string
}
</script>
