export { MockProvider, LiveProvider } from './MockProvider';
export type { Channel, Cluster, BundleDeployment, DriftEvent, PolicyViolation } from './MockProvider';

export interface DataProvider {
  getChannels(): Promise<Channel[]>;
  getClusters(): Promise<Cluster[]>;
  getBundleDeployments(channelId?: string): Promise<BundleDeployment[]>;
  getDriftEvents(channelId?: string): Promise<DriftEvent[]>;
  getPolicyViolations(clusterId?: string): Promise<PolicyViolation[]>;
  createChannel(channel: Omit<Channel, 'id' | 'lastUpdated'>): Promise<Channel>;
  updateChannel(id: string, updates: Partial<Channel>): Promise<Channel>;
  deleteChannel(id: string): Promise<void>;
}

export function createDataProvider(useMock: boolean = true): DataProvider {
  if (useMock) {
    return new MockProvider();
  } else {
    return new LiveProvider();
  }
}
