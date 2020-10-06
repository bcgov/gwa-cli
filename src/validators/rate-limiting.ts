import * as Joi from 'joi';

const rateLimitingSchema = Joi.object({
  fault_tolerant: Joi.boolean().optional(),
  hide_client_headers: Joi.boolean().optional(),
  limit_by: Joi.string()
    .valid('consumer', 'credential', 'ip', 'service', 'header')
    .optional(),
  minute: Joi.number().optional(),
  policy: Joi.string().empty().optional(),
  header_name: Joi.number().optional(),
  second: Joi.number().optional(),
  hour: Joi.number().optional(),
  day: Joi.number().optional(),
  month: Joi.number().optional(),
  year: Joi.number().optional(),
  redis_database: Joi.number().optional(),
  redis_host: Joi.string().empty('').optional(),
  redis_password: Joi.number().optional(),
  redis_port: Joi.number().optional(),
  redis_timeout: Joi.number().optional(),
});

export default rateLimitingSchema;
