import Joi from 'joi';

const oidcSchema = Joi.object({
  response_type: Joi.string().empty('').optional(),
  introspection_endpoint: Joi.string().uri(),
  filters: Joi.string().empty('').optional(),
  bearer_only: Joi.string().valid('yes', 'no').optional(),
  ssl_verify: Joi.string().valid('yes', 'no').optional(),
  session_secret: Joi.string().empty('').optional(),
  introspection_endpoint_auth_method: Joi.string().empty('').optional(),
  realm: Joi.string().empty('').optional(),
  redirect_after_logout_uri: Joi.string().empty('').optional(),
  scope: Joi.string().empty('').optional(),
  token_endpoint_auth_method: Joi.string().empty('').optional(),
  logout_path: Joi.string().empty('').optional(),
  client_id: Joi.string().empty('').required(),
  client_secret: Joi.string().empty('').required(),
  discovery: Joi.string().uri().empty('').optional(),
  recovery_page_path: Joi.string().empty('').optional(),
  redirect_uri_path: Joi.string().empty('').optional(),
});

export default oidcSchema;
