<template>
  <div :class="{'app-container': true, 'dark-mode': isDarkMode}">
    <!-- Sidebar Menu -->
    <aside class="sidebar">
      <div class="sidebar-content">
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
        </div>
              <!-- Toggle Button for Dark Mode at the bottom -->
        <div class="dark-mode-toggle-wrapper">
          <button @click="toggleDarkMode" class="dark-mode-toggle" style="width: 50px; height: 50px; background-color: #fff;" title="Dark mode coming soon" disabled>
            <!-- Inline SVG content -->
            <svg v-if ="!isDarkMode" xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 512 256.04" style="fill: currentColor;">
              <path d="M128.02 0h.18c22.03 0 42.83 5.66 61 15.6h210.38c30.89 0 59 12.65 79.38 33.04C499.35 68.99 512 97.1 512 128.02c0 30.92-12.66 59.03-33.02 79.4l-.42.38c-20.34 20.15-48.29 32.64-78.98 32.64H189.24c-18.17 9.93-38.98 15.6-61.04 15.6h-.18c-35.2 0-67.22-14.41-90.42-37.6C14.41 195.25 0 163.24 0 128.02s14.4-67.24 37.59-90.43l.91-.83C61.65 14.05 93.29 0 128.02 0zm-5.95 54.42c0-1.95.8-3.73 2.08-5 2.74-2.77 7.27-2.76 10.02-.01l.14.16a7.042 7.042 0 0 1 1.94 4.85v12.95c0 1.95-.8 3.73-2.08 5.01-2.75 2.75-7.27 2.75-10.02 0a7.084 7.084 0 0 1-2.08-5.01V54.42zm6.05 31.17c11.72 0 22.32 4.75 30 12.43 7.67 7.68 12.43 18.29 12.43 30 0 11.72-4.75 22.32-12.43 30s-18.28 12.43-30 12.43c-11.72 0-22.32-4.75-30.01-12.43-7.67-7.68-12.43-18.28-12.43-30 0-11.72 4.76-22.32 12.43-30 7.69-7.67 18.3-12.43 30.01-12.43zm-56.33-5.34a7.114 7.114 0 0 1-2.07-5.01c0-3.9 3.18-7.09 7.09-7.09 1.81 0 3.62.69 5 2.07l9.16 9.16a7.065 7.065 0 0 1 2.08 5.01c0 1.8-.7 3.62-2.08 5.01a7.057 7.057 0 0 1-5.01 2.08c-1.8 0-3.61-.7-5-2.07l-9.17-9.16zm-17.28 53.81c-1.95 0-3.73-.8-5-2.08-2.77-2.74-2.76-7.27-.01-10.01l.15-.14a7.04 7.04 0 0 1 4.86-1.94h12.94a7.082 7.082 0 0 1 7.09 7.09c0 1.95-.8 3.73-2.07 5.01a7.099 7.099 0 0 1-5.02 2.07H54.51zm25.82 50.28a7.049 7.049 0 0 1-5 2.07c-3.91 0-7.09-3.16-7.09-7.08 0-1.81.68-3.62 2.07-5.01l9.31-9.29a7.02 7.02 0 0 1 4.86-1.94 7.09 7.09 0 0 1 7.09 7.09c0 1.79-.69 3.6-2.08 4.99l-9.16 9.17zm53.82 17.29c0 1.94-.8 3.73-2.08 5-2.74 2.76-7.27 2.75-10.02 0l-.13-.15a7.033 7.033 0 0 1-1.94-4.85v-12.95c0-1.96.8-3.73 2.07-5.01 2.76-2.75 7.27-2.75 10.03 0a7.1 7.1 0 0 1 2.07 5.01v12.95zm50.28-25.83a7.055 7.055 0 0 1 2.07 5.01c0 3.89-3.18 7.09-7.08 7.09-1.81 0-3.63-.69-5.01-2.07l-9.16-9.16a7.095 7.095 0 0 1-2.07-5.02c0-3.9 3.18-7.09 7.08-7.09 1.8 0 3.61.7 5 2.08l9.17 9.16zm17.29-53.82c1.93 0 3.73.81 5 2.08 2.76 2.75 2.75 7.27 0 10.02l-.15.14a7.098 7.098 0 0 1-4.85 1.94h-12.95c-1.96 0-3.74-.8-5.01-2.08-2.76-2.75-2.76-7.27 0-10.02a7.049 7.049 0 0 1 5.01-2.08h12.95zM175.89 71.7a7.074 7.074 0 0 1 5-2.07c3.9 0 7.1 3.19 7.1 7.09 0 1.81-.69 3.62-2.07 5l-9.32 9.31a7.12 7.12 0 0 1-4.86 1.93c-3.91 0-7.09-3.18-7.09-7.09 0-1.8.7-3.61 2.08-5l9.16-9.17zm34.17-41.87c2.96 2.47 5.81 5.07 8.53 7.8 23.22 23.15 37.63 55.17 37.63 90.39s-14.42 67.23-37.6 90.42a130.2 130.2 0 0 1-8.5 7.77h189.46c26.83 0 51.24-10.91 69.02-28.5l.32-.35c17.79-17.79 28.85-42.35 28.85-69.34 0-26.99-11.06-51.55-28.85-69.35-17.77-17.8-42.33-28.84-69.34-28.84H210.06zm-82.04-14.71h.18c62.09 0 112.89 50.81 112.89 112.9 0 62.1-50.86 112.9-112.89 112.9h-.18c-62.03 0-112.9-50.8-112.9-112.9 0-62.09 50.81-112.9 112.9-112.9z"/>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" shape-rendering="geometricPrecision" text-rendering="geometricPrecision" image-rendering="optimizeQuality" fill-rule="evenodd" clip-rule="evenodd" viewBox="0 0 512 262.86">
              <path fill-rule="nonzero" d="M316.78 16.55h-205.9c-30.5 0-58.22 12.48-78.31 32.57C12.47 69.21 0 96.93 0 127.44c0 30.5 12.47 58.22 32.57 78.31 20.09 20.1 47.81 32.57 78.31 32.57h193.25c21.54 15.43 47.9 24.54 76.26 24.54h.18c36.14 0 69.02-14.79 92.83-38.6 23.8-23.81 38.6-56.67 38.6-92.83 0-36.15-14.78-69.03-38.63-92.8C449.53 14.8 416.67 0 380.57 0h-.18c-23.02 0-44.72 6.02-63.61 16.55zm70.62 97.17.43.09c.82-3.45 2.83-6.19 6.04-8.16 3.2-1.98 6.53-2.57 10.01-1.75l.1-.43c-3.47-.82-6.2-2.83-8.17-6.03-1.98-3.22-2.57-6.55-1.75-10.01l-.43-.1c-.82 3.47-2.83 6.2-6.03 8.18-3.21 1.98-6.55 2.56-10.02 1.74l-.1.43c3.47.82 6.2 2.84 8.18 6.04 1.99 3.19 2.56 6.52 1.74 10zm36.87 16.77.53.12c1.02-4.35 3.55-7.78 7.58-10.26 4.02-2.49 8.2-3.22 12.56-2.19l.13-.53c-4.35-1.03-7.78-3.55-10.26-7.59-2.49-4.03-3.22-8.22-2.2-12.56l-.53-.12c-1.02 4.35-3.55 7.77-7.58 10.26-4.02 2.49-8.21 3.22-12.56 2.19l-.13.53c4.36 1.03 7.78 3.55 10.26 7.58 2.49 4.02 3.22 8.22 2.2 12.57zm-38.79-61.01c-15.69 7.67-26.98 23.26-28.29 41.93-1.96 27.88 19.05 52.06 46.92 54.02 13.23.93 25.64-3.32 35.22-11.02 4.75-3.82 9.66-.45 7.59 4.36-11.33 26.42-38.45 44.04-68.74 41.91-38.29-2.69-67.14-35.91-64.45-74.19C316.3 89.8 347.05 61.67 383.44 62c6.71.06 8.13 4.5 2.04 7.48zm-5.09-53.95h.18c63.75 0 115.91 52.15 115.91 115.9 0 63.75-52.23 115.91-115.91 115.91h-.18c-63.68 0-115.91-52.16-115.91-115.91s52.16-115.9 115.91-115.9z"/>
            </svg>
          </button>
        </div>

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
        <button @click="generateClusterNetworkMap" class="scan-btn">Generate Network Map for Cluster</button>
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
        <h2 v-if="scanInitiated && unprotectedPods.length === 0 && !isShowClusterMap" class="no-policies-message">
          All good!
          No unproteced pods found.
        </h2>
      </div>

    <!-- Loading visualization message --> 
    <div v-if="isLoadingVisualization" class="loading-message">
      {{ loadingMessage }}
    </div>
      
    <!-- Visualization components based on isShowClusterMap -->
    <div v-if="isShowClusterMap && clusterVisualizationData.length > 0">
      <network-policy-visualization
        :key="`cluster-map-${clusterMapGenerationCount}`"
        :clusterData="clusterVisualizationData"
        visualizationType="cluster">
      </network-policy-visualization>
    </div>
      
    <!-- Visualization components for namespace -->
    <div v-if="!isShowClusterMap">
      <div v-for="(vizData, namespace) in namespaceVisualizationData" :key="namespace">
        <network-policy-visualization 
          v-if="vizData && vizData.length > 0"
          :policies="vizData"
          visualizationType="namespace">
        </network-policy-visualization>
      </div>
    </div>
    </main>
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
      isDarkMode: false,
      scanResults: null,
      unprotectedPods: [],
      scanInitiated: false,
      message: null,
      netfetchScore: null,
      menuVisible: false,
      currentPage: 1,
      pageSize: 10,
      clusterMapGenerationCount: 0,
      expandedNamespaces: {},
      selectedNamespace: '',
      allNamespaces: [],
      lastScanType: 'cluster',
      namespaceVisualizationData: {},
      isLoadingVisualization: false,
      isShowClusterMap: false,
      isClusterMapLoading: false,
      clusterVisualizationData: [],
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
    loadingMessage() {
      if (this.isClusterMapLoading) {
        return 'Loading cluster visualization data...';
      }
      if (this.isLoadingVisualization) {
        return 'Loading visualization data...';
      }
      return '';
    }
  },
  methods: {
    toggleDarkMode() {
    this.isDarkMode = !this.isDarkMode;
    },
    toggleMenu() {
      this.menuVisible = !this.menuVisible;
    },
    async fetchScanResults() {
    this.isShowClusterMap = false;
    this.scanInitiated = true;
    this.lastScanType = 'cluster';
    try {
      const response = await axios.get('http://localhost:8080/scan');
      this.scanResults = response.data;
      this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
      this.netfetchScore = response.data.Score;
      this.updateExpandedNamespaces();
      this.namespaceVisualizationData = {};
    } catch (error) {
      console.error('Error fetching scan results:', error);
      }
    },
    async fetchNamespacesWithPolicies() {
      try {
        const response = await axios.get('http://localhost:8080/namespaces-with-policies');
        return response.data.namespaces;
      } catch (error) {
        console.error('Error fetching namespaces with policies:', error);
        return [];
      }
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
          this.showSuccessMessage(namespace);
          // Update the unprotectedPods list by removing the remediated namespace's pods.
          this.unprotectedPods = this.unprotectedPods.filter(pod => pod.namespace !== namespace);

          if (this.lastScanType === 'cluster') {
            // If the last scan was a cluster-wide scan, fetch cluster-wide scan results.
            await this.fetchScanResults();
            this.isShowClusterMap = false;
            this.namespaceVisualizationData = {};
          } else {
            // If the last scan was a namespace-specific scan, fetch scan results for that namespace.
            await this.fetchScanResultsForNamespace();
          }
        } else {
          this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}. Status code: ${response.status}` };
        }
      } catch (error) {
        this.message = { type: 'error', text: `Failed to apply policy to namespace: ${namespace}. Error: ${error.message}` };
        console.error('Error applying policy to', namespace, ':', error);
      }
    },
    showSuccessMessage(namespace) {
    this.message = { type: 'success', text: `Policy successfully applied to namespace: ${namespace}` };
    setTimeout(() => {
      this.message = null;
      }, 10000); // Clear success message after 10s
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

      this.isShowClusterMap = false;
      this.namespaceVisualizationData = {};
      this.clusterVisualizationData = [];
      this.lastScanType = 'namespace';
      this.scanInitiated = true;
      this.isLoadingVisualization = true;

      try {
        const scanResponse = await axios.get(`http://localhost:8080/scan?namespace=${namespace}`);
        this.scanResults = scanResponse.data;
        if (scanResponse.data.UnprotectedPods && scanResponse.data.UnprotectedPods.length > 0) {
          this.unprotectedPods = this.parseUnprotectedPods(scanResponse.data.UnprotectedPods);
          this.netfetchScore = scanResponse.data.Score || null;
        } else {
          this.unprotectedPods = [];
          this.netfetchScore = 42;
        }
        this.updateExpandedNamespaces();

        // Fetch visualization data only for the scanned namespace
        await this.fetchVisualizationDataForNamespaces([namespace]);
      } catch (error) {
        console.error('Error scanning namespace:', namespace, error);
        this.message = { type: 'error', text: `Failed to scan namespace: ${namespace}. Error: ${error.message}` };
      } finally {
        this.isLoadingVisualization = false;
      }
    },
    // Fetch and update visualization data for multiple namespaces
    async fetchVisualizationDataForNamespaces(namespaces) {
      this.isLoadingVisualization = true;
      if (!Array.isArray(namespaces)) {
        console.error('Invalid namespaces array:', namespaces);
        return;
      }
      for (const namespace of namespaces) {
        try {
          const response = await axios.get(`http://localhost:8080/visualization?namespace=${namespace}`);
          if (response.data && Array.isArray(response.data.policies)) {
            this.namespaceVisualizationData[namespace] = response.data.policies;
          } else {
            this.namespaceVisualizationData, namespace, [];
          }
        } catch (error) {
          console.error(`Error fetching visualization data for namespace ${namespace}:`, error);
          this.namespaceVisualizationData, namespace, [];
        }
      }
      this.isLoadingVisualization = false;
    },
    async generateClusterNetworkMap() {
      this.isShowClusterMap = true;
      this.isLoadingVisualization = true;
      this.namespaceVisualizationData = {};
      this.clusterMapGenerationCount++;

      try {
      const response = await axios.get('http://localhost:8080/visualization/cluster');
      if (response.data && Array.isArray(response.data)) {
        this.clusterVisualizationData = response.data;
      } else {
        console.error('Received data is not an array:', response.data);
        this.clusterVisualizationData = []; // Set to an empty array if data is not correct
      }
    } catch (error) {
      console.error('Error fetching cluster visualization data:', error);
    } finally {
      this.isLoadingVisualization = false;
      }
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
  background-color: white;
  color: #333;
}

.app-container {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 100vh;
  width: 200px;
  border-right: 1px solid #87CEEB;
  padding: 20px;
  flex: 0 0 200px;
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
  border: 8px solid #87CEEB;
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
  border: 1px solid #87CEEB;
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
  background-color: #5cb7db;
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

.loading-message {
text-align: center;
padding: 20px;
font-size: 1.2em;
color: black;
}

.network-policy-visualization {
  width: 100vw;
  height: 60vh;
  margin: 0 auto;
  padding-top: 5px;
  border: 1px solid #87CEEB;
  border-radius: 5px;
  background-color: #fff;
   display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

/* Dark Mode Specific Styles */
.dark-mode-toggle-container {
  position: absolute;
  bottom: 1px;
  left: 20px;
  margin-top: auto;
}

.dark-mode-toggle {
  background: transparent;
  border: none;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  margin-top: 20px;
  width: 50px;
  height: 50px;
  transition: background-color 0.3s;
}

.dark-mode-toggle svg {
  fill: currentColor;
  width: 30px;
  height: 30px; 
}

.app-container.dark-mode .dark-mode-toggle {
  background-color: #444;
  color: #fff;
}

.app-container.dark-mode .dark-mode-toggle svg {
  fill: #fff;
}

.dark-mode-toggle:hover {
  background-color: #ddd;
  color: #333;
}

.app-container.dark-mode .dark-mode-toggle:hover svg {
  fill: #555;
}

/* Dark Mode Specific Styles */
.app-container.dark-mode {
  background-color: #000; /* Black background for the entire app */
  color: #fff; /* White text color for the entire app */
}

.app-container.dark-mode .sidebar {
  background-color: #1a1a1a; /* Slightly lighter black for sidebar */
  color: #fff; /* White text for sidebar */
}

.app-container.dark-mode .menu-item a {
  color: #fff; /* White text for menu items */
}

.app-container.dark-mode .menu-item a.active {
  background-color: #333; /* Darker background for active menu item */
}

.app-container.dark-mode .content {
  color: #fff; /* White text for main content */
}

.app-container.dark-mode .header {
  color: #fff; /* White text for header */
}

.app-container.dark-mode .score-container {
  border: 5px solid #444; /* Dark border for score container */
}

.app-container.dark-mode .buttons .scan-btn {
  border: 2px solid #444; /* Dark border for buttons */
  color: #fff; /* White text for buttons */
}

.app-container.dark-mode .buttons .scan-btn:hover {
  background-color: #333; /* Darker background on hover for buttons */
}

.app-container.dark-mode .message-container .success-message,
.app-container.dark-mode .message-container .error-message {
  color: #fff; /* White text for messages */
}

.app-container.dark-mode .table-container {
  background-color: #1a1a1a; /* Slightly lighter black for table container */
  color: #fff; /* White text for table */
}

.app-container.dark-mode .pods-table th,
.app-container.dark-mode .pods-table td {
  border: 1px solid #444; /* Dark border for table cells */
}

.app-container.dark-mode .pods-table th {
  background-color: #333; /* Dark background for table header */
}

.app-container.dark-mode .remediate-btn {
  background-color: #333; /* Dark background for remediate button */
}

.app-container.dark-mode .pagination-controls .pagination-btn {
  border: 1px solid #444; /* Dark border for pagination buttons */
}

.app-container.dark-mode .network-policy-visualization {
  border: 2px solid #444; /* Dark border for visualization */
  background-color: #1a1a1a; /* Slightly lighter black for visualization */
}

/* Adjust toggle button appearance in dark mode */
.app-container.dark-mode .dark-mode-toggle {
  background-color: #333; /* Dark background for toggle button */
  color: #fff; /* White icon color */
}

.app-container.dark-mode .dark-mode-toggle:hover {
  background-color: #444; /* Slightly lighter black on hover for toggle button */
}

</style>

