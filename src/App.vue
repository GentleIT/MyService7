<template>
  <div class="wrapper">
    <div v-if="!isAuthenticated" class="auth-container">
      <h1>Авторизация</h1>
      <form @submit.prevent="handleLogin" class="auth-form">
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
    </div>

    <div style="display: flex; justify-content: center; align-items: center; height: 200px;">
      <div v-if="isAuthenticated" style="margin-bottom: 50px; background-color: #f9f9f9; border: 1px solid #ddd; width: 200px; height: 100px; display: flex; flex-direction: column; justify-content: center;">
        <button @click="getWeather" style="background-color:#3498db; width: 150px; align-self: center">Get Weather</button>
        <p v-if="weather" style="text-align: center;">{{ weather }}</p>
      </div>
    </div>

    <div v-if="isAuthenticated" class="addEquipment-container">
      <div class="addEquipment-items"> <!-- Для центрирорования -->
        <h2>Добавить технику</h2>

        <form @submit.prevent="addEquipment">
          <div class="addEquipment-form-item">
            <label for="name">Название:</label>
            <input type="text" v-model="newEquipment.name" id="name" required />
          </div>
          <div class="addEquipment-form-item">
            <label for="driver">Водитель:</label>
            <input type="text" v-model="newEquipment.driver" id="driver" required />
          </div>
          <div class="addEquipment-form-item">
            <label for="parked">Припаркован:</label>
            <select v-model="newEquipment.parked" id="parked" required>
              <option :value="true">Да</option>
              <option :value="false">Нет</option>
            </select>
          </div>
          <button type="submit">Добавить</button>
        </form>

        <div v-if="responseMessage" class="response-message">
          {{ responseMessage }}
        </div>
      </div>
    </div>

    <div v-if="isAuthenticated" class="getEquipment-form">
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
      loading: false,
      newEquipment: {
        name: "",
        driver: "",
        parked: false
      },
      responseMessage: "",

      weather: ""
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
    handleLogin() {
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
    },

    addEquipment() {
      axios.post('http://localhost:8080/equipment', this.newEquipment, {
        headers: {
          'Content-Type': 'application/json',
          'token': this.token
        }
      })
      .then(response => {
        if (response.data === "Successfully logged into AddEquipment") {
          this.responseMessage = 'Техника успешно добавлена!';
          this.newEquipment = { name: '', driver: '', parked: false }; 
        } else {
          this.responseMessage = 'Ошибка при добавлении техники!';
        }
      })
      .catch(error => {
        this.responseMessage = 'Ошибка при добавлении техники!';
        console.error(error);
      });
    },

    getWeather() {
    axios.get('http://localhost:8080/getWeather')
      .then(response => {
        this.weather = response.data.weather;
      })
      .catch(error => {
        console.error("Ошибка при получении погоды:", error);
      });
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
  flex-direction: column;
}

.auth-container h1{
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
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

.getEquipment-form button{
  height: 40px;
  margin-left: 5px;
  border-radius: 5px;
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

.response-message {
  margin-top: 15px;
  font-size: 18px;
  font-weight: bold;
  color: green;
}

.addEquipment-container {
  background-color: #f4f4f4;
  border: 1px solid #ddd;
  border-radius: 25px;
  height: 200px;
  width: 400px;
  margin:auto;
  margin-bottom: 50px;

  display: flex;
  align-items: center;
  justify-content: center;
}

.addEquipment-items {
  background-color: #3498db;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
}

.addEquipment-items form {
  width: 300px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 0px 0px 10px 10px;
  background-color: #f9f9f9;
}

.addEquipment-form-item {
  display: flex;
  justify-content: space-between;
  margin-top: 4px;
}

.addEquipment-items button {
  width: 100px;
}

.error {
  color: red;
  margin-top: 10px;
}

</style>