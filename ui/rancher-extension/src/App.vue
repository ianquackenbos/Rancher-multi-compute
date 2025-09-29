<template>
  <div id="app">
    <header class="app-header">
      <h1>Rancher Multi-Compute</h1>
      <div class="provider-info">
        <span class="provider-badge" :class="{ mock: useMock, live: !useMock }">
          {{ useMock ? 'Mock Data' : 'Live Data' }}
        </span>
      </div>
    </header>

    <main class="app-main">
      <div class="dashboard-grid">
        <!-- Channels Overview -->
        <section class="dashboard-section">
          <h2>Channels</h2>
          <div class="channels-grid">
            <div 
              v-for="channel in channels" 
              :key="channel.id"
              class="channel-card"
              :class="channel.phase.toLowerCase()"
            >
              <div class="channel-header">
                <h3>{{ channel.name }}</h3>
                <span class="vendor-badge" :class="channel.vendor">
                  {{ channel.vendor.toUpperCase() }}
                </span>
              </div>
              <div class="channel-details">
                <p><strong>Channel:</strong> {{ channel.channel }}</p>
                <p><strong>Version:</strong> {{ channel.observedVersion }}</p>
                <p><strong>Clusters:</strong> {{ channel.clusterCount }}</p>
                <p><strong>Status:</strong> {{ channel.phase }}</p>
              </div>
            </div>
          </div>
        </section>

        <!-- Clusters Overview -->
        <section class="dashboard-section">
          <h2>Clusters</h2>
          <div class="clusters-grid">
            <div 
              v-for="cluster in clusters" 
              :key="cluster.id"
              class="cluster-card"
              :class="cluster.status.toLowerCase()"
            >
              <div class="cluster-header">
                <h3>{{ cluster.name }}</h3>
                <span class="status-badge" :class="cluster.status.toLowerCase()">
                  {{ cluster.status }}
                </span>
              </div>
              <div class="cluster-details">
                <p><strong>Vendor:</strong> {{ cluster.vendor.toUpperCase() }}</p>
                <p><strong>GPU Type:</strong> {{ cluster.gpuType }}</p>
                <p><strong>GPU Count:</strong> {{ cluster.gpuCount }}</p>
                <p v-if="cluster.migEnabled"><strong>MIG:</strong> Enabled</p>
              </div>
            </div>
          </div>
        </section>

        <!-- Drift Events -->
        <section class="dashboard-section">
          <h2>Drift Events</h2>
          <div class="drift-events">
            <div 
              v-for="event in driftEvents" 
              :key="event.id"
              class="drift-event"
              :class="event.severity.toLowerCase()"
            >
              <div class="event-header">
                <span class="event-type">{{ event.type }}</span>
                <span class="severity-badge" :class="event.severity.toLowerCase()">
                  {{ event.severity }}
                </span>
              </div>
              <p class="event-message">{{ event.message }}</p>
              <p class="event-time">{{ formatTime(event.timestamp) }}</p>
            </div>
          </div>
        </section>

        <!-- Policy Violations -->
        <section class="dashboard-section">
          <h2>Policy Violations</h2>
          <div class="policy-violations">
            <div 
              v-for="violation in policyViolations" 
              :key="violation.id"
              class="policy-violation"
            >
              <div class="violation-header">
                <span class="policy-name">{{ violation.policy }}</span>
                <span class="resource-name">{{ violation.resource }}</span>
              </div>
              <p class="violation-message">{{ violation.message }}</p>
              <p class="violation-time">{{ formatTime(violation.timestamp) }}</p>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, inject } from 'vue';
import type { DataProvider, Channel, Cluster, DriftEvent, PolicyViolation } from './providers';

const dataProvider = inject<DataProvider>('dataProvider')!;
const useMock = import.meta.env.VITE_USE_MOCK_PROVIDER !== 'false';

const channels = ref<Channel[]>([]);
const clusters = ref<Cluster[]>([]);
const driftEvents = ref<DriftEvent[]>([]);
const policyViolations = ref<PolicyViolation[]>([]);

onMounted(async () => {
  try {
    [channels.value, clusters.value, driftEvents.value, policyViolations.value] = await Promise.all([
      dataProvider.getChannels(),
      dataProvider.getClusters(),
      dataProvider.getDriftEvents(),
      dataProvider.getPolicyViolations(),
    ]);
  } catch (error) {
    console.error('Failed to load data:', error);
  }
});

function formatTime(timestamp: string): string {
  return new Date(timestamp).toLocaleString();
}
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.app-header h1 {
  margin: 0;
  color: #495057;
}

.provider-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.provider-badge.mock {
  background: #fff3cd;
  color: #856404;
}

.provider-badge.live {
  background: #d1ecf1;
  color: #0c5460;
}

.app-main {
  padding: 2rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
}

.dashboard-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.dashboard-section h2 {
  margin: 0 0 1rem 0;
  color: #495057;
  font-size: 1.25rem;
}

.channels-grid, .clusters-grid {
  display: grid;
  gap: 1rem;
}

.channel-card, .cluster-card {
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 1rem;
  transition: box-shadow 0.2s;
}

.channel-card:hover, .cluster-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.channel-card.succeeded {
  border-left: 4px solid #28a745;
}

.channel-card.progressing {
  border-left: 4px solid #ffc107;
}

.channel-card.failed {
  border-left: 4px solid #dc3545;
}

.cluster-card.ready {
  border-left: 4px solid #28a745;
}

.cluster-card.notready {
  border-left: 4px solid #dc3545;
}

.channel-header, .cluster-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.channel-header h3, .cluster-header h3 {
  margin: 0;
  font-size: 1rem;
  color: #495057;
}

.vendor-badge, .status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
}

.vendor-badge.nvidia {
  background: #76b900;
  color: white;
}

.vendor-badge.amd {
  background: #ed1c24;
  color: white;
}

.vendor-badge.intel {
  background: #0071c5;
  color: white;
}

.status-badge.ready {
  background: #28a745;
  color: white;
}

.status-badge.notready {
  background: #dc3545;
  color: white;
}

.channel-details, .cluster-details {
  font-size: 0.875rem;
  color: #6c757d;
}

.channel-details p, .cluster-details p {
  margin: 0.25rem 0;
}

.drift-events, .policy-violations {
  display: grid;
  gap: 1rem;
}

.drift-event, .policy-violation {
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 1rem;
}

.drift-event.high {
  border-left: 4px solid #dc3545;
}

.drift-event.medium {
  border-left: 4px solid #ffc107;
}

.drift-event.low {
  border-left: 4px solid #28a745;
}

.event-header, .violation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.event-type, .policy-name {
  font-weight: 500;
  color: #495057;
}

.severity-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
}

.severity-badge.high {
  background: #f8d7da;
  color: #721c24;
}

.severity-badge.medium {
  background: #fff3cd;
  color: #856404;
}

.severity-badge.low {
  background: #d1ecf1;
  color: #0c5460;
}

.resource-name {
  font-size: 0.875rem;
  color: #6c757d;
  font-family: monospace;
}

.event-message, .violation-message {
  margin: 0.5rem 0;
  color: #495057;
}

.event-time, .violation-time {
  margin: 0;
  font-size: 0.75rem;
  color: #6c757d;
}
</style>
