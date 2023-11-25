<template>
  <div id="app">
    <h1>Netfetch Dashboard</h1>
    <button @click="fetchScanResults">Scan Network Policies</button>

    <!-- Score Display -->
    <div v-if="netfetchScore !== null" class="score-display">
      <h2>Your Netfetch Score</h2>
      <p class="score">{{ netfetchScore }} out of 42</p>
    </div>

    <!-- Display Success or Error Message -->
    <div v-if="message" :class="{ 'success-message': message.type === 'success', 'error-message': message.type === 'error' }">
      {{ message.text }}
    </div>

    <!-- Message when no unprotected pods are found -->
    <div v-if="unprotectedPods.length === 0 && scanInitiated">
      <h2>No network policies missing. You are good to go!</h2>
    </div>

    <!-- Table for Unprotected Pods -->
    <div v-else-if="unprotectedPods.length > 0">
      <h2>Unprotected Pods</h2>
      <table>
        <tr>
          <th>Namespace</th>
          <th>Pod Name</th>
          <th>Pod IP</th>
          <th>Action</th>
        </tr>
        <tr v-for="pod in unprotectedPods" :key="pod.name">
          <td>{{ pod.namespace }}</td>
          <td>{{ pod.name }}</td>
          <td>{{ pod.ip }}</td>
          <td>
            <button @click="remediate(pod.namespace)" class="remediate-btn">
              Remediate
              <span class="tooltip">Apply a default deny all network policy to this namespace</span>
            </button>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      scanResults: null,
      unprotectedPods: [],
      scanInitiated: false,
      message: null,
      netfetchScore: null,
    };
  },
  methods: {
    async fetchScanResults() {
      this.scanInitiated = true;
      try {
        const response = await axios.get('http://localhost:8080/scan');
        this.scanResults = response.data;
        this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
        this.netfetchScore = response.data.Score;
      } catch (error) {
        console.error('Error fetching scan results:', error);
      }
    },
    parseUnprotectedPods(data) {
      return data.map(podDetail => {
        const [namespace, name, ip] = podDetail.split(' ');
        return { namespace, name, ip };
      });
    },
    async remediate(namespace) {
    try {
      await axios.post('http://localhost:8080/add-policy', { namespace });
      this.message = { type: 'success', text: `Policy successfully applied to namespace: ${namespace}` };
      this.unprotectedPods = this.unprotectedPods.filter(pod => pod.namespace !== namespace);
      await this.fetchScanResults();
    } catch (error) {
      this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}` };
      console.error('Error applying policy to', namespace, ':', error);
      }
    }
  },
};
</script>

<style>
/* Table Styles */
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
  transition: 0.3s;
  border-radius: 5px;
}

th, td {
  border-bottom: 1px solid #ddd;
  padding: 12px 15px;
  text-align: left;
  color: #333;
}

th {
  background-color: #4CAF50;
  color: white;
  font-weight: bold;
}

tr:hover {
  background-color: #f5f5f5;
}

/* Button Styles */
button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  margin: 10px 0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: 0.3s;
}

button:hover {
  background-color: #45a049;
}

/* General Styles */
#app {
  font-family: Arial, sans-serif;
  padding: 20px;
}

h1 {
  color: #333;
}

h2 {
  color: #444;
  margin-top: 30px;
}

.remediate-btn {
  background-color: #ff9800; /* Orange color for attention */
  color: white;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  position: relative; /* Required for tooltip positioning */
}

.remediate-btn:hover {
  background-color: #e68a00; /* Darker shade for hover effect */
}

.tooltip {
  visibility: hidden;
  width: 220px;
  background-color: black;
  color: #fff;
  text-align: center;
  border-radius: 6px;
  padding: 5px;
  position: absolute;
  z-index: 1;
  bottom: 125%;
  left: 50%;
  margin-left: -110px; /* Half of the width to align it */
  opacity: 0;
  transition: opacity 0.3s
}

.remediate-btn:hover .tooltip {
  visibility: visible;
  opacity: 1;
}

.score-display {
  background-color: #4CAF50;
  color: white;
  padding: 20px;
  margin-top: 20px;
  border-radius: 10px;
  text-align: center;
}

.score {
  font-size: 24px;
  font-weight: bold;
}
</style>
