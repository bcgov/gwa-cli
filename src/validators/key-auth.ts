import Joi from 'joi';

const keyAuthSchema = Joi.object({
  key_names: Joi.array()
    .items(Joi.string().pattern(/^[a-zA-Z0-9_-]+$/))
    .optional(),
  key_in_body: Joi.boolean().optional(),
  hide_credentials: Joi.boolean().optional(),
  anonymous: Joi.string().empty('').optional(),
  run_on_preflight: Joi.boolean().optional(),
});

export default keyAuthSchema;
