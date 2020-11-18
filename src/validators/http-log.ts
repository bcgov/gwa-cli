import Joi from 'joi';

const httpLogSchema = Joi.object({
  http_endpoint: Joi.string().uri().required(),
  method: Joi.string().valid('POST', 'PUT', 'PATCH').optional(),
  timeout: Joi.number().optional(),
  keepalive: Joi.number().optional(),
});

export default httpLogSchema;
