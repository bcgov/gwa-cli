import * as Joi from 'joi';

const requestTransformerSchema = Joi.object({
  http_method: Joi.string().empty('').optional(),
  remove: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    querystring: Joi.array().items(Joi.string().required()).optional(),
    body: Joi.array().items(Joi.string().required()).optional(),
  }),
  replace: Joi.object({
    uri: Joi.string().uri().optional(),
    headers: Joi.array().items(Joi.string().required()).optional(),
    querystring: Joi.array().items(Joi.string().required()).optional(),
    body: Joi.array().items(Joi.string().required()).optional(),
  }),
  add: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    querystring: Joi.array().items(Joi.string().required()).optional(),
    body: Joi.array().items(Joi.string().required()).optional(),
  }),
  append: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    querystring: Joi.array().items(Joi.string().required()).optional(),
    body: Joi.array().items(Joi.string().required()).optional(),
  }),
});

export default requestTransformerSchema;
