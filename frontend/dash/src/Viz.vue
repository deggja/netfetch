<template>
    <div class="network-policy-visualization" ref="vizContainer"></div>
  </template>  
  
  <script>
  import * as d3 from 'd3';

  // Define the custom cluster force
  function forceCluster(centers) {
    let nodes;
    let strength = 0.1;

    function force(alpha) {
      const l = alpha * strength;
      nodes.forEach(node => {
        const center = centers[node.cluster];
        if (center) {
          node.vx -= (node.x - center.x) * l;
          node.vy -= (node.y - center.y) * l;
        }
      });
    }

    force.initialize = function(_) {
      nodes = _;
    };

    force.strength = function(_) {
      if (typeof _ === 'number') {
        strength = _;
      } else {
        const clusterNodeCount = {};
        nodes.forEach(node => {
          strength = Math.min(1, 10 / (clusterNodeCount[node.cluster] || 1));
        });
      }
      return force;
    };

    return force;
  }


  export default {
  name: 'NetworkPolicyVisualization',
  props: {
    policies: Array,
    clusterData: {
      type: Array,
      default: () => [],
    },
    visualizationType: {
      type: String,
      default: 'namespace' // Possible values: 'namespace', 'cluster'
    },
  },
  mounted() {
    this.createNetworkMap();
  },
  methods: {
    createNetworkMap() {
    let data;
    if (this.visualizationType === 'cluster') {
      data = this.clusterData.reduce((acc, item) => {
        return item.policies && Array.isArray(item.policies) ? [...acc, ...item.policies] : acc;
      }, []);
    } else {
      data = this.policies;
    }

    if (!Array.isArray(data)) {
      console.error('Visualization data is undefined or not an array');
      return;
    }

    const color = d3.scaleOrdinal(d3.schemeCategory10);
    const width = this.$refs.vizContainer.clientWidth;
    const height = this.$refs.vizContainer.clientHeight;
    const nodes = [];
    const links = [];
    const namespaceClusterMap = new Map();

    // Group policies by namespace and prepare nodes and links
    let clusterIndex = 0;
    data.forEach(policy => {
      if (!namespaceClusterMap.has(policy.namespace)) {
        namespaceClusterMap.set(policy.namespace, clusterIndex++);
      }
      const namespaceCluster = namespaceClusterMap.get(policy.namespace);

      const policyNode = { id: policy.name, type: 'policy', cluster: namespaceCluster };
      nodes.push(policyNode);

      policy.targetPods.forEach(podName => {
        const podNode = { id: podName, type: 'pod', cluster: namespaceCluster };
        if (!nodes.some(n => n.id === podName)) {
          nodes.push(podNode);
        }
        links.push({ source: policy.name, target: podName });
      });
    });

    const clusterNodeCounts = nodes.reduce((counts, node) => {
      counts[node.cluster] = (counts[node.cluster] || 0) + 1;
      return counts;
    }, {});

    // Tooltip for full text display
    const tooltip = d3.select('body').append('div')
          .attr('class', 'tooltip')
          .style('visibility', 'hidden')
          .style('position', 'absolute')
          .style('background', 'white')
          .style('border', '1px solid black')
          .style('padding', '5px')
          .style('pointer-events', 'none');

    // Calculate cluster centers
    const gridColumns = Math.ceil(Math.sqrt(namespaceClusterMap.size));
    const gridRows = Math.ceil(namespaceClusterMap.size / gridColumns);
    const sectionWidth = width / gridColumns;
    const sectionHeight = height / gridRows;
    const boundaryPadding = Math.min(sectionWidth, sectionHeight) / 4;

    const clusterCenters = {};
    const namespaces = Array.from(namespaceClusterMap.keys());
    namespaces.forEach((namespace, index) => {
      const columnIndex = index % gridColumns;
      const rowIndex = Math.floor(index / gridColumns);
      const x = columnIndex * sectionWidth + sectionWidth / 2;
      const y = rowIndex * sectionHeight + sectionHeight / 2;
      clusterCenters[namespaceClusterMap.get(namespace)] = { x, y };
    });

    // Assign initial positions
    nodes.forEach(node => {
      const center = clusterCenters[node.cluster];
      node.x = center ? center.x : width / 2;
      node.y = center ? center.y : height / 2;
    });

    console.log('Nodes after initial position:', nodes);
    console.log('Cluster Centers:', clusterCenters);

    d3.select(this.$refs.vizContainer).selectAll('svg').remove();

    // Create SVG element
    const svg = d3.select(this.$refs.vizContainer)
      .append('svg')
      .attr('width', '100%')
      .attr('height', '100%')
      .attr('viewBox', `0 0 ${width} ${height}`);

    const containerGroup = svg.append('g');

    function drawClusterBoundaries() {
      containerGroup.selectAll('.cluster-boundary')
        .data(namespaces)
        .enter().append('rect')
        .attr('class', 'cluster-boundary')
        .attr('x', d => clusterCenters[namespaceClusterMap.get(d)].x - boundaryPadding / 2)
        .attr('y', d => clusterCenters[namespaceClusterMap.get(d)].y - boundaryPadding / 2)
        .attr('width', boundaryPadding)
        .attr('height', boundaryPadding)
        .attr('fill', 'none')
        .attr('stroke', '#ccc')
        .attr('stroke-dasharray', '4');
      }

    drawClusterBoundaries();

    createNamespaceLabels(containerGroup, namespaces, clusterCenters, namespaceClusterMap);

    // Console log for debugging
    console.log("Nodes before simulation:", nodes);
    console.log("Links:", links);
    console.log("Cluster Centers:", clusterCenters);

    const nodeRadius = 10;
    const collisionRadius = nodeRadius * 1.2;

    const linkDistance = (d) => {
    const clusterSize = clusterNodeCounts[d.source.cluster] || 1;
    const baseDistance = 50;
    const additionalDistancePerNode = 10;
    
    return baseDistance + (clusterSize * additionalDistancePerNode);
  };

    // Create the simulation with appropriate forces
    const simulation = d3.forceSimulation(nodes)
      .force('link', d3.forceLink(links).id(d => d.id).distance(linkDistance))
      .force('charge', d3.forceManyBody().strength(-50))
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collide', d3.forceCollide(collisionRadius).strength(1))
      .force('cluster', forceCluster(clusterCenters).strength(0.8))
      .force('x', d3.forceX().x(d => clusterCenters[d.cluster].x).strength(0.5))
      .force('y', d3.forceY().y(d => clusterCenters[d.cluster].y).strength(0.5));

    const node = containerGroup.append('g')
      .selectAll('circle')
      .data(nodes)
      .join('circle')
      .attr('r', 5)
      .attr('fill', d => color(namespaceClusterMap.get(d.namespace)))
      .call(drag(simulation));


    // Legends
    const legendGroup = svg.append('g')
      .attr('class', 'legend')
      .attr('transform', 'translate(10, 20)');

    legendGroup.append('circle')
      .attr('r', 5)
      .attr('fill', '#28a745')
      .attr('cx', 0)
      .attr('cy', 0);

    legendGroup.append('text')
      .attr('x', 10)
      .attr('y', 5)
      .text('Pod');

    legendGroup.append('circle')
      .attr('r', 5)
      .attr('fill', '#007bff')
      .attr('cx', 0)
      .attr('cy', 20);

    legendGroup.append('text')
      .attr('x', 10)
      .attr('y', 25)
      .text('Policy');

    // Setup zoom behavior
    const zoom = d3.zoom()
        .scaleExtent([0.1, 10])
        .translateExtent([[-width, -height], [2 * width, 2 * height]]) // This should cover the zoomed area
        .on('zoom', (event) => {
            containerGroup.attr('transform', event.transform);
        });

        svg.call(zoom);

    // Create links and nodes inside the containerGroup
    const link = containerGroup.append('g')
          .attr('stroke', '#999')
          .attr('stroke-opacity', 0.6)
          .selectAll('line')
          .data(links)
          .join('line')
          .attr('stroke-width', d => Math.sqrt(d.value));

    // Add labels with hover functionality to show full text
    const labels = containerGroup.append('g')
        .attr('class', 'labels')
        .selectAll('text')
        .data(nodes)
        .enter().append('text')
        .attr('dx', 12)
        .attr('dy', '.35em')
        .text(d => {
            // Display full name for policy and truncate for pod
            return d.type === 'policy' ? d.id : (d.id.substring(0, 10) + (d.id.length > 10 ? '...' : ''));
        })
        .on('mouseover', (event, d) => {
            tooltip.html(d.id)
                .style('visibility', 'visible')
                .style('left', (event.pageX + 10) + 'px')
                .style('top', (event.pageY - 10) + 'px');
        })
        .on('mouseout', () => {
            tooltip.style('visibility', 'hidden');
        });

        simulation.on('tick', () => {
          const padding = boundaryPadding / 2;
          // Keep nodes within the boundaries
          nodes.forEach(node => {
            const center = clusterCenters[node.cluster];
            if (!center) {
              console.error('No center found for node', node);
              return; // Skip this node to avoid errors
            }
            // Ensure nodes are within the boundary minus the padding
            node.x = Math.max(center.x - padding, Math.min(center.x + padding, node.x));
            node.y = Math.max(center.y - padding, Math.min(center.y + padding, node.y));
          });

          // Update link positions, checking for NaN values
          link
            .attr('x1', d => isNaN(d.source.x) ? 0 : d.source.x)
            .attr('y1', d => isNaN(d.source.y) ? 0 : d.source.y)
            .attr('x2', d => isNaN(d.target.x) ? 0 : d.target.x)
            .attr('y2', d => isNaN(d.target.y) ? 0 : d.target.y);

          // Update node positions, checking for NaN values
          node
            .attr('cx', d => isNaN(d.x) ? 0 : d.x)
            .attr('cy', d => isNaN(d.y) ? 0 : d.y);

          // Update label positions, checking for NaN values
          labels
            .attr('x', d => isNaN(d.x) ? 0 : d.x)
            .attr('y', d => isNaN(d.y) ? 0 : d.y);
        });


    // Drag functionality
    function drag(simulation) {
      function dragstarted(event) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
      }

      function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
      }

      function dragended(event) {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
      }

      return d3.drag()
          .on('start', dragstarted)
          .on('drag', dragged)
          .on('end', dragended);
          }

      function createNamespaceLabels() {
        containerGroup.selectAll('.namespace-label')
          .data(namespaces)
          .enter().append('text')
          .attr('class', 'namespace-label')
          .attr('x', d => clusterCenters[namespaceClusterMap.get(d)].x)
          .attr('y', d => clusterCenters[namespaceClusterMap.get(d)].y - boundaryPadding / 2)
          .attr('text-anchor', 'middle')
          .attr('fill', '#555')
          .text(d => d)
          .attr('font-size', '12px')
          .attr('font-weight', 'bold');
        }

        createNamespaceLabels();
      },
    },
    watch: {
        policies(newVal) {
            console.log('Policies data received:', newVal);
            this.createNetworkMap();
        }
    }
};
</script>
  
<style scoped>
.network-policy-visualization svg {
  border: 1px solid #87CEEB;
  border-radius: 5px;
  background-color: #f9f9f9;
}

.network-policy-visualization circle {
  stroke: #87CEEB;
  stroke-width: 1.5px;
}

.network-policy-visualization text {
  font-size: 10px;
  pointer-events: none;
}

.network-policy-visualization {
  width: 100%;
  padding: 20px;
  height: 600px;
}

.cluster-boundary {
  stroke-opacity: 0.6;
  stroke-width: 1px;
}

.namespace-label {
  font-size: 12px;
  font-weight: bold;
}

.tooltip {
    position: absolute;
    visibility: hidden;
    background: white;
    border: 1px solid black;
    padding: 5px;
    pointer-events: none;
    z-index: 10;
  }
</style>
  