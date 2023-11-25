<template>
  <div id="app" class="white-background">
    <!-- Sidebar Menu -->
    <aside class="sidebar-menu">
      <ul>
        <li>
          <a href="#" class="active">Overview</a>
        </li>
        <li>
          <a href="https://github.com/deggja/netfetch" target="_blank">GitHub</a>
        </li>
      </ul>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <div class="header-content">
          <!-- Main Dashboard Header -->
          <header class="app-header">
            <img src="@/assets/logo.png" alt="Netfetch Logo" class="logo">
          </header>

        <div class="dashboard-title-score">
          <h1 class="white-text">Netfetch Dashboard</h1>
          <!-- Score Display -->
          <div v-if="scanInitiated" class="score-display">
            <svg class="donut-chart" width="100" height="100" viewBox="0 0 42 42">
              <circle class="donut-ring" cx="21" cy="21" r="15.91549430918954" fill="transparent" stroke="#fff" stroke-width="3"></circle>
              <circle class="donut-segment" cx="21" cy="21" r="15.91549430918954" fill="transparent" stroke="#ce4b99" stroke-width="3" stroke-dasharray="70 30" stroke-dashoffset="25"></circle>
              <text x="50%" y="50%" class="donut-score white-text" text-anchor="middle" dy=".3em">{{ netfetchScore !== null ? netfetchScore : 'Calculating...' }}</text>
            </svg>
          </div>
        </div>

        <div class="scan-buttons">
          <button @click="fetchScanResults" class="scan-btn light-blue">Scan Cluster</button>
          <button @click="fetchScanResults" class="scan-btn light-blue">Scan Namespace</button>
        </div>
      </div>

      <!-- Success or Error Message Display -->
      <div v-if="message" :class="{'success-message': message.type === 'success', 'error-message': message.type === 'error'}">
        {{ message.text }}
      </div>

      <!-- Unprotected Pods Table -->
      <section v-if="unprotectedPods.length > 0 && scanInitiated" class="pods-table-section">
        <div v-for="(pods, namespace) in paginatedPods" :key="namespace">
          <h3 @click="toggleNamespace(namespace)">
            {{ namespace }}
            <span class="namespace-toggle-indicator">{{ isNamespaceExpanded(namespace) ? '▲' : '▼' }}</span>
          </h3>
          <table v-show="isNamespaceExpanded(namespace)">
            <thead>
              <tr>
                <th>Pod Name</th>
                <th>Pod IP</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="pod in pods" :key="pod.name">
                <td>{{ pod.name }}</td>
                <td>{{ pod.ip }}</td>
                <td>
                  <button @click="remediate(pod.namespace)" class="remediate-btn light-blue">Remediate</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="pagination-controls">
          <button @click="changePage(-1)" :disabled="currentPage === 1" class="pagination-btn">Previous</button>
          <button @click="changePage(1)" :disabled="currentPage * pageSize >= totalPods" class="pagination-btn">Next</button>
        </div>
      </section>


      <!-- Message for No Missing Policies -->
      <h2 v-else class="white-text">No network policies missing. You are good to go!</h2>
    </main>
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
      currentPage: 1,
      pageSize: 10,
      expandedNamespaces: {},
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
  computed: {
    groupedPods() {
      return this.unprotectedPods.reduce((acc, pod) => {
        if (!acc[pod.namespace]) {
          acc[pod.namespace] = [];
        }
        acc[pod.namespace].push(pod);
        return acc;
        }, {});
      },
      paginatedPods() {
      const start = (this.currentPage - 1) * this.pageSize;
      const namespaces = Object.keys(this.groupedPods);
      const paginatedNamespaces = namespaces.slice(start, start + this.pageSize);
      let paginated = {};
      paginatedNamespaces.forEach(namespace => {
        paginated[namespace] = this.groupedPods[namespace];
      });
      return paginated;
    },
    totalPods() {
      return this.unprotectedPods.length;
    }
  },
  methods: {
    toggleMenu() {
      this.menuVisible = !this.menuVisible;
    },
    async fetchScanResults() {
    this.scanInitiated = true;
    try {
      const response = await axios.get('http://localhost:8080/scan');
      this.scanResults = response.data;
      this.unprotectedPods = [];
      this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
      this.netfetchScore = response.data.Score;
      const namespaces = new Set(this.unprotectedPods.map(pod => pod.namespace));
      namespaces.forEach(namespace => {
        this.expandedNamespaces, namespace, true;
      });
    } catch (error) {
      console.error('Error fetching scan results:', error);
      }
      this.updateExpandedNamespaces();
    },
    updateExpandedNamespaces() {
      const namespaces = new Set(this.unprotectedPods.map(pod => pod.namespace));
      namespaces.forEach(namespace => {
        this.expandedNamespaces, namespace, true;
      });
    },
    parseUnprotectedPods(data) {
      if (!data || !Array.isArray(data)) {
        return [];
      }

      const uniquePods = {};
      data.forEach(podDetail => {
        const [namespace, name, ip] = podDetail.split(' ');
        const key = `${namespace}-${name}-${ip}`;
        if (!uniquePods[key]) {
          uniquePods[key] = { namespace, name, ip };
        }
      });

      return Object.values(uniquePods);
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
    toggleNamespace(namespace) {
      this.expandedNamespaces[namespace] = !this.expandedNamespaces[namespace];
    },
    isNamespaceExpanded(namespace) {
      return !!this.expandedNamespaces[namespace];
    },
    changePage(step) {
      this.currentPage += step;
    },
    mounted() {
    // Call updateExpandedNamespaces on mount to handle initial data
      this.updateExpandedNamespaces();
    },
  },
};
</script>

<style>

/* Main Content Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  color: black;
}

.main-content {
  padding-left: 220px;
}

/* Main app background */
.white-background {
  background-color: #fff;
}

.white-text {
  color: black;
}


/* Header styles */
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 80px;
  padding: 0 20px;
}

/* Dashboard Title and Score */
.dashboard-title-score {
  display: flex;
  align-items: center;
}

.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  z-index: 9;
  padding: 10px 20px;
}

.header {
  display: flex;
  align-items: center;
  padding: 10px;
}

.logo {
  height: 70px;
  width: auto;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 15;
  padding: 10px;
  background-color: #fff;
}

.sidebar-menu {
  width: 180px;
  position: fixed;
  top: 80px;
  bottom: 0;
  left: 0;
  background-color: #fff;
  z-index: 10;
  padding-top: 10px;
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
  width: 80%;
  margin: auto;
}

.pods-table-section {
  background-color: #fff;
  padding: 20px;
}

th, td {
  padding: 12px 15px;
  text-align: left;
}

th {
  background-color: #add8e6;
}

tr:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

/* Style for buttons */
.scan-buttons {
  display: flex;
}

button {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s;
}

.light-blue {
  background-color: #add8e6; /* Light blue background */
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
  margin: 40px auto 20px 0;
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

h3 {
  cursor: pointer;
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
  display: flex;
  justify-content: center;
  align-items: center;
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

/* Pagination */
.pagination-controls {
  display: flex;
  justify-content: center;
}

.pagination-btn {
  margin: 0 10px;
}

.namespace-toggle-indicator {
  font-size: 0.8em; /* smaller size than namespace name */
  margin-left: 5px; /* space between name and indicator */
  cursor: pointer; /* indicates interactiveness */
}

</style>
