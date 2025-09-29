export interface Channel {
  id: string;
  name: string;
  vendor: 'nvidia' | 'amd' | 'intel';
  channel: 'stable' | 'lts' | 'canary';
  phase: 'Pending' | 'Progressing' | 'Succeeded' | 'Failed';
  observedVersion: string;
  clusterCount: number;
  lastUpdated: string;
}

export interface Cluster {
  id: string;
  name: string;
  vendor: 'nvidia' | 'amd' | 'intel';
  gpuCount: number;
  gpuType: string;
  migEnabled: boolean;
  status: 'Ready' | 'NotReady' | 'Unknown';
}

export interface BundleDeployment {
  id: string;
  channelId: string;
  clusterId: string;
  status: 'Ready' | 'Failed' | 'Progressing';
  message: string;
  lastTransitionTime: string;
}

export interface DriftEvent {
  id: string;
  channelId: string;
  clusterId: string;
  type: 'VersionMismatch' | 'ConfigurationDrift' | 'PolicyViolation';
  severity: 'Low' | 'Medium' | 'High';
  message: string;
  timestamp: string;
}

export interface PolicyViolation {
  id: string;
  clusterId: string;
  policy: string;
  resource: string;
  message: string;
  timestamp: string;
}

export class MockProvider {
  private channels: Channel[] = [
    {
      id: 'nvidia-stable',
      name: 'NVIDIA Stable',
      vendor: 'nvidia',
      channel: 'stable',
      phase: 'Succeeded',
      observedVersion: 'v24.9.0/12.4.1',
      clusterCount: 3,
      lastUpdated: '2024-01-15T10:30:00Z'
    },
    {
      id: 'amd-lts',
      name: 'AMD LTS',
      vendor: 'amd',
      channel: 'lts',
      phase: 'Progressing',
      observedVersion: 'v0.9.0/5.6.0',
      clusterCount: 2,
      lastUpdated: '2024-01-15T09:15:00Z'
    },
    {
      id: 'intel-canary',
      name: 'Intel Canary',
      vendor: 'intel',
      channel: 'canary',
      phase: 'Failed',
      observedVersion: 'v0.5.0-rc1/23.4.0-rc1',
      clusterCount: 1,
      lastUpdated: '2024-01-15T08:45:00Z'
    }
  ];

  private clusters: Cluster[] = [
    {
      id: 'cluster-1',
      name: 'GPU Cluster 1',
      vendor: 'nvidia',
      gpuCount: 8,
      gpuType: 'NVIDIA A100 80GB',
      migEnabled: true,
      status: 'Ready'
    },
    {
      id: 'cluster-2',
      name: 'GPU Cluster 2',
      vendor: 'nvidia',
      gpuCount: 4,
      gpuType: 'NVIDIA RTX 4090',
      migEnabled: false,
      status: 'Ready'
    },
    {
      id: 'cluster-3',
      name: 'AMD Cluster 1',
      vendor: 'amd',
      gpuCount: 4,
      gpuType: 'AMD Radeon RX 7900 XTX',
      migEnabled: false,
      status: 'Ready'
    },
    {
      id: 'cluster-4',
      name: 'Intel Cluster 1',
      vendor: 'intel',
      gpuCount: 2,
      gpuType: 'Intel Arc A770',
      migEnabled: false,
      status: 'NotReady'
    }
  ];

  private bundleDeployments: BundleDeployment[] = [
    {
      id: 'bd-1',
      channelId: 'nvidia-stable',
      clusterId: 'cluster-1',
      status: 'Ready',
      message: 'Bundle deployed successfully',
      lastTransitionTime: '2024-01-15T10:30:00Z'
    },
    {
      id: 'bd-2',
      channelId: 'nvidia-stable',
      clusterId: 'cluster-2',
      status: 'Ready',
      message: 'Bundle deployed successfully',
      lastTransitionTime: '2024-01-15T10:25:00Z'
    },
    {
      id: 'bd-3',
      channelId: 'amd-lts',
      clusterId: 'cluster-3',
      status: 'Progressing',
      message: 'Bundle deployment in progress',
      lastTransitionTime: '2024-01-15T09:15:00Z'
    },
    {
      id: 'bd-4',
      channelId: 'intel-canary',
      clusterId: 'cluster-4',
      status: 'Failed',
      message: 'Bundle deployment failed: GPU driver not found',
      lastTransitionTime: '2024-01-15T08:45:00Z'
    }
  ];

  private driftEvents: DriftEvent[] = [
    {
      id: 'drift-1',
      channelId: 'nvidia-stable',
      clusterId: 'cluster-1',
      type: 'VersionMismatch',
      severity: 'Medium',
      message: 'Cluster running v24.8.0 but channel specifies v24.9.0',
      timestamp: '2024-01-15T11:00:00Z'
    },
    {
      id: 'drift-2',
      channelId: 'amd-lts',
      clusterId: 'cluster-3',
      type: 'ConfigurationDrift',
      severity: 'Low',
      message: 'GPU memory limit changed from 8GB to 16GB',
      timestamp: '2024-01-15T10:45:00Z'
    }
  ];

  private policyViolations: PolicyViolation[] = [
    {
      id: 'violation-1',
      clusterId: 'cluster-1',
      policy: 'limit-gpu-per-pod',
      resource: 'pod/vllm-inference',
      message: 'Pod requests 8 GPUs but policy limits to 4',
      timestamp: '2024-01-15T12:00:00Z'
    },
    {
      id: 'violation-2',
      clusterId: 'cluster-2',
      policy: 'require-runtime-class',
      resource: 'pod/pytorch-training',
      message: 'Pod missing required runtime class',
      timestamp: '2024-01-15T11:30:00Z'
    }
  ];

  async getChannels(): Promise<Channel[]> {
    return Promise.resolve(this.channels);
  }

  async getClusters(): Promise<Cluster[]> {
    return Promise.resolve(this.clusters);
  }

  async getBundleDeployments(channelId?: string): Promise<BundleDeployment[]> {
    if (channelId) {
      return Promise.resolve(this.bundleDeployments.filter(bd => bd.channelId === channelId));
    }
    return Promise.resolve(this.bundleDeployments);
  }

  async getDriftEvents(channelId?: string): Promise<DriftEvent[]> {
    if (channelId) {
      return Promise.resolve(this.driftEvents.filter(event => event.channelId === channelId));
    }
    return Promise.resolve(this.driftEvents);
  }

  async getPolicyViolations(clusterId?: string): Promise<PolicyViolation[]> {
    if (clusterId) {
      return Promise.resolve(this.policyViolations.filter(violation => violation.clusterId === clusterId));
    }
    return Promise.resolve(this.policyViolations);
  }

  async createChannel(channel: Omit<Channel, 'id' | 'lastUpdated'>): Promise<Channel> {
    const newChannel: Channel = {
      ...channel,
      id: `channel-${Date.now()}`,
      lastUpdated: new Date().toISOString()
    };
    this.channels.push(newChannel);
    return Promise.resolve(newChannel);
  }

  async updateChannel(id: string, updates: Partial<Channel>): Promise<Channel> {
    const index = this.channels.findIndex(c => c.id === id);
    if (index === -1) {
      throw new Error(`Channel ${id} not found`);
    }
    this.channels[index] = { ...this.channels[index], ...updates, lastUpdated: new Date().toISOString() };
    return Promise.resolve(this.channels[index]);
  }

  async deleteChannel(id: string): Promise<void> {
    const index = this.channels.findIndex(c => c.id === id);
    if (index === -1) {
      throw new Error(`Channel ${id} not found`);
    }
    this.channels.splice(index, 1);
    return Promise.resolve();
  }
}
