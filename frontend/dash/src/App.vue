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
          <div class="left-buttons">
          <button @click="fetchScanResults" class="scan-btn">Scan cluster</button>
          <button @click="fetchScanResultsForNamespace" class="scan-btn">Scan namespace</button>
          <select v-model="selectedNamespace" class="namespace-select">
            <option disabled value="">Select a namespace</option>
            <option v-for="namespace in allNamespaces" :key="namespace" :value="namespace">
              {{ namespace }}
            </option>
          </select>
          <button @click="suggestPolicy" class="scan-btn">Suggest policy</button>
          </div>
          <button @click="generateClusterNetworkMap" class="scan-btn create-cluster-map-btn">Create cluster map</button>
        </div>

        <!-- Checkbox for selecting "native" network policy type -->
        <div class="policy-type-checkboxes">
          <label class="checkbox-container">
            kubernetes
            <input type="checkbox" v-model="isScanForNative" checked disabled>
            <span class="checkmark"></span>
          </label>
          <label class="checkbox-container">
            cilium
            <input type="checkbox" v-model="isScanForCilium" disabled>
            <span class="checkmark"></span>
          </label>
          <label class="checkbox-container">
            calico
            <input type="checkbox" v-model="isScanForCalico" disabled>
            <span class="checkmark"></span>
          </label>
        </div>

        <div class="message-container">
          <!-- Success or Error Message Display -->
          <div v-if="message" :class="{'notification success': message.type === 'success', 'error': message.type === 'notification error'}">
            {{ message.text }}
          </div>
        </div>
        <section v-if="suggestedNetworkPolicies.length > 0" class="policy-container">
          <h3 class="policy-header">
            Suggested network policies for {{ activeNamespaceForPolicies || 'cert-manager' }}
              <span class="tooltip-icon" @mouseover="showTooltip" @mouseleave="hideTooltip">
                <!-- New SVG inline -->
                <svg class="svg-icon" style="width: 1em; height: 1em; vertical-align: middle; fill: currentColor; overflow: hidden;" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg">
                  <path d="M512 81.408a422.4 422.4 0 1 0 422.4 422.4A422.4 422.4 0 0 0 512 81.408z m26.624 629.76a45.056 45.056 0 0 1-31.232 12.288 42.496 42.496 0 0 1-31.232-12.8 41.984 41.984 0 0 1-12.8-30.72 39.424 39.424 0 0 1 12.8-30.72 42.496 42.496 0 0 1 31.232-12.288 43.008 43.008 0 0 1 31.744 12.288 39.424 39.424 0 0 1 12.8 30.72 43.008 43.008 0 0 1-13.312 31.744z m87.04-235.52a617.472 617.472 0 0 1-51.2 47.104 93.184 93.184 0 0 0-25.088 31.232 80.896 80.896 0 0 0-9.728 39.936v10.24h-64v-10.24a119.808 119.808 0 0 1 12.288-57.344A311.296 311.296 0 0 1 555.52 460.8l10.24-11.264a71.168 71.168 0 0 0 16.896-44.032A69.632 69.632 0 0 0 563.2 358.4a69.632 69.632 0 0 0-51.2-17.92 67.072 67.072 0 0 0-58.88 26.112 102.4 102.4 0 0 0-16.384 61.44h-61.44a140.288 140.288 0 0 1 37.888-102.4 140.8 140.8 0 0 1 104.96-38.4 135.68 135.68 0 0 1 96.256 29.184 108.032 108.032 0 0 1 36.352 86.528 116.736 116.736 0 0 1-25.088 73.216z"/>
                </svg>
                <span class="tooltiptext" v-show="showTooltip">Important note: These suggestions are generated by analyzing your existing applications. You should inspect the YAML and make the neccessary adjustments. It is possible to edit the policy inline before applying them to the namespace by using edit.</span>
              </span>
            </h3>
            <div v-if="creationMessage" class="notification success">{{ creationMessage }}</div>
            <div v-if="creationError" class="notification error">{{ creationError }}</div>
          <div class="policy-cards-container">
            <div class="policy-card" v-for="(policy, index) in displayedPolicies" :key="index">
              <!-- Display textarea if editIndex matches, otherwise display policy -->
              <textarea v-if="editIndex === index" v-model="editablePolicy" class="policy-edit" ref="textarea"></textarea>
              <pre v-else class="policy-yaml">{{ policy }}</pre>
              <div class="copy-container">
                <button @click="copyToClipboard(policy, index, $event)" class="button-copy" :title="copyButtonTexts[index] || 'Copy'">
                  <!-- Inline SVG content -->
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 111.07 122.88" class="copy-icon" ref="copy-icon-${index}">
                    <path class="st0" d="M97.67,20.81L97.67,20.81l0.01,0.02c3.7,0.01,7.04,1.51,9.46,3.93c2.4,2.41,3.9,5.74,3.9,9.42h0.02v0.02v75.28 v0.01h-0.02c-0.01,3.68-1.51,7.03-3.93,9.46c-2.41,2.4-5.74,3.9-9.42,3.9v0.02h-0.02H38.48h-0.01v-0.02 c-3.69-0.01-7.04-1.5-9.46-3.93c-2.4-2.41-3.9-5.74-3.91-9.42H25.1c0-25.96,0-49.34,0-75.3v-0.01h0.02 c0.01-3.69,1.52-7.04,3.94-9.46c2.41-2.4,5.73-3.9,9.42-3.91v-0.02h0.02C58.22,20.81,77.95,20.81,97.67,20.81L97.67,20.81z M0.02,75.38L0,13.39v-0.01h0.02c0.01-3.69,1.52-7.04,3.93-9.46c2.41-2.4,5.74-3.9,9.42-3.91V0h0.02h59.19 c7.69,0,8.9,9.96,0.01,10.16H13.4h-0.02v-0.02c-0.88,0-1.68,0.37-2.27,0.97c-0.59,0.58-0.96,1.4-0.96,2.27h0.02v0.01v3.17 c0,19.61,0,39.21,0,58.81C10.17,83.63,0.02,84.09,0.02,75.38L0.02,75.38z M100.91,109.49V34.2v-0.02h0.02 c0-0.87-0.37-1.68-0.97-2.27c-0.59-0.58-1.4-0.96-2.28-0.96v0.02h-0.01H38.48h-0.02v-0.02c-0.88,0-1.68,0.38-2.27,0.97 c-0.59,0.58-0.96,1.4-0.96,2.27h0.02v0.01v75.28v0.02h-0.02c0,0.88,0.38,1.68,0.97,2.27c0.59,0.59,1.4,0.96,2.27,0.96v-0.02h0.01 h59.19h0.02v0.02c0.87,0,1.68-0.38,2.27-0.97c0.59-0.58,0.96-1.4,0.96-2.27L100.91,109.49L100.91,109.49L100.91,109.49 L100.91,109.49z"/>
                  </svg>
                </button>
              </div>
              <div class="action-buttons">
                <button v-if="editIndex === index" @click="savePolicyChanges">
                  Save
                </button>
                <button v-if="editIndex === index" @click="cancelEdit">
                  Cancel
                </button>
                <button v-else @click="editPolicy(index)">
                  Edit
                </button> 
                <button @click="createPolicy(policy, selectedNamespace)">
                  Create
                </button>
              </div>
            </div>
          </div>
          <div class="policy-pagination">
            <button @click="previousPolicySet" :disabled="currentPolicySetStart === 0" class="pagination-btn">Previous</button>
            <button @click="nextPolicySet" :disabled="currentPolicySetStart + policySetSize >= suggestedNetworkPolicies.length" class="pagination-btn">Next</button>
          </div>
        </section>
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
                    <th>Action
                      <span class="tooltip-icon" @mouseover="showRemediateTooltip" @mouseleave="hideTooltip">
                        <svg class="svg-icon" style="width: 1em; height: 1em; vertical-align: middle; fill: currentColor; overflow: hidden;" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg">
                          <path d="M512 81.408a422.4 422.4 0 1 0 422.4 422.4A422.4 422.4 0 0 0 512 81.408z m26.624 629.76a45.056 45.056 0 0 1-31.232 12.288 42.496 42.496 0 0 1-31.232-12.8 41.984 41.984 0 0 1-12.8-30.72 39.424 39.424 0 0 1 12.8-30.72 42.496 42.496 0 0 1 31.232-12.288 43.008 43.008 0 0 1 31.744 12.288 39.424 39.424 0 0 1 12.8 30.72 43.008 43.008 0 0 1-13.312 31.744z m87.04-235.52a617.472 617.472 0 0 1-51.2 47.104 93.184 93.184 0 0 0-25.088 31.232 80.896 80.896 0 0 0-9.728 39.936v10.24h-64v-10.24a119.808 119.808 0 0 1 12.288-57.344A311.296 311.296 0 0 1 555.52 460.8l10.24-11.264a71.168 71.168 0 0 0 16.896-44.032A69.632 69.632 0 0 0 563.2 358.4a69.632 69.632 0 0 0-51.2-17.92 67.072 67.072 0 0 0-58.88 26.112 102.4 102.4 0 0 0-16.384 61.44h-61.44a140.288 140.288 0 0 1 37.888-102.4 140.8 140.8 0 0 1 104.96-38.4 135.68 135.68 0 0 1 96.256 29.184 108.032 108.032 0 0 1 36.352 86.528 116.736 116.736 0 0 1-25.088 73.216z"/>
                        </svg>
                        <span class="tooltiptext">{{ remediateTooltipText }}</span>
                      </span>
                    </th>
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
            Scan completed
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
  import yaml from 'js-yaml';

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
        isShowScanCompletedMessage: false,
        clusterVisualizationData: [],
        suggestedNetworkPolicies: [],
        currentPolicyIndex: 0,
        policySetSize: 3,
        currentPolicySetStart: 0,
        copyButtonTexts: {},
        editButtonTexts: {},
        editMode: {},
        editablePolicies: {},
        editablePolicy: '',
        editIndex: null,
        creationMessage: '',
        creationError: '',
        isShowTooltip: false,
        activeNamespaceForPolicies: '',
        isScanForNative: true,
        remediateTooltipText: 'Remediate will create a default deny all ingress and egress network policy in the namespace. This will deny all traffic coming to and from the pods. In addition to doing this, you must create network policies to allow the required traffic from and to your pods. You can do this by using the Suggest policy button.'
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
      displayedPolicies() {
        return this.suggestedNetworkPolicies.slice(this.currentPolicySetStart, this.currentPolicySetStart + this.policySetSize);
      },
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
      showTooltip() {
        this.isShowTooltip = true;
      },
      hideTooltip() {
        this.isShowTooltip = false;
      },
      editPolicy(index) {
        this.editIndex = index;
        this.editablePolicy = this.suggestedNetworkPolicies[index];
      },  
      // Called when the "Save" button is clicked
      savePolicyChanges() {
        if (this.editIndex !== null) {
          this.suggestedNetworkPolicies[this.editIndex] = this.editablePolicy;
          
          this.editIndex = null;
          this.editablePolicy = '';
        }
      },
      cancelEdit() {
        this.editIndex = null;
        this.editablePolicy = '';
      },
      async createPolicy(policyYaml, namespace) {
        try {
          const response = await axios.post('/create-policy', {
            yaml: policyYaml,
            namespace: namespace,
          });
          if (response.status === 200) {
            this.creationMessage = 'Policy created successfully!';
            this.creationError = '';

            // Reset visualization data for all other namespaces
            Object.keys(this.namespaceVisualizationData).forEach(ns => {
              if (ns !== namespace) {
                // Assuming the default state for no visualization data is an empty array or null
                this.namespaceVisualizationData[ns] = [];
              }
            });
            
            // Refresh the suggested policies and visualization data for the namespace
            await this.suggestPolicy();
            await this.fetchVisualizationDataForNamespaces([namespace]);
          } else {
            this.creationError = `Policy creation failed with status: ${response.status}`;
            this.creationMessage = '';
          }
        } catch (error) {
          this.creationError = 'Policy creation failed. ' + (error.response?.data || error.message);
          this.creationMessage = ''; // Clear any previous messages
        }
      },
      copyToClipboard(policyYaml, index, event) {
        navigator.clipboard.writeText(policyYaml)
          .then(() => {
            // Use the event target to find the closest button ancestor in case the SVG is not a direct child
            const button = event.target.closest('button');
            const svgIcon = button ? button.querySelector('svg') : null;
            if (svgIcon) {
              svgIcon.classList.add('copy-success');
              setTimeout(() => {
                svgIcon.classList.remove('copy-success');
              }, 1000);
            }
          })
          .catch(err => {
            console.error('Failed to copy text: ', err);
          });
      },
      nextPolicySet() {
      const nextStart = this.currentPolicySetStart + this.policySetSize;
      if (nextStart < this.suggestedNetworkPolicies.length) {
        this.currentPolicySetStart = nextStart;
        }
      },
      previousPolicySet() {
        const prevStart = this.currentPolicySetStart - this.policySetSize;
        if (prevStart >= 0) {
          this.currentPolicySetStart = prevStart;
        }
      },
      toggleDarkMode() {
      this.isDarkMode = !this.isDarkMode;
      },
      toggleMenu() {
        this.menuVisible = !this.menuVisible;
      },
      async fetchScanResults() {
        if (this.isScanForNative) {
          this.isShowClusterMap = false;
          this.scanInitiated = true;
          this.lastScanType = 'cluster';
          this.suggestedNetworkPolicies = [];
          try {
            const response = await axios.get('/scan');
            this.scanResults = response.data;
            this.unprotectedPods = this.parseUnprotectedPods(response.data.UnprotectedPods);
            this.netfetchScore = response.data.Score;
            this.updateExpandedNamespaces();
            this.namespaceVisualizationData = {};
          } catch (error) {
            console.error('Error fetching scan results:', error);
            }
          }
        },
        async fetchNamespacesWithPolicies() {
          try {
            const response = await axios.get('/namespaces-with-policies');
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
          const response = await axios.post('/add-policy', { namespace });
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
      async suggestPolicy() {
        if (!this.selectedNamespace) {
          alert('Please select a namespace.');
          return;
        }

        this.unprotectedPods = [];
        this.scanInitiated = false;
        this.message = null;

        try {
          const podInfo = await this.fetchPodInfo(this.selectedNamespace);
          if (!podInfo || !Array.isArray(podInfo)) {
            console.error('No pod information available for the namespace:', this.selectedNamespace);
            return; // Return early if no pod information is available
          }

          const existingPolicies = await this.fetchExistingPolicies(this.selectedNamespace);

          // Generate network policy suggestions
          let suggestedPolicies = this.suggestNetworkPolicy(podInfo, existingPolicies);

          // Filter out policies that already exist
          suggestedPolicies = suggestedPolicies.filter(policy => {
            const policyName = yaml.load(policy).metadata.name;
            return !existingPolicies.includes(policyName);
          });

          this.suggestedNetworkPolicies = suggestedPolicies;

          console.log(this.suggestedNetworkPolicies);
        } catch (error) {
          console.error('Error suggesting policies:', error);
        }
        if (this.selectedNamespace && this.suggestedNetworkPolicies.length > 0) {
          this.activeNamespaceForPolicies = this.selectedNamespace;
        }
      },
      suggestNetworkPolicy(podData, existingPolicies) {
        let policies = [];
        const existingPolicyIdentifiers = new Set(existingPolicies);

        if (!Array.isArray(podData)) {
          console.error('Invalid podData: Expected an array, received:', podData);
          return [];
        }

        podData.forEach(pod => {
          let appName = pod.Labels['app.kubernetes.io/name'] || 'unknown-app';
          let matchLabels = {};
          if (pod.Labels['app.kubernetes.io/name']) {
            matchLabels['app.kubernetes.io/name'] = pod.Labels['app.kubernetes.io/name'];
          }

          if (Array.isArray(pod.Ports) && pod.Ports.length > 0) {
            pod.Ports.forEach(port => {
              const policyName = `allow-${appName}-${port.containerPort}-${port.protocol.toLowerCase()}-nfpol`.replace(/[^a-zA-Z0-9-]/g, '-');

              // Generate an identifier for each potential policy
              const policyIdentifier = JSON.stringify({
                namespace: pod.Namespace,
                podSelector: matchLabels,
                port: port.containerPort,
                protocol: port.protocol
              });

              if (!existingPolicyIdentifiers.has(policyIdentifier)) {
                policies.push({
                  apiVersion: 'networking.k8s.io/v1',
                  kind: 'NetworkPolicy',
                  metadata: {
                    name: policyName,
                    namespace: pod.Namespace
                  },
                  spec: {
                    podSelector: {
                      matchLabels: matchLabels
                    },
                    policyTypes: ['Ingress', 'Egress'],
                    ingress: [{
                      from: [],
                      ports: [{ protocol: port.protocol, port: port.containerPort }]
                    }],
                    egress: [{
                      to: [{ ipBlock: { cidr: '0.0.0.0/0' } }]
                    }]
                  }
                });
                // Add the identifier to the set to prevent duplicate creation
                existingPolicyIdentifiers.add(policyIdentifier);
              }
            });
          }
        });

        return policies.map(policy => yaml.dump(policy));
    },
      async fetchExistingPolicies(namespace) {
        try {
          const response = await axios.get(`/namespace-policies?namespace=${namespace}`);
          return response.data.items.map(policy => policy.metadata.name);
        } catch (error) {
          console.error('Error fetching existing policies:', error);
          return [];
        }
      },
      async fetchPodInfo(namespace) {
        try {
          const response = await axios.get(`/pod-info?namespace=${namespace}`);
          // Check if the data is present and is an array, else return an empty array
          return Array.isArray(response.data) ? response.data : [];
        } catch (error) {
          console.error('Error fetching pod information:', error);
          // Return an empty array to indicate no pods were found or there was an error
          return [];
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
          const response = await axios.get('/namespaces');
          this.allNamespaces = response.data.namespaces;
          console.log('Namespaces fetched:', this.allNamespaces);
        } catch (error) {
          console.error('Error fetching namespaces:', error);
        }
      },
      async fetchScanResultsForNamespace() {
        if (this.isScanForNative) {
          this.suggestedNetworkPolicies = [];
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
            const scanResponse = await axios.get(`/scan?namespace=${namespace}`);
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
            const response = await axios.get(`/visualization?namespace=${namespace}`);
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
          const response = await axios.get('/visualization/cluster');
        if (response.data && Array.isArray(response.data)) {
          this.clusterVisualizationData = response.data;
        } else {
          console.error('Received data is not an array:', response.data);
          this.clusterVisualizationData = [];
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
          const response = await axios.get(`/visualization?namespace=${namespace}`);
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
    color: black;
    text-decoration: none;
    border-left: 5px solid transparent;
  }

  .menu-item a.active {
    border-left: 5px solid #87CEEB;
    background-color: #eee;
  }

  .menu-item a:hover {
    background-color: #87CEEB;
  }

  .content {
    flex-grow: 1;
    padding: 40px;
  }

  .header {
    display: flex;
    min-height: 100px;
    justify-content: space-between;
    align-items: center;
  }

  .dashboard-title {
    color: black;
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
    margin-left: auto;
    margin-right: 30px;
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
    justify-content: space-between;
    align-items: center;
    gap: 10px;
    margin-bottom: 20px;
    margin-top: 10px;
  }

  .left-buttons {
    display: flex;
    gap: 10px;
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
    padding: 10px 15px;
    border-radius: 5px;
    border: 1px solid #87CEEB;
    background-color: white;
    font-size: 14px;
    color: #333;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    background-image: url('data:image/svg+xml;utf8,<svg fill="%2387CEEB" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M7 10l5 5 5-5z"/></svg>');
    background-repeat: no-repeat;
    background-position: right 10px center;
    cursor: pointer;
    transition: border-color 0.3s;
  }

  .namespace-select:hover {
    border-color: #5cb7db;
  }

  .namespace-select:focus {
    outline: none;
    border-color: #5cb7db;
  }

  .scan-btn:hover {
    background-color:#87CEEB;
    color: #fff;
  }

  .message-container {
    margin-bottom: 20px;
  }

  .notification {
    padding: 10px;
    margin: 10px 0;
    border-radius: 5px;
    text-align: center;
  }
  .success {
    background-color: #d4edda;
    color: #155724;
    border: 1px solid #c3e6cb;
  }

  .error {
    background-color: #f8d7da;
    color: #721c24;
    border: 1px solid #f5c6cb;
  }

  /* Policy container */
  .policy-container {
    margin-top: 20px;
    padding: 20px;
    background-color: #fff;
    border-radius: 5px;
  }

  .policy-header {
    margin-bottom: 20px;
  }

  .tooltip-icon {
  position: relative;
  display: inline-block;
  cursor: pointer;
  z-index: 1500;
}

.tooltiptext {
  visibility: hidden;
  width: 600px;
  background-color: white;
  color: rgb(55, 54, 54);
  text-align: left;
  border-radius: 4px;
  font-size: 13px;
  padding: 8px 12px;
  position: absolute;
  z-index: 1400;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%);
  border: 1px solid #87CEEB;
  transition: visibility 0.2s, opacity 0.2s ease-in-out;
  opacity: 0;
}

.tooltip-icon:hover .tooltiptext {
  visibility: visible;
  opacity: 1;
}

  .policy-cards-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-start;
    gap: 30px;
    align-items: stretch;
    min-height: 400px;
  }

  .policy-card {
    flex: 0 0 calc(33.333% - 30px);
    border: 1px solid #87CEEB;
    border-radius: 5px;
    padding: 10px;
    background-color: #fff;
    overflow: hidden; 
    min-height: 150px;
    position: relative;
    justify-content: space-between;
  }

  .policy-card textarea {
    width: 100%;
    height: 93%;
    padding: 10px;
    border-radius: 5px;
    border: 1px solid #87CEEB;
    font-family: monospace;
    resize: none;
  }

  .policy-yaml {
    white-space: pre-wrap;
    margin-bottom: 0;
    flex-grow: 1;
  }

  .policy-pagination {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-top: 20px;
  }

  .policy-card button {
    padding: 3px 5px;
    margin: 0;
    background: white;
    color: black;
    border: 1px solid #87CEEB;
    cursor: pointer;
    border-radius: 5px;
  }

  .policy-card-buttons {
    display: flex;
    justify-content: space-between;
    padding-top: 10px;
  }

  .policy-card button:hover {
    background-color: #5cb7db;
  }

  .copy-container {
    position: absolute;
    top: 10px;
    right: 10px;
  }

  .copy-icon {
    width: 12px;
    height: 12px;
    fill: currentColor;
  }

  .copy-success {
    fill: #d4edda;
  }

  .action-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: auto;
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

  .checkbox-container {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 25px;
  margin-top: 5px;
  margin-bottom: 0;
  cursor: pointer;
  font-size: 12px;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.policy-type-checkboxes {
  display: flex;
  justify-content: flex-start;
  gap: 2.5px;
}

.checkbox-container input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.checkmark {
  position: absolute;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  height: 18px;
  width: 18px;
  background-color: #eee;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.checkbox-container:hover input ~ .checkmark {
  background-color: #ccc;
}

.checkbox-container input:checked ~ .checkmark {
  background-color: #2196F3;
}

.checkmark:after {
  content: "";
  position: absolute;
  display: none;
}

.checkbox-container input:checked ~ .checkmark:after {
  display: block;
}

.checkbox-container .checkmark:after {
  left: 6px;
  top: 3px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  -webkit-transform: rotate(45deg);
  -ms-transform: rotate(45deg);
  transform: rotate(45deg);
}

.checkbox-container:not(:last-child) {
  margin-right: 10px;
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

  .app-container.dark-mode {
    background-color: #000;
    color: #fff;
  }

  .app-container.dark-mode .sidebar {
    background-color: #1a1a1a;
    color: #fff;
  }

  .app-container.dark-mode .menu-item a {
    color: #fff;
  }

  .app-container.dark-mode .menu-item a.active {
    background-color: #333;
  }

  .app-container.dark-mode .content {
    color: #fff;
  }

  .app-container.dark-mode .header {
    color: #fff;
  }

  .app-container.dark-mode .score-container {
    border: 5px solid #444;
  }

  .app-container.dark-mode .buttons .scan-btn {
    border: 2px solid #444;
    color: #fff;
  }

  .app-container.dark-mode .buttons .scan-btn:hover {
    background-color: #333;
  }

  .app-container.dark-mode .message-container .success-message,
  .app-container.dark-mode .message-container .error-message {
    color: #fff;
  }

  .app-container.dark-mode .table-container {
    background-color: #1a1a1a;
    color: #fff;
  }

  .app-container.dark-mode .pods-table th,
  .app-container.dark-mode .pods-table td {
    border: 1px solid #444;
  }

  .app-container.dark-mode .pods-table th {
    background-color: #333;
  }

  .app-container.dark-mode .remediate-btn {
    background-color: #333;
  }

  .app-container.dark-mode .pagination-controls .pagination-btn {
    border: 1px solid #444;
  }

  .app-container.dark-mode .network-policy-visualization {
    border: 2px solid #444;
    background-color: #1a1a1a;
  }

  .app-container.dark-mode .dark-mode-toggle {
    background-color: #333;
    color: #fff;
  }

  .app-container.dark-mode .dark-mode-toggle:hover {
    background-color: #444;
  }

  </style>

