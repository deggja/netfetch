<template>
    <div class="network-policy-visualization" ref="vizContainer"></div>
  </template>  
  
  <script>
  import * as d3 from 'd3';

  // Define the custom cluster force
  function forceCluster(centers) {
  let nodes;
  let strength = 0.1; // Default strength value

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
    // This allows us to chain strength like other d3 forces if we want to set a strength
    if (typeof _ === 'number') {
      strength = _;
      return force;
    }
    return strength;
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

    const width = this.$refs.vizContainer.clientWidth;
    const height = this.$refs.vizContainer.clientHeight;
    const nodes = [];
    const links = [];
    const clusterMap = new Map();

    data.forEach((policy, index) => {
      const policyNode = { id: policy.name, type: 'policy', cluster: index };
      nodes.push(policyNode);
      clusterMap.set(policy.name, index);

      policy.targetPods.forEach(podName => {
        const podNode = { id: podName, type: 'pod', cluster: index };
        if (!nodes.some(n => n.id === podName)) {
          nodes.push(podNode);
        }
        links.push({ source: policy.name, target: podName });
      });
    });

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
    const uniqueClusters = Array.from(new Set(data.map((_, index) => index)));
    const numberOfClusters = uniqueClusters.length;
    const gridSize = Math.ceil(Math.sqrt(numberOfClusters));
    const sectionWidth = width / gridSize;
    const sectionHeight = height / gridSize;
    const clusterCenters = {};

    uniqueClusters.forEach((cluster, index) => {
      const x = (index % gridSize) * sectionWidth + sectionWidth / 2;
      const y = Math.floor(index / gridSize) * sectionHeight + sectionHeight / 2;
      clusterCenters[cluster] = { x, y };
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

    // Console log for debugging
    console.log("Nodes before simulation:", nodes);
    console.log("Links:", links);
    console.log("Cluster Centers:", clusterCenters);

    // Create the simulation with appropriate forces
    const simulation = d3.forceSimulation(nodes)
      .force('link', d3.forceLink(links).id(d => d.id).distance(50))
      .force('charge', d3.forceManyBody().strength(-50))
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('cluster', forceCluster(clusterCenters).strength(0.5))
      .force('x', d3.forceX().x(d => clusterCenters[d.cluster].x))
      .force('y', d3.forceY().y(d => clusterCenters[d.cluster].y));

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

    const node = containerGroup.append('g')
        .attr('stroke', '#fff')
        .attr('stroke-width', 1.5)
        .selectAll('circle')
        .data(nodes)
        .join('circle')
        .attr('r', 5)
        .attr('fill', d => d.type === 'policy' ? '#007bff' : '#28a745')
        .call(drag(simulation));

    // Add labels with hover functionality to show full text
    const labels = containerGroup.append('g')
          .attr('class', 'labels')
          .selectAll('text')
          .data(nodes)
          .enter().append('text')
            .attr('dx', 12)
            .attr('dy', '.35em')
            .text(d => d.id.substring(0, 10) + '...') // Truncate text
            .on('mouseover', (event, d) => {
              tooltip.html(d.id) // Set the full name for the tooltip
                .style('visibility', 'visible')
                .style('left', (event.pageX + 10) + 'px')
                .style('top', (event.pageY - 10) + 'px');
            })
            .on('mouseout', () => {
              tooltip.style('visibility', 'hidden');
            });

    // Update positions on each tick
    simulation.on('tick', () => {
      nodes.forEach(node => {
        if (isNaN(node.x) || isNaN(node.y)) {
          console.error('NaN position for node:', node);
        }
      });

      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

      node
        .attr('cx', d => d.x)
        .attr('cy', d => d.y);

      // Adjust label positions to prevent overlap
      labels.attr('x', d => d.x)
          .attr('y', d => d.y)
          .each(function (d, i) {
            for (let j = i + 1; j < nodes.length; ++j) {
              const other = labels.nodes()[j];
              const thisBox = this.getBBox();
              const otherBox = other.getBBox();
              
              const deltaX = d.x - nodes[j].x;
              const deltaY = d.y - nodes[j].y;
              const dist = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
              
              const padding = 2;
              if (dist < thisBox.width + otherBox.width + padding) {
                d.x += deltaX / dist * padding;
                d.y += deltaY / dist * padding;
                nodes[j].x -= deltaX / dist * padding;
                nodes[j].y -= deltaY / dist * padding;
              }
            }
        });
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
  