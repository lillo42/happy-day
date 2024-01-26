<template>
  <div class="flex flex-row justify-content-center">
    <div class="flex flex-column gap-5">
      <h1>Product Details</h1>
      <span class="p-float-label" v-if="!isNew">
        <InputText v-model="product.id" id="id" disabled />
        <label for="id">ID</label>
      </span>

      <span class="p-float-label">
        <InputText
          v-model="product.name"
          id="name"
          minlength="1"
          maxlength="50"
          :class="nameIsInvalid"
        />
        <label for="name">Name</label>
      </span>

      <span class="p-float-label">
        <InputNumber
          v-model="product.price"
          id="price"
          :minFractionDigits="2"
          mode="currency"
          currency="USD"
          locale="en-US"
          :min="0"
        />
        <label for="price">Price</label>
      </span>

      <span class="p-float-label" v-if="!isNew">
        <InputText v-model="product.createdAt" id="createdAt" disabled />
        <label for="createdAt">Created at</label>
      </span>

      <span class="p-float-label" v-if="!isNew">
        <InputText v-model="product.updatedAt" id="updatedAt" disabled />
        <label for="updatedAt">Update at</label>
      </span>

      <div class="flex flew-row gap-1">
        <Button label="Save" icon="pi pi-check" @click="save()" />
        <Button label="Cancel" severity="secondary" icon="pi pi-times" @click="cancel()" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'

import type { Product } from '@/models/product'
import type { ProductService } from '@/services/product.service'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const service = inject<ProductService>('ProductService')
if (service === null || service === undefined) {
  throw new Error('service is null or undefined')
}

const isNew = computed(() => route.params.id === 'new')

const product = reactive({
  id: '',
  name: '',
  price: 0,
  createdAt: '',
  updatedAt: ''
})

const nameIsInvalid = computed<string>(() => {
  if (product.name.length == 0 || product.name.length > 50) {
    return 'p-invalid'
  }

  return ''
})

onMounted(async () => await fetchProduct(route.params.id as string))

const fetchProduct = async (id: string | null) => {
  if (id === null || id === 'new') {
    return
  }

  try {
    const response = await service.get(id)
    if (response.success) {
      // sometimes even the US needs 24-hour time
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

      product.id = response.data!.id
      product.name = response.data!.name
      product.price = response.data!.price
      product.createdAt = formatter.format(response.data!.createdAt)
      product.updatedAt = formatter.format(response.data!.updatedAt)
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error',
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

const save = async () => {
  try {
    const data: Product = {
      id: product.id,
      name: product.name,
      price: product.price,
      createdAt: new Date(),
      updatedAt: new Date()
    }

    const response = isNew.value ? await service.create(data) : await service.update(data)

    if (response.success) {
      router.push('/products')
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error',
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
  router.push('/products')
}
</script>
