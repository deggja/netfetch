<template>
    <div class="network-policy-visualization" ref="vizContainer"></div>
  </template>  
  
  <script>
  import * as d3 from 'd3';

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
        // Flatten the policies array from the nested structure
        data = this.clusterVisualizationData.reduce((acc, item) => {
          if (item.policies && Array.isArray(item.policies)) {
            return [...acc, ...item.policies];
          }
          return acc;
        }, []);
      } else {
        data = this.policies;
      }

    console.log("Visualization data:", data);
    const container = d3.select(this.$refs.vizContainer);
    const width = 800;
    const height = 600;

    const nodes = [];
    const links = [];

    // Iterate over flattened policies
    data.forEach(policy => {
      if (!policy.targetPods || !Array.isArray(policy.targetPods)) {
        console.warn(`Skipping policy ${policy.name} as it has no targetPods or targetPods is not an array`);
        return;
      }

      const policyNode = { id: policy.name, type: 'policy' };
      nodes.push(policyNode);

      policy.targetPods.forEach(podName => {
        const podNode = { id: podName, type: 'pod' };
        nodes.push(podNode);
        links.push({ source: policy.name, target: podName });
      });
    });

    const svg = container.append('svg')
      .attr('width', width)
      .attr('height', height)
      .attr("viewBox", [0, 0, width, height]);

    // Create the simulation
    const simulation = d3.forceSimulation(nodes)
      .force('link', d3.forceLink(links).id(d => d.id))
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter(width / 2, height / 2));

    // Create links
    const link = svg.append('g')
      .attr('stroke', '#999')
      .attr('stroke-opacity', 0.6)
      .selectAll('line')
      .data(links)
      .join('line')
      .attr('stroke-width', d => Math.sqrt(d.value));

    // Create nodes
    const node = svg.append('g')
      .attr('stroke', '#fff')
      .attr('stroke-width', 1.5)
      .selectAll('circle')
      .data(nodes)
      .join('circle')
      .attr('r', 5)
      .attr('fill', d => d.type === 'policy' ? '#007bff' : '#28a745')
      .call(drag(simulation));

    // Node labels
    const labels = svg.append('g')
      .attr('class', 'labels')
      .selectAll('text')
      .data(nodes)
      .enter().append('text')
        .attr('dx', 12)
        .attr('dy', '.35em')
        .text(d => d.id);

    // Update positions on each tick
    simulation.on('tick', () => {
      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

      node
        .attr('cx', d => d.x)
        .attr('cy', d => d.y);

      labels
        .attr('x', d => d.x)
        .attr('y', d => d.y);
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
  width: 800px;
  height: 600px;
}
</style>
  