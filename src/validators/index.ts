import acl from './acl';
import basicAuth from './basic-auth';
import bcgovGwaEndpoint from './bcgov-gwa-endpoint';
import cors from './cors';
import gwaIpAnonymity from './gwa-ip-anonymity';
import httpLog from './http-log';
import ipRestriction from './ip-restriction';
import jwt from './jwt';
import keyAuth from './key-auth';
import kongSpecExpose from './kong-spec-expose';
import oidc from './oidc';
import rateLimiting from './rate-limiting';
import referer from './referer';
import requestTermination from './request-termination';
import requestTransformer from './request-transformer';
import responseTransformer from './response-transformer';
import statsd from './statsd';

type Validators = {
  [key: string]: any;
};

const validators: Validators = {
  acl,
  'basic-auth': basicAuth,
  'bcgov-gwa-endpoint': bcgovGwaEndpoint,
  cors,
  'gwa-ip-anonymity': gwaIpAnonymity,
  'http-log': httpLog,
  'ip-restriction': ipRestriction,
  jwt,
  'key-auth': keyAuth,
  'kong-spec-expose': kongSpecExpose,
  oidc,
  'rate-limiting': rateLimiting,
  referer,
  'request-termination': requestTermination,
  'request-transformation': requestTransformer,
  'response-transformation': responseTransformer,
  statsd,
};

export default validators;
