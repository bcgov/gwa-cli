import gwaEndpoint from './plugins/gwa-endpoint';
import gwaIpAnonymity from './plugins/gwa-ip-anonymity';
import ipRestriction from './plugins/ip-restriction';
import oidc from './plugins/oidc';
import rateLimiting from './plugins/rate-limiting';
import specExpose from './plugins/spec-expose';

export default {
  'bcgov-gwa-endpoint': gwaEndpoint,
  'gwa-ip-anonymity': gwaIpAnonymity,
  'ip-restriction': ipRestriction,
  oidc,
  'rate-limiting': rateLimiting,
  'spec-expose': specExpose,
};
