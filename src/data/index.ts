import gwaEndpoint from './plugins/gwa-endpoint';
import gwaIpAnonymity from './plugins/gwa-ip-anonymity';
import ipRestriction from './plugins/ip-restriction';
import OIDC from './plugins/oidc';
import rateLimiting from './plugins/rate-limiting';
import specExpose from './plugins/spec-expose';

const plugins = {
  [gwaEndpoint.id]: gwaEndpoint,
  [gwaIpAnonymity.id]: gwaIpAnonymity,
  [ipRestriction.id]: ipRestriction,
  [OIDC.id]: OIDC,
  [rateLimiting.id]: rateLimiting,
  [specExpose.id]: specExpose,
};

export default {
  plugins,
};
