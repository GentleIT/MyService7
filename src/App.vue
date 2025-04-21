<template>
  <div class="container">
    <h1>🛒 Каталог товаров</h1>

    <form class="form" @submit.prevent="addProduct">
      <input v-model="product.name" placeholder="Название товара" required />
      <input v-model="product.description" placeholder="Описание" />
      <input v-model.number="product.price" type="number" placeholder="Цена (₽)" required min="0" step="0.01" />
      <button :disabled="loading">Добавить товар</button>
    </form>

    <p v-if="message" class="message">{{ message }}</p>

    <div class="list">
      <h2>📦 Список товаров</h2>
      <div v-if="products.length === 0">Нет добавленных товаров.</div>
      <ul>
        <li v-for="(item, index) in products" :key="index" class="product-card">
          <strong>{{ item.name }}</strong>
          <p>{{ item.description }}</p>
          <span class="price">{{ item.price }} ₽</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const product = ref({
  name: '',
  description: '',
  price: 0
})

const products = ref([])
const message = ref('')
const loading = ref(false)

const getProducts = async () => {
  try {
    const response = await axios.get('http://localhost:8080/products')
    products.value = response.data
  } catch (err) {
    console.error("Ошибка при получении товаров:", err)
  }
}

const addProduct = async () => {
  if (!product.value.name || product.value.price <= 0) {
    message.value = 'Пожалуйста, заполните все поля корректно.'
    return
  }

  loading.value = true
  try {
    await axios.post('http://localhost:8080/products', product.value)
    await getProducts()
    message.value = '✅ Товар успешно добавлен!'
    product.value = { name: '', description: '', price: 0 }
  } catch (err) {
    console.error("Ошибка при добавлении товара:", err)
    message.value = '❌ Ошибка при добавлении товара.'
  } finally {
    loading.value = false
    setTimeout(() => (message.value = ''), 3000)
  }
}

onMounted(getProducts)
</script>

<style scoped>
body {
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  background: #f3f4f6;
  margin: 0;
  padding: 0;
}

.container {
  max-width: 700px;
  margin: 40px auto;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
}

h1, h2 {
  text-align: center;
  color: #333;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

input {
  padding: 10px;
  font-size: 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  transition: 0.2s;
}

input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

button {
  padding: 10px;
  font-size: 16px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:disabled {
  background-color: #93c5fd;
  cursor: not-allowed;
}

.message {
  text-align: center;
  font-size: 14px;
  color: #10b981;
  margin-bottom: 15px;
}

.list ul {
  list-style: none;
  padding: 0;
}

.product-card {
  background: #f9fafb;
  padding: 15px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  margin-bottom: 10px;
  transition: 0.3s;
}

.product-card:hover {
  background-color: #f3f4f6;
}

.product-card .price {
  display: block;
  margin-top: 8px;
  font-weight: bold;
  color: #111827;
}
</style>
