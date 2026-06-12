<script setup>
import { ref } from 'vue'

const status = ref('Ready')
const isLoading = ref(false)

async function wake() {
  isLoading.value = true
  status.value = 'Sending Wake-on-LAN...'
  
  try {
    const response = await fetch('/api/wake', {
      method: 'POST'
    })
    
    if (response.ok) {
      const text = await response.text()
      status.value = text
    } else {
      const error = await response.text()
      status.value = `Error: ${error}`
    }
  } catch (err) {
    status.value = `Failed: ${err.message}`
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="container">
    <h1>Wake-on-LAN</h1>
    <p class="subtitle">Send a magic packet to wake up a device on your local network</p>
    
    <button @click="wake" :disabled="isLoading" class="wake-button">
      <span v-if="isLoading">Sending...</span>
      <span v-else>Wake Up Device</span>
    </button>
    
    <div class="status" :class="{ 'error': status.startsWith('Error') || status.startsWith('Failed') }">
      {{ status }}
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 600px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

h1 {
  color: #333;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #666;
  margin-bottom: 2rem;
}

.wake-button {
  background: #0076ff;
  color: white;
  border: none;
  padding: 12px 24px;
  font-size: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  margin-bottom: 1rem;
}

.wake-button:hover {
  background: #0056cc;
}

.wake-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.status {
  padding: 12px;
  border-radius: 8px;
  background: #f0f0f0;
  color: #333;
}

.status.error {
  background: #ffebee;
  color: #c62828;
}
</style>
