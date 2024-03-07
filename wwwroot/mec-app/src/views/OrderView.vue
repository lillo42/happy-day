<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-5 w-3">
      <h1 class="flex flex-row justify-content-center">Order Details</h1>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="order.id" id="id" disabled />
        <label for="id">Id</label>
      </span>

      <span class="p-float-label">
        <Textarea
          v-model="order.address"
          row="5"
          id="address"
          maxlength="1000"
          :class="{ 'p-invalid': isAddressInvalid }"
        />
        <label for="address">Address</label>
        <small class="invalid" v-if="errors.includes('order-address-is-empty')">
          Address is empty
        </small>
        <small class="invalid" v-if="errors.includes('order-address-is-too-large')">
          Address is too large
        </small>
      </span>

      <span class="p-float-label">
        <Textarea v-model="order.comment" row="5" id="comment" />
        <label for="comment">Comment</label>
      </span>

      <span class="p-float-label">
        <Calendar
          v-model="order.delivery"
          selection-mode="range"
          show-time
          hour-format="24"
          id="delivery"
          :class="{ 'p-invalid': isDeliveryInvalid }"
        />
        <label for="delivery">Delivery/Pickup</label>
        <small class="invalid" v-if="errors.includes('order-delivery-at-is-invalid')">
          Delivery/Pickup is invalid
        </small>
      </span>

      <span class="p-float-label">
        <InputNumber
          v-model="order.totalPrice"
          id="totalPrice"
          :min-fraction-digits="2"
          mode="currency"
          currency="USD"
          locale="en-US"
          :class="{ 'p-invalid': isTotalPriceInvalid }"
          disabled
        />
        <label for="totalPrice">Total Price</label>
        <span class="invalid" v-if="errors.includes('order-total-price-is-invalid')">
          Total Price is invalid
        </span>
      </span>

      <span class="p-float-label">
        <InputNumber
          v-model="order.discount"
          id="discount"
          :min-fraction-digits="2"
          :min="0"
          mode="currency"
          currency="USD"
          locale="en-US"
          :class="{ 'p-invalid': isDiscountInvalid }"
        />
        <label for="discount">Discount</label>
        <span class="invalid" v-if="errors.includes('order-discount-is-invalid')">
          Discount is invalid
        </span>
      </span>

      <span class="p-float-label">
        <InputNumber
          v-model="order.finalPrice"
          id="finalPrice"
          :min-fraction-digits="2"
          :min="0.01"
          mode="currency"
          currency="USD"
          locale="en-US"
          :class="{ 'p-invalid': isFinalPriceInvalid }"
        />
        <label for="finalPrice">Final price</label>
        <span class="invalid" v-if="errors.includes('order-final-price-is-invalid')">
          Final Price is invalid
        </span>
      </span>

      <div>
        <DataTable
          v-model:value="order.products"
          edit-mode="cell"
          @cell-edit-complete="onCellEditComplete"
        >
          <template #header>
            <div class="flex flex-wrap align-items-center justify-content-between gap-2">
              <span class="text-xl text-900 font-bold">Products</span>
              <AutoComplete
                v-model="productFilter"
                dropdown
                optionLabel="name"
                :suggestions="apiProducts"
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
                @click="removeProduct(slotProps.data.id)"
              />
            </template>
          </Column>
        </DataTable>
      </div>

      <div>
        <DataTable
          v-model:value="order.payments"
          edit-mode="cell"
          @cell-edit-complete="onCellEditComplete"
        >
          <template #header>
            <div class="flex flex-wrap align-items-center justify-content-between gap-2">
              <span class="text-xl text-900 font-bold">Payments</span>
              <Button icon="pi pi-plus" rounded raised @click="addPayment()" />
            </div>
          </template>

          <Column field="at" header="At">
            <template #editor="slotProps">
              <Calendar v-model="slotProps.data.at" show-time hour-format="24" id="at" />
            </template>
          </Column>
          <Column field="method" header="Method">
            <template #editor="slotProps">
              <Dropdown
                v-model="slotProps.data.method"
                :options="paymentMethods"
                optionLabel="label"
              />
            </template>
          </Column>
          <Column field="amount" header="Amount">
            <template #editor="slotProps">
              <InputNumber
                v-model="slotProps.data.amount"
                autofocus
                :min-fraction-digits="2"
                :min="0.01"
                mode="currency"
                currency="USD"
                locale="en-US"
              />
            </template>
          </Column>

          <Column field="info" header="Info">
            <template #editor="slotProps">
              <InputText v-model="slotProps.data.info" autofocus />
            </template>
          </Column>

          <Column header="Actions">
            <template #body="slotProps">
              <Button
                icon="pi pi-trash"
                severity="danger"
                rounded
                raised
                @click="removePayment(slotProps.index)"
              />
            </template>
          </Column>
        </DataTable>
      </div>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="order.createdAt" id="createdAt" disabled />
        <label for="createdAt">Created at</label>
      </span>

      <span class="p-float-label flex flex-column" v-if="!isNew">
        <InputText v-model="order.updatedAt" id="updatedAt" disabled />
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
import type { AutoCompleteCompleteEvent } from 'primevue/autocomplete'
import type { DataTableCellEditCompleteEvent } from 'primevue/datatable'

import type { OrderService } from '@/services/order.service'
import type { OrderCustomer, OrderPayment, OrderProduct } from '@/models/order'
import type { ProductService } from '@/services/product.service'
import type { CustomerService } from '@/services/customer.service'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const orderService = inject<OrderService>('OrderService')
const productService = inject<ProductService>('ProductService')
const customerService = inject<CustomerService>('CustomerService')
if (orderService === undefined || productService === undefined || customerService === undefined) {
  throw new Error('Service is not provided')
}

const isNew = computed(() => route.params.id === 'new')
const order = reactive({
  id: null as string | null,
  address: '',
  comment: null as string | null,

  delivery: [] as Date[],

  totalPrice: 0,
  discount: 0,
  finalPrice: 0,

  products: [] as OrderProduct[],
  payments: [] as OrderPayment[],
  customer: null as OrderCustomer | null,

  createdAt: null as string | null,
  updatedAt: null as string | null
})

const errors = ref<string[]>([])
const isInvalid = computed(() => errors.value.length > 0)
const isAddressInvalid = computed(() => errors.value.findIndex((x) => x.indexOf('address') > -1))
const isDeliveryInvalid = computed(() => errors.value.findIndex((x) => x.indexOf('at') > -1))
const isTotalPriceInvalid = computed(() =>
  errors.value.findIndex((x) => x.indexOf('total-price') > -1)
)
const isDiscountInvalid = computed(() => errors.value.findIndex((x) => x.indexOf('discount') > -1))
const isFinalPriceInvalid = computed(() =>
  errors.value.findIndex((x) => x.indexOf('final-price') > -1)
)

const onCellEditComplete = (event: DataTableCellEditCompleteEvent) => {
  let { data, newValue, field } = event
  if (field === 'quantity' || field === 'amount') {
    data[field] = Number.parseInt(newValue)
  } else {
    data[field] = newValue
  }
}

onMounted(async () => await fetch())
watch([order.products], async () => await quote())
watch([order.totalPrice, order.discount], () => {
  order.finalPrice = order.totalPrice - order.discount
})

// begin Payment
const paymentMethods = ref([
  { label: 'Pix', value: 'pix' },
  { label: 'Bank Transfer', value: 'bank-transfer' },
  { label: 'Cash', value: 'cash' }
])

const addPayment = () => {
  order.payments.push({
    amount: 0,
    at: new Date(),
    info: null,
    method: 'pix'
  })
}

const removePayment = (index: number) => {
  order.payments.splice(index, 1)
}
// end Payment

// beging Products
const apiProducts = ref<OrderProduct[]>([])
const productFilter = ref<{ id: string; name: string }>()

const addProduct = () => {
  if (productFilter.value === undefined || productFilter.value === null) {
    return
  }

  const product = apiProducts.value.find((product) => product.id === productFilter.value!.id)
  if (product === undefined || product === null) {
    return
  }

  productFilter.value = undefined
  if (order.products.find((p) => p.id === product.id) !== undefined) {
    return
  }

  order.products.push({
    id: product.id,
    name: product.name,
    quantity: 1,
    price: product.price
  })
}

const removeProduct = (id: string) => {
  order.products = order.products.filter((product) => product.id !== id)
}

const fetchProduct = async (event: AutoCompleteCompleteEvent | null) => {
  try {
    const response = await productService.getAll(event?.query || '', 0, 50)
    if (response.success) {
      apiProducts.value =
        response.data!.items?.map((product) => ({
          id: product.id,
          name: product.name,
          quantity: 0,
          price: product.price
        })) ?? []
    } else {
      apiProducts.value = []
      toast.add({
        severity: 'error',
        summary: 'Error to connect API',
        detail: response.error?.message ?? 'unexpected error',
        life: 3000
      })
    }
  } catch (err: any) {
    apiProducts.value = []

    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error',
      life: 3000
    })
  }
}

// end Products

// beging Customer
const apiCustomers = ref<OrderCustomer[]>([])

const fetchCustomer = async (event: AutoCompleteCompleteEvent | null) => {
  try {
    const response = await customerService.fetchAll(event?.query || '', null, null, 0, 50)
    if (response.success) {
      const items = response.data!.items ?? []
      apiCustomers.value = items.map((customer) => ({
        id: customer.id,
        name: customer.name,
        comment: customer.comment,
        phones: customer.phones
      }))
    } else {
      apiCustomers.value = []
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: response.error?.message ?? 'unexpected error',
        life: 3000
      })
    }
  } catch (err: any) {
    apiCustomers.value = []
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error',
      life: 3000
    })
  }
}

// end Customer

const cancel = () => {
  router.push({ name: 'orders' })
}

const save = async () => {}

const fetch = async () => {
  if (isNew.value) {
    return
  }

  try {
    const response = await orderService.fetch(route.params.id as string)
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

      const data = response.data!
      order.id = data.id
      order.address = data.address
      order.comment = data.comment
      order.delivery = [data.deliveryAt, data.pickUpAt]
      order.totalPrice = data.totalPrice
      order.discount = data.discount
      order.finalPrice = data.finalPrice
      order.products = data.products
      order.payments = data.payments
      order.customer = data.customer
      order.createdAt = formatter.format(data.createAt)
      order.updatedAt = formatter.format(data.updateAt)
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error to connect API',
        detail: response.error?.message ?? 'unexpected error',
        life: 3000
      })
    }
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error',
      life: 3000
    })
  }
}

const quote = async () => {
  try {
    const response = await orderService.quote(order.products)
    if (response.success) {
      order.totalPrice = response.data!
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: response.error?.message ?? 'unexpected error',
        life: 3000
      })
    }
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error to connect API',
      detail: err.toString() ?? 'unexpected error',
      life: 3000
    })
  }
}
</script>
