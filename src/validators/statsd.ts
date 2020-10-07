import Joi from 'joi';

const statsdSchema = Joi.object({
  host: Joi.string().ip().optional(),
  port: Joi.number().optional(),
  metrics: Joi.array().items(Joi.string().required()).optional(),
  prefix: Joi.string().empty('').optional(),
});

export default statsdSchema;
