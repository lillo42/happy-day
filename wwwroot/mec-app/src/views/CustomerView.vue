<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-5 w-3">
      <h1 class="flex flex-row justify-content-center">Customer Details</h1>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="customer.id" id="id" disabled />
        <label for="id">Id</label>
      </span>

      <span class="p-float-label flex flex-column">
        <InputText
          v-model="customer.name"
          id="name"
          rows="5"
          minlength="1"
          maxlength="255"
          :class="{ 'p-invalid': nameIsInvalid }"
        />
        <label for="name">Name</label>
        <small class="invalid" v-if="errors.includes('customer-name-is-empty')">
          Name is empty
        </small>
        <small class="invalid" v-if="errors.includes('customer-name-is-too-large')">
          Name is large than 255
        </small>
      </span>

      <span class="p-float-label flex flex-column">
        <Textarea v-model="customer.comment" id="comment" maxlength="500" />
        <label for="comment">Comment</label>
      </span>

      <div>
        <DataTable
          v-model:value="customer.phones"
          edit-mode="cell"
          @cell-edit-complete="onCellEditComplete"
        >
          <template #header>
            <div class="flex flex-wrap align-items-center justify-content-between gap-2">
              <span class="text-xl text-900 font-bold">Phones</span>
              <Button icon="pi pi-plus" rounded raised @click="addPhone()" />
            </div>
          </template>

          <Column field="number" header="Phone" editor="text">
            <template #editor="slotProps">
              <InputMask
                v-model="slotProps.data.number"
                mask="(99) 9999-9999?9"
                plcaceholder="(99) 9999-9999"
                autofocus
              />
            </template>
          </Column>

          <Column header="actions">
            <template #body="slotProps">
              <div class="flex flex-row gap-2">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  @click="deletePhone(slotProps.index)"
                />
              </div>
            </template>
          </Column>
        </DataTable>

        <small class="invalid" v-if="errors.includes('customer-phone-is-invalid')">
          one or more phone is invalid
        </small>
      </div>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="customer.createdAt" id="createdAt" disabled />
        <label for="createdAt">Created at</label>
      </span>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="customer.updatedAt" id="updatedAt" disabled />
        <label for="updatedAt">Update at</label>
      </span>

      <div class="flex flew-row justify-content-center gap-1">
        <Button label="Save" icon="pi pi-check" @click="save()" :disabled="isInvalid" />
        <Button label="Cancel" severity="secondary" icon="pi pi-times" @click="cancel()" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useToast } from 'primevue/usetoast'

import type { CustomerService } from '@/services/customer.service'
import type { DataTableCellEditCompleteEvent } from 'primevue/datatable'
import type { Customer } from '@/models/customer'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const service = inject<CustomerService>('CustomerService')

if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const isNew = computed(() => route.params.id === 'new')
const customer = reactive({
  id: null as string | null,
  name: '',
  comment: null as string | null,
  phones: [] as { number: string }[],
  createdAt: null as string | null,
  updatedAt: null as string | null
})

const errors = ref<string[]>([])
const isInvalid = computed(() => errors.value.length > 0)
const validate = () => {
  errors.value = []
  if (customer.name.length === 0) {
    errors.value.push('customer-name-is-empty')
  }

  if (customer.name.length > 255) {
    errors.value.push('customer-name-is-too-large')
  }

  for (const phone of customer.phones) {
    if (phone.number.length === 0) {
      errors.value.push('customer-phone-is-invalid')
      break
    }
  }
}

onMounted(async () => await fetchData())
watch(customer, () => validate())

const nameIsInvalid = computed(
  () =>
    errors.value.includes('customer-name-is-empty') ||
    errors.value.includes('customer-name-is-too-large')
)

const addPhone = () => {
  customer.phones.push({ number: '' })
}

const deletePhone = (index: number) => {
  customer.phones.splice(index, 1)
}

const onCellEditComplete = (event: DataTableCellEditCompleteEvent) => {
  let { data, newValue, field } = event
  if (field === 'number') {
    data[field] = newValue
  }
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
  if (isNew.value) {
    return
  }

  const id = route.params.id
  try {
    const response = await service.fetch(id as string)
    if (response.success) {
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

      const formatter = new Intl.DateTimeFormat('en-US', options)
      customer.id = response.data!.id
      customer.name = response.data!.name
      customer.comment = response.data!.comment
      customer.phones = response.data!.phones.map((phone: string) => {
        return { number: format(phone, '(##) ####-#####') }
      })
      customer.createdAt = formatter.format(new Date(response.data!.createAt))
      customer.updatedAt = formatter.format(new Date(response.data!.updateAt))
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error during fetch discount',
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

const cancel = () => {
  router.push({ name: 'customers' })
}

const save = async () => {
  validate()
  if (isInvalid.value) {
    return
  }
  try {
    const data: Customer = {
      id: customer.id!,
      name: customer.name,
      comment: customer.comment,
      phones: customer.phones.map((phone) => {
        let number = phone.number.replace('(', '')
        number = number.replace(')', '')
        number = number.replace('-', '')
        number = number.replace(' ', '')
        return number
      }),
      createAt: new Date(),
      updateAt: new Date()
    }

    const response = isNew.value ? await service.create(data) : await service.update(data)
    if (response.success) {
      router.push({ name: 'customers' })
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error during fetch discount',
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
</script>

<style scoped>
.invalid {
  color: #e24c4c;
}
</style>
