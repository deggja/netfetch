<template>
  <div id="app" class="app-container">
    <!-- Sidebar Menu -->
    <aside class="sidebar">
      <div class="logo-container">
        <img src="@/assets/logo.png" alt="Netfetch Logo" class="logo">
      </div>
      <ul class="menu">
        <li class="menu-item">
          <a href="#" class="active">Overview</a>
        </li>
        <li class="menu-item">
          <a href="https://github.com/deggja/netfetch" target="_blank">GitHub</a>
        </li>
      </ul>
    </aside>

    <!-- Main Content -->
    <main class="content">
      <div class="header">
        <h1 class="dashboard-title">Netfetch Dashboard</h1>
        <div class="score-container" v-if="scanInitiated">
          <span class="score">{{ netfetchScore !== null ? netfetchScore : '...' }}</span>
          <span class="score-label">Score</span>
        </div>
      </div>

      <div class="buttons">
        <button @click="fetchScanResults" class="scan-btn">Scan cluster</button>
        <button @click="fetchScanResultsForNamespace" class="scan-btn">Scan namespace</button>
        <select v-model="selectedNamespace" class="namespace-select">
          <option disabled value="">Select a namespace</option>
          <option v-for="namespace in allNamespaces" :key="namespace" :value="namespace">
            {{ namespace }}
          </option>
        </select>
      </div>

      <div class="message-container">
        <!-- Success or Error Message Display -->
        <div v-if="message" :class="{'success-message': message.type === 'success', 'error-message': message.type === 'error'}">
          {{ message.text }}
        </div>
      </div>

      <div class="table-container">
        <!-- Unprotected Pods Table -->
        <section v-if="unprotectedPods.length > 0 && scanInitiated">
          <div v-for="(pods, namespace) in paginatedPods" :key="namespace">
            <h3 @click="toggleNamespace(namespace)" class="namespace-header">
              {{ namespace }}
              <span class="namespace-toggle-indicator">{{ isNamespaceExpanded(namespace) ? '▲' : '▼' }}</span>
            </h3>
            <table v-show="isNamespaceExpanded(namespace)" class="pods-table">
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
                    <button @click="remediate(pod.namespace)" class="remediate-btn">Remediate</button>
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
        <h2 v-else class="no-policies-message">No network policies missing. You are good to go!</h2>
      </div>
    </main>
  </div>
  <div v-for="(vizData, namespace) in namespaceVisualizationData" :key="namespace">
  <network-policy-visualization 
    v-if="vizData.length > 0"
    :policies="vizData">
  </network-policy-visualization>
</div>
</template>



<script>
import axios from 'axios';
import NetworkPolicyVisualization from './Viz.vue';

export default {
  name: 'App',
  components: {
    NetworkPolicyVisualization,
  },
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
      selectedNamespace: '',
      allNamespaces: [],
      lastScanType: 'cluster',
      namespaceVisualizationData: {},
    };
  },
  watch: {
  selectedNamespace(newNamespace, oldNamespace) {
    if (newNamespace !== oldNamespace) {
      this.fetchVisualizationData(newNamespace);
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
      const grouped = this.groupedPods;
      const paginated = {};
      let startIndex = 0;
      let namespaceIndex = 0;

      let podsCounter = 0;
      for (const ns in grouped) {
        const podsInNamespace = grouped[ns].length;
        if (podsCounter + podsInNamespace >= (this.currentPage - 1) * this.pageSize) {
          startIndex = (this.currentPage - 1) * this.pageSize - podsCounter;
          namespaceIndex = Object.keys(grouped).indexOf(ns);
          break;
        }
        podsCounter += podsInNamespace;
      }

      let podsAdded = 0;
      while (podsAdded < this.pageSize && namespaceIndex < Object.keys(grouped).length) {
        const currentNamespace = Object.keys(grouped)[namespaceIndex];
        const podsInNamespace = grouped[currentNamespace].slice(startIndex, startIndex + this.pageSize - podsAdded);
        paginated[currentNamespace] = podsInNamespace;
        podsAdded += podsInNamespace.length;

        startIndex = 0;
        namespaceIndex++;
      }

      return paginated;
    },
    totalPods() {
    return this.unprotectedPods.length;
    },
    totalPages() {
    return Math.ceil(this.totalPods / this.pageSize);
    },
    startIndexForCurrentPage() {
      let count = 0;
      for (let i = 0; i < this.unprotectedPods.length; i++) {
        if (count >= (this.currentPage - 1) * this.pageSize) {
          return i;
        }
        count += this.unprotectedPods[i].length;
      }
      return 0;
    },
  },
  methods: {
    toggleMenu() {
      this.menuVisible = !this.menuVisible;
    },
    async fetchScanResults() {
    this.scanInitiated = true;
    this.lastScanType = 'cluster';
    try {
      const response = await axios.get('http://localhost:8080/scan');
      this.scanResults = response.data;
      this.unprotectedPods = [];
      this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
      this.netfetchScore = response.data.Score;
      this.updateExpandedNamespaces();
      const namespaces = new Set(this.unprotectedPods.map(pod => pod.namespace));
      this.fetchVisualizationDataForNamespaces(Array.from(namespaces));
      namespaces.forEach(namespace => {
        this.expandedNamespaces, namespace, true;
      });
      this.fetchVisualizationData('');
      this.fetchVisualizationDataForNamespaces(this.unprotectedPods.map(pod => pod.namespace));
    } catch (error) {
      console.error('Error fetching scan results:', error);
      }
      this.updateExpandedNamespaces();
    },
    updateExpandedNamespaces() {
    const namespaces = new Set(this.unprotectedPods.map(pod => pod.namespace));
    namespaces.forEach(namespace => {
      this.expandedNamespaces[namespace] = true;
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
        const response = await axios.post('http://localhost:8080/add-policy', { namespace });
        if (response.status === 200) {
          this.message = { type: 'success', text: `Policy successfully applied to namespace: ${namespace}` };
          this.unprotectedPods = this.unprotectedPods.filter(pod => pod.namespace !== namespace);

          if (this.lastScanType === 'cluster') {
            await this.fetchScanResults();
          } else {
            await this.fetchScanResultsForNamespace(namespace);
          }

          // Refresh the network visualization for the current namespace
          this.fetchVisualizationData(namespace);
        } else {
          this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}. Status code: ${response.status}` };
        }
        
        this.removeVisualizationDataForNamespace(namespace);
            // Re-fetch visualization data for the remaining namespaces
            const remainingNamespaces = this.unprotectedPods
                .filter(pod => pod.namespace !== namespace)
                .map(pod => pod.namespace);
            this.fetchVisualizationDataForNamespaces(remainingNamespaces);
      } catch (error) {
        this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}. Error: ${error.message}` };
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
      this.currentPage = Math.max(1, Math.min(this.currentPage + step, this.totalPages));
    },
    removeVisualizationDataForNamespace(namespace) {
        if (this.visualizationData && this.visualizationData.policies) {
            this.visualizationData.policies = this.visualizationData.policies.filter(policy => {
                return policy.targetPods.some(pod => pod.namespace === namespace);
            });
        }
    },
    async fetchAllNamespaces() {
      console.log('fetchAllNamespaces called');
      try {
        const response = await axios.get('http://localhost:8080/namespaces');
        this.allNamespaces = response.data.namespaces;
        console.log('Namespaces fetched:', this.allNamespaces);
      } catch (error) {
        console.error('Error fetching namespaces:', error);
      }
    },
    async fetchScanResultsForNamespace() {
      const namespace = this.selectedNamespace;
      if (!namespace) {
        alert('Please select a namespace.');
        return;
      }

      this.lastScanType = 'namespace';
      this.scanInitiated = true;
      try {
        const response = await axios.get(`http://localhost:8080/scan?namespace=${namespace}`);
        this.scanResults = response.data;
        if (response.data.UnprotectedPods && response.data.UnprotectedPods.length > 0) {
          this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
          this.netfetchScore = response.data.Score || null;
          this.updateExpandedNamespaces();
        } else {
          this.unprotectedPods = [];
          this.netfetchScore = 42;
        }
        this.fetchVisualizationDataForNamespaces([namespace]);
      } catch (error) {
        console.error('Error scanning namespace:', namespace, error);
        this.message = { type: 'error', text: `Failed to scan namespace: ${namespace}. Error: ${error.message}` };
      }
    },
    // Fetch and update visualization data for multiple namespaces
    fetchVisualizationDataForNamespaces(namespaces) {
      if (!Array.isArray(namespaces)) {
        console.error('Invalid namespaces array:', namespaces);
        return;
      }
      namespaces.forEach(async (namespace) => {
        try {
          const response = await axios.get(`http://localhost:8080/visualization?namespace=${namespace}`);
          console.log(response); // For debugging
   
          if (response.data && Array.isArray(response.data.policies)) {
            this.namespaceVisualizationData, namespace, response.data.policies;
          } else {
            this.namespaceVisualizationData, namespace, [];
          }
        } catch (error) {
          console.error(`Error fetching visualization data for namespace ${namespace}:`, error);
          this.namespaceVisualizationData, namespace, [];
        }
      });
    },
    // Viz
    async fetchVisualizationData(namespace) {
      if (!namespace) return;
      try {
        const response = await axios.get(`http://localhost:8080/visualization?namespace=${namespace}`);
        if (response.data && Array.isArray(response.data.policies)) {
          // Your existing logic
        } else {
          console.warn('Visualization data for namespace is null or not in expected format:', namespace);
        }
      } catch (error) {
        console.error('Error fetching visualization data:', error);
      }
    },
  },
  mounted() {
      this.updateExpandedNamespaces();
      this.fetchAllNamespaces();
  },
};
</script>

<style>
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: 'Helvetica', sans-serif;
  }

  body, html {
    background-color: #f5f5f5;
    color: #333;
  }

  .app-container {
    display: flex;
    min-height: 100vh;
  }

  .sidebar {
    background-color: #fff;
    width: 200px;
    border-right: 5px solid #87CEEB;
    border-width: 1px;
    padding: 20px;
  }

  .logo-container {
    margin-bottom: 20px;
  }

  .menu {
    list-style-type: none;
  }

  .menu-item a {
    display: block;
    padding: 10px;
    color: #333;
    text-decoration: none;
    border-left: 5px solid transparent;
  }

  .menu-item a.active {
    border-left: 5px solid #87CEEB;
    background-color: #eee;
  }

  .menu-item a:hover {
    background-color: #ddd;
  }

  .content {
    flex-grow: 1;
    padding: 40px;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .dashboard-title {
    color: #333;
    font-size: 24px;
  }

  .score-container {
    width: 100px;
    height: 100px;
    border: 5px solid #87CEEB;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    flex-direction: column;
  }

  .score {
    font-size: 24px;
    font-weight: bold;
  }

  .score-label {
  font-size: 0.8em;
  position: absolute;
  bottom: 10px;
  left: 0;
  right: 0;
  text-align: center;
}

  .buttons {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  .scan-btn {
    padding: 10px 20px;
    border: 2px solid #87CEEB;
    border-radius: 5px;
    background-color: #fff;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .namespace-select {
    padding: 10px;
    border-radius: 5px;
    border: 2px solid #87CEEB;
    text-align: center;
    -moz-appearance: none;
    appearance: none;
  }

  .scan-btn:hover {
    background-color:#87CEEB;
    color: #fff;
  }

  .message-container {
    margin-bottom: 20px;
  }

  .success-message {
    background-color: #28a745;
    color: #fff;
    padding: 10px;
    border-radius: 5px;
    text-align: center;
  }

  .error-message {
    background-color: #dc3545;
    color: #fff;
    padding: 10px;
    border-radius: 5px;
    text-align: center;
  }

  /* Table container */
  .table-container {
    background-color: #fff;
    border: 2px solid #87CEEB;
    padding: 20px;
    border-radius: 1px;
    width: 100%;
    overflow-x: auto;
  }

  .namespace-header {
    cursor: pointer;
    margin: 20px 0;
    font-size: 18px;
    min-width: 200px;
  }

  .namespace-toggle-indicator {
    font-size: 0.8em;
    margin-left: 5px;
  }

  .pods-table {
    width: 100%;
    table-layout: fixed;
    border-collapse: collapse;
    margin-bottom: 20px;
  }

  .pods-table th, .pods-table td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }

  .pods-table th {
    background-color: #f0f0f0;
  }

  .pods-table tr:nth-child(even) {
    background-color: #f9f9f9;
  }

  .pods-table tr:hover {
    background-color: #f1f1f1;
  }

  .remediate-btn {
    background-color: #87CEEB;
    color: white;
    padding: 6px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .remediate-btn:hover {
    background-color: #87CEEB;
  }

  .pagination-controls {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }

  .pagination-btn {
    padding: 5px 10px;
    margin: 0 10px;
    background-color: #fff;
    border: 1px solid #87CEEB;
    border-radius: 5px;
    cursor: pointer;
  }

  .pagination-btn:hover {
    background-color: #87CEEB;
  }

  .pagination-btn:disabled {
    color: #aaa;
    cursor: not-allowed;
  }

  .no-policies-message {
    color: black;
    text-align: center;
    margin-top: 20px;
  }
</style>

