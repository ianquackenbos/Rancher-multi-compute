import { Channel, Cluster, BundleDeployment, DriftEvent, PolicyViolation } from './MockProvider';

export class LiveProvider {
  private baseUrl: string;

  constructor(baseUrl: string = '/api/v1') {
    this.baseUrl = baseUrl;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async getChannels(): Promise<Channel[]> {
    return this.request<Channel[]>('/channels');
  }

  async getClusters(): Promise<Cluster[]> {
    return this.request<Cluster[]>('/clusters');
  }

  async getBundleDeployments(channelId?: string): Promise<BundleDeployment[]> {
    const endpoint = channelId ? `/channels/${channelId}/bundle-deployments` : '/bundle-deployments';
    return this.request<BundleDeployment[]>(endpoint);
  }

  async getDriftEvents(channelId?: string): Promise<DriftEvent[]> {
    const endpoint = channelId ? `/channels/${channelId}/drift-events` : '/drift-events';
    return this.request<DriftEvent[]>(endpoint);
  }

  async getPolicyViolations(clusterId?: string): Promise<PolicyViolation[]> {
    const endpoint = clusterId ? `/clusters/${clusterId}/policy-violations` : '/policy-violations';
    return this.request<PolicyViolation[]>(endpoint);
  }

  async createChannel(channel: Omit<Channel, 'id' | 'lastUpdated'>): Promise<Channel> {
    return this.request<Channel>('/channels', {
      method: 'POST',
      body: JSON.stringify(channel),
    });
  }

  async updateChannel(id: string, updates: Partial<Channel>): Promise<Channel> {
    return this.request<Channel>(`/channels/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
  }

  async deleteChannel(id: string): Promise<void> {
    await this.request<void>(`/channels/${id}`, {
      method: 'DELETE',
    });
  }
}
