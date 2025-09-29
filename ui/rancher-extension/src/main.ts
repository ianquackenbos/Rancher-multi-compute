import { createApp } from 'vue';
import App from './App.vue';
import { createDataProvider } from './providers';

// Create data provider based on environment
const useMock = import.meta.env.VITE_USE_MOCK_PROVIDER !== 'false';
const dataProvider = createDataProvider(useMock);

// Create Vue app
const app = createApp(App);

// Provide data provider to all components
app.provide('dataProvider', dataProvider);

// Mount the app
app.mount('#app');
