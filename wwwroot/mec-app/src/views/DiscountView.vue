<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-5 w-6">
      <h1 class="flex flex-row justify-content-center">Discount Details</h1>
      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="discount.id" id="id" disabled />
        <label for="id">Id</label>
      </span>

      <span class="p-float-label flex flex-column">
        <InputText
          v-model="discount.name"
          id="name"
          minlength="1"
          maxlength="50"
          :class="{ 'p-invalid': nameIsInvalid }"
        />
        <label for="name">Name</label>
        <small class="invalid" v-if="errors.includes('discount-name-is-empty')">
          Name is empty
        </small>
        <small class="invalid" v-if="errors.includes('discount-name-is-too-large')">
          Name is large than 50
        </small>
      </span>

      <span class="p-float-label flex flex-column">
        <InputNumber
          v-model="discount.price"
          id="price"
          :min="0.01"
          :minFractionDigits="2"
          mode="currency"
          currency="USD"
          locale="en-US"
        />
        <label for="price">Price</label>
        <small class="invalid" v-if="errors.includes('discount-price-is-invalid')">
          Price is invalid
        </small>
      </span>

      <div>
        <DataTable
          v-model:value="discount.products"
          editMode="cell"
          @cell-edit-complete="onCellEditComplete"
        >
          <template #header>
            <div class="flex flex-wrap align-items-center justify-content-between gap-2">
              <span class="text-xl text-900 font-bold">Products</span>
              <AutoComplete
                v-model="productFilter"
                dropdown
                optionLabel="name"
                :suggestions="products"
                @complete="fetchProduct"
              />
              <Button icon="pi pi-plus" rounded raised @click="addProduct()" />
            </div>
          </template>
          <Column field="id" header="ID" />
          <Column field="name" header="Name" />
          <Column field="quantity" header="Quantity">
            <template #editor="slotProps">
              <InputNumber v-model="slotProps.data.quantity" autofocus :min="0" />
            </template>
          </Column>
          <Column header="Actions">
            <template #body="slotProps">
              <Button
                icon="pi pi-trash"
                severity="danger"
                rounded
                raised
                @click="deleteProduct(slotProps.data.id)"
              />
            </template>
          </Column>
        </DataTable>

        <small class="invalid" v-if="errors.includes('discount-products-is-missing')">
          Products is missing
        </small>
      </div>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="discount.createdAt" id="createdAt" disabled />
        <label for="createdAt">Created at</label>
      </span>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="discount.updatedAt" id="updatedAt" disabled />
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
import { useRoute, useRouter } from 'vue-router'

import { useToast } from 'primevue/usetoast'
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'

import type { Discount } from '@/models/discount'
import type { DiscountService } from '@/services/discount.service'
import type { ProductService } from '@/services/product.service'
import type { AutoCompleteCompleteEvent } from 'primevue/autocomplete'
import type { DataTableCellEditCompleteEvent } from 'primevue/datatable'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const discountService = inject<DiscountService>('DiscountService')
const productService = inject<ProductService>('ProductService')
if (
  discountService === null ||
  discountService === undefined ||
  productService === null ||
  productService === undefined
) {
  throw new Error('service is null or undefined')
}

const products = ref<DiscountProduct[]>([])
const productFilter = ref<{ id: string; name: string }>()

const isNew = computed(() => route.params.id === 'new')
const discount = reactive({
  id: null as string | null,
  name: '',
  price: 0.01,
  products: [] as DiscountProduct[],
  createdAt: null as string | null,
  updatedAt: null as string | null
})

const errors = ref<string[]>([])
const nameIsInvalid = computed(() => {
  return (
    errors.value.includes('discount-name-is-empty') ||
    errors.value.includes('discount-name-is-too-large')
  )
})

watch(discount, () => validate())

onMounted(async () => {
  await fetchData()
})

const addProduct = () => {
  if (productFilter.value === undefined || productFilter.value === null) {
    return
  }

  const product = products.value.find((product) => product.id === productFilter.value!.id)
  if (product === undefined || product === null) {
    return
  }

  productFilter.value = undefined
  if (discount.products.find((p) => p.id === product.id) !== undefined) {
    return
  }

  discount.products.push({
    id: product.id,
    name: product.name,
    quantity: 0
  })
}

const deleteProduct = (id: string) => {
  discount.products = discount.products.filter((product) => product.id !== id)
}

const onCellEditComplete = (event: DataTableCellEditCompleteEvent) => {
  let { data, newValue, field } = event
  if (field === 'quantity') {
    data[field] = Number.parseInt(newValue)
  }
}

const fetchProduct = async (event: AutoCompleteCompleteEvent | null) => {
  try {
    const response = await productService.getAll(event?.query ?? '', 0, 50)
    if (response.success) {
      products.value = response.data!.items!.map((product) => {
        return {
          id: product.id,
          name: product.name,
          quantity: 0
        }
      })
    } else {
      products.value = []
      toast.add({
        severity: 'error',
        summary: 'Error during fetch products',
        detail: response.error?.message ?? 'unexpected error'
      })
    }
  } catch (err: any) {
    products.value = []
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error'
    })
  }
}

const fetchData = async () => {
  if (isNew.value) {
    return
  }

  const id = route.params.id
  try {
    const response = await discountService.fetch(id as string)
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

      discount.id = response.data!.id
      discount.name = response.data!.name
      discount.price = response.data!.price
      discount.createdAt = formatter.format(new Date(response.data!.createAt))
      discount.updatedAt = formatter.format(new Date(response.data!.updateAt))
      discount.products = response.data!.products.map((product) => {
        return {
          id: product.id,
          name: product.name,
          quantity: product.quantity
        }
      })
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

const isInvalid = computed(() => errors.value.length > 0)
const validate = () => {
  errors.value = []
  if (discount.name.length === 0) {
    errors.value.push('discount-name-is-empty')
  }

  if (discount.name.length > 50) {
    errors.value.push('discount-name-is-too-large')
  }

  if (discount.price < 0) {
    errors.value.push('discount-price-is-invalid')
  }

  if (discount.products.length === 0) {
    errors.value.push('discount-products-is-missing')
  }
}

const save = async () => {
  validate()
  if (isInvalid.value) {
    return
  }

  try {
    const data: Discount = {
      id: discount.id!,
      name: discount.name,
      price: discount.price,
      products: discount.products.map((product) => {
        return {
          id: product.id,
          name: product.name,
          quantity: product.quantity
        }
      }),
      createAt: new Date(),
      updateAt: new Date()
    }

    const response = isNew.value
      ? await discountService.create(data)
      : await discountService.update(data)

    if (response.success) {
      router.push('/discounts')
    } else {
      errors.value.push(response.error!.message)
      toast.add({
        severity: 'error',
        summary: 'Error during saving discount',
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
  router.push('/discounts')
}

interface DiscountProduct {
  id: string
  name: string
  quantity: number
}
</script>

<style scoped>
.invalid {
  color: #e24c4c;
}
</style>
