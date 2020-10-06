import * as Joi from 'joi';

const jwtSchema = Joi.object({
  uri_param_names: Joi.array().items(Join.string().required()).optional(),
  cookie_names: Joi.array().items(Join.string().required()).optional(),
  header_names: Joi.array().items(Join.string().required()).optional(),
  claims_to_verify: Joi.array().items(Join.string().required()).optional(),
  secret_is_base64: Joi.boolean().optional(),
  anonymous: Joi.string().empty('').optional(),
  run_on_preflight: Joi.boolean().optional(),
  maximum_expiration: Joi.number().optional(),
});

export default jwtSchema;
