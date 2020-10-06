import * as Joi from 'joi';

const oidcSchema = Joi.object({
  response_type: Joi.string().empty('').optional(),
  introspection_endpoint: Joi.string().uri(),
  filters: Joi.string().optional(),
  bearer_only: Joi.string().valid('yes', 'no').optional(),
  ssl_verify: Joi.string().valid('yes', 'no').optional(),
  session_secret: Joi.string().optional(),
  introspection_endpoint_auth_method: Joi.string().optional(),
  realm: Joi.string().optional(),
  redirect_after_logout_uri: Joi.string().optional(),
  scope: Joi.string().optional(),
  token_endpoint_auth_method: Joi.string().optional(),
  logout_path: Joi.string().optional(),
  client_id: Joi.string().required(),
  client_secret: Joi.string().required(),
  discovery: Joi.string().uri().optional(),
  recovery_page_path: Joi.string().optional(),
  redirect_uri_path: Joi.string().optional(),
});

export default oidSchema;
