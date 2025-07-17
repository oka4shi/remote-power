export type SpinnerStatuses = "done" | "loading" | "none";

export type NetworkStatus = -1 | 0 | 1;

export type NetworkStatusResponse = {
  ifaces: {
    name: string;
    status: NetworkStatus;
  }[];
  updated_at: string;
  error: boolean;
};
