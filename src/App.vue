<template>
  <div class="wrapper">
    <!-- <h1>Авторизация</h1>

    <div v-if="!isAuthenticated" class="auth-container">
      <form @submit.prevent="login" class="auth-form">
        <div>
          <label for="login">Логин:</label>
          <input type="text" v-model="login" id="login" required />
        </div>
        <div>
          <label for="password">Пароль:</label>
          <input type="password" v-model="password" id="password" required />
        </div>
        <button type="submit">Войти</button>
      </form>
      <div v-if="error" class="error">{{ error }}</div>
    </div> -->

    <div>
      <h2>Список техники</h2>

      <input 
        type="text" 
        v-model="searchQuery" 
        @input="filterEquipment" 
        placeholder="Поиск по названию..."
        class="search-input"
      />

      <button @click="getEquipment">Загрузить все данные</button>

      <div v-if="loading">Загрузка...</div>
      <div v-if="error" class="error">{{ error }}</div>

      <div class="cards">
        <div v-for="item in filteredEquipment" :key="item.id" class="card">
          <h2>{{ item.name }}</h2>
          <p><strong>Водитель:</strong> {{ item.driver }}</p>
          <p><strong>Дата:</strong> {{ formatDate(item.day) }}</p>
          <p><strong>GPS:</strong> {{ item.gps }}</p>
          <p><strong>Припаркован:</strong> {{ item.parked ? 'Да' : 'Нет' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      login: "",
      password: "",
      isAuthenticated: false,
      token: null,
      error: null,
      equipment: [],
      filteredEquipment: [],
      searchQuery: "",
      loading: false
    };
  },
  created() {
    const savedToken = localStorage.getItem('authToken');
    if (savedToken) {
      this.token = savedToken;
      this.isAuthenticated = true;
    }
  },
  methods: {
    login() {
      console.log('Logging in with:', this.login, this.password);

      axios.post('http://localhost:8080/authorize', {
        login: this.login,
        password: this.password
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      })
      .then(res => {
        console.log('Response:', res.data);
        if (res.data.includes("successfully entered")) {
          this.token = "123456789";
          localStorage.setItem('authToken', this.token);
          this.isAuthenticated = true;
          this.error = null;
        } else {
          this.error = "Неверный логин или пароль.";
        }
      })
      .catch(err => {
        this.error = 'Ошибка авторизации';
        console.error(err);
      });
    },

    getEquipment() {
      this.loading = true;
      this.error = null;

      axios.get('http://localhost:8080/getAll', {
        headers: {
          'Content-Type': 'application/json',
          'token': this.token
        }
      })
      .then(res => {
        this.equipment = res.data;
        this.filteredEquipment = res.data;
        this.loading = false;
      })
      .catch(err => {
        this.error = 'Ошибка загрузки данных';
        console.error(err);
        this.loading = false;
      });
    },

    filterEquipment() {
      this.filteredEquipment = this.equipment.filter(item =>
        item.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },

    formatDate(dateStr) {
      const d = new Date(dateStr);
      return d.toLocaleDateString('ru-RU') + ' ' + d.toLocaleTimeString('ru-RU');
    }
  }
};
</script>

<style scoped>
.wrapper {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  font-family: sans-serif;
  text-align: center;
}

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
  width: 300px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 10px;
  background-color: #f9f9f9;
}

.auth-form input {
  padding: 8px;
  font-size: 16px;
}

.auth-form button {
  padding: 10px;
  font-size: 16px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.search-input {
  margin: 20px 0;
  padding: 10px;
  width: 300px;
}

.cards {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15px;
  margin-top: 20px;
}

.card {
  background-color: #f4f4f4;
  border: 1px solid #ddd;
  border-radius: 10px;
  padding: 15px;
  text-align: left;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>
