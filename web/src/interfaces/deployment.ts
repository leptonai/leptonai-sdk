export interface Deployment {
  id: string;
  name: string;
  photon_id: string;
  created_at: number;
  status: {
    endpoint: { internal_endpoint: string; external_endpoint: string };
    state: string;
  };
  resource_requirement: {
    memory: number;
    cpu: number;
    min_replicas: number;
    accelerator_type?: string;
    accelerator_num?: number;
  };
  envs?: Array<DeploymentEnv | DeploymentSecretEnv>;
  mounts?: Array<DeploymentMount>;
}

export interface DeploymentMount {
  mount_path: string;
  path: string;
}

export interface DeploymentEnv {
  name: string;
  value: string;
}

export interface DeploymentSecretEnv {
  name: string;
  value_from: { secret_name_ref: string };
}

export interface Replica {
  id: string;
}

export interface Metric {
  metric: { name: string; handler?: string };
  values: Array<[number, string]>;
}

/**
 * https://pkg.go.dev/k8s.io/api/events/v1#Event
 */
export enum DeploymentEventTypes {
  Normal = "Normal",
  Warning = "Warning",
}
export interface DeploymentEvent {
  type: string;
  reason: DeploymentEventTypes | string;
  count: number;
  last_observed_time: string;
}
