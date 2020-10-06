import * as Joi from 'joi';

const corsSchema = Joi.object({
  origins: Joi.string().optional(),
  methods: Joi.array()
    .items(
      Joi.string().valid(
        'GET',
        'HEAD',
        'PUT',
        'PATCH',
        'POST',
        'DELETE',
        'OPTIONS',
        'TRACE',
        'CONNECT'
      )
    )
    .optional(),
  headers: Joi.array().items(Joi.string().required()).optional(),
  exposed_headers: Joi.array().items(Joi.string().required()).optional(),
  credentials: Joi.boolean().optional(),
  max_age: Joi.number().optional(),
  preflight_continue: Joi.boolean().optional(),
});

export default corsSchema;
