import * as Joi from 'joi';

const basicAuthSchema = Joi.object({
  hide_credentials: Joi.boolean().optional(),
  anonymous: Joi.string().optional(),
});

export default basicAuthSchema;
