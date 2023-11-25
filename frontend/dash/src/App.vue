<template>
  <div id="app">

    <!-- Logo and Hamburger Menu -->
    <header class="app-header">
      <img src="@/assets/logo.png" alt="Netfetch Logo" class="logo">
      <button class="hamburger" @click="toggleMenu">&#9776;</button>
      <nav v-if="menuVisible" class="menu">
        <ul>
          <li><a href="#">Overview</a></li>
        </ul>
      </nav>
    </header>

    <h1>Netfetch Dashboard</h1>
    <button @click="fetchScanResults">Scan Network Policies</button>

    <!-- Score Display -->
    <div v-if="scanInitiated" class="score-display">
      <h2>Your Netfetch Score</h2>
      <div class="donut-chart-container">
        <svg class="donut-chart" width="200" height="200" viewBox="0 0 42 42">
          <circle class="donut-ring" cx="21" cy="21" r="15.91549430918954" fill="transparent" stroke="#d2d3d4" stroke-width="3"></circle>
          <circle class="donut-segment" cx="21" cy="21" r="15.91549430918954" fill="transparent" stroke="#ce4b99" stroke-width="3" stroke-dasharray="0 100" stroke-dashoffset="25"></circle>
          <text x="50%" y="50%" class="donut-score" text-anchor="middle" dy=".3em">{{ netfetchScore !== null ? netfetchScore : 'Calculating...' }}</text>
        </svg>
      </div>
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
      menuVisible: false,
    };
  },
  watch: {
  netfetchScore(newScore) {
    if (newScore !== null) {
      const maxScore = 42;
      const normalizedScore = (newScore / maxScore) * 100;
      const donutSegment = this.$el.querySelector('.donut-segment');
      if (donutSegment) {
        donutSegment.style.strokeDasharray = `${normalizedScore} ${100 - normalizedScore}`;
        }
      }
    },
  },
  methods: {
    toggleMenu() {
      this.menuVisible = !this.menuVisible;
    },
    async fetchScanResults() {
    this.scanInitiated = true;
    try {
      const response = await axios.get('http://localhost:8080/scan');
      console.log("Backend response:", response.data); // Add this line
      this.scanResults = response.data;
      this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
      this.netfetchScore = response.data.Score;
    } catch (error) {
      console.error('Error fetching scan results:', error);
      }
    },
    parseUnprotectedPods(data) {
    if (!data || !Array.isArray(data)) {
      return []; // Return an empty array if data is null or not an array
    }
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

      // Fetch updated scan results to refresh the score and the list of unprotected pods
      const rescanResponse = await axios.get('http://localhost:8080/scan');
      this.netfetchScore = rescanResponse.data.Score;
      this.unprotectedPods = this.parseUnprotectedPods(rescanResponse.data.UnprotectedPods);
    } catch (error) {
      this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}` };
      console.error('Error applying policy to', namespace, ':', error);
      }
    },
  },
};
</script>

<style>

/* Header styles */
.app-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 10px;
  background-color: #fff;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1000;
}

.header {
  display: flex;
  align-items: center;
  padding: 10px;
}

.logo {
  height: 70px;
  margin-bottom: 10px;
}

.hamburger {
  color: #000;
  background: none;
  border: none;
  font-size: 34px;
  cursor: pointer;
  margin-left: 5px;
}

.menu {
  display: none;
  position: absolute;
  background-color: none;
  min-width: 160px;
  box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
  z-index: 1;
}

.menu ul {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.menu li a {
  padding: 12px 16px;
  text-decoration: none;
  display: block;
  color: black;
}

.menu li a:hover {
  background-color: #ddd;
}

/* Table Styles */
table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px 15px;
  text-align: left;
}

th {
  background-color: rgba(0, 123, 255, 0.5);
  color: white;
  font-weight: bold;
}

tr:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

/* Style for buttons */
button {
  background-color: #28a745;
  color: 000;
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #218838;
}

/* General Styles */
body, html {
  margin: 0;
  color: 000;
  padding: 0;
  background-color: #fff;
}

a {
  color: #000;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

#app {
  max-width: 1200px;
  margin: 40px auto;
  padding: 20px;
  padding-top: 60px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}

h1 {
  color: #333;
}

h2 {
  color: white;
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

/* Score Display Styles */
.score-display {
  /* Remove background and box-shadow if not needed */
  margin: 20px auto;
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
}

.score {
  font-size: 24px;
  font-weight: bold;
  color: white;
}

.card {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
}

/* Success and Error Message Styles */
.success-message, .error-message {
  padding: 10px;
  border-radius: 5px;
  text-align: center;
  margin-top: 20px;
}

.success-message {
  background-color: #28a745;
}

.error-message {
  background-color: #dc3545;
}

/* Style the Donut */
.donut-chart-container {
  position: relative;
  width: 200px;
  height: 200px;
}

.donut-chart {
  transform: rotate(-90deg);
}

.donut-ring {
  stroke-width: 3;
}

.donut-segment {
  transition: stroke-dasharray 0.3s;
}

.donut-score {
  fill: #000;
  font-size: 0.3em;
  font-weight: bold;
  text-anchor: middle;
  transform: rotate(90deg);
  transform-origin: center
}

/* Ensures all text is black */
body, html, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td,
article, aside, canvas, details, embed, 
figure, figcaption, footer, header, hgroup, 
menu, nav, output, ruby, section, summary,
time, mark, audio, video {
  color: black;
  border: none;
}

</style>
