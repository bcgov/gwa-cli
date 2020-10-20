export interface IIpRestriction {
  name: string;
  config: {
    allow: string[];
  };
}

export interface IGwaIpAnonymity {
  name: string;
  config: {
    ipv4Mask: number;
    ipv6Mask: number;
  };
}

export interface IGwaEndpoint {
  name: string;
  config: {
    apiOwners: string[];
  };
}

export interface IOidc {
  name: string;
  config: {
    responseType: string;
    introspectionEndpoint: string;
    filters: string | null;
    bearerOnly: string;
    sslVerify: string;
    sessionSecret: string | null;
    introspectionEndpointAuthMethod: string | null;
    realm: string;
    redirectAfterLogoutUri: string;
    scope: string;
    tokenEndpointAuthMethod: string;
    logoutPath: string;
    clientId: string;
    clientSecret: string;
    discovery: string;
    recoveryPagePath: string | null;
    redirectUriPath: string | null;
  };
}

export interface IRateLimiting {
  name: string;
  config: {
    faultTolerant: boolean;
    hideClientHeaders: boolean;
    limitBy: string;
    minute: number;
    policy: string;
    headerName: string | null;
    second: number | null;
    hour: number | null;
    day: number | null;
    month: number | null;
    year: number | null;
    redisDatabase: number;
    redisHost: string | null;
    redisPassword: string | null;
    redisPort: number;
    redisTimeout: number;
    protocols: ['http', 'https'];
  };
}

export interface ISpecExpose {
  name: string;
  config: {
    specUrl: string;
  };
}
