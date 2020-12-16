export interface StatusData {
  name: string;
  upstream: string;
  status: 'UP' | 'DOWN';
  reason: string;
  host: string;
  envHost: string;
}
