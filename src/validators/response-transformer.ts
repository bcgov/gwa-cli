import * as Joi from 'joi';

const responseTransformerSchema = Joi.object({
  http_method: Joi.string().empty('').optional(),
  remove: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    json: Joi.array().items(Joi.string().required()).optional(),
  }),
  replace: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    json: Joi.array().items(Joi.string().required()).optional(),
    json_types: Joi.array().items(Joi.string().required()).optional(),
  }),
  add: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    json: Joi.array().items(Joi.string().required()).optional(),
    json_types: Joi.array().items(Joi.string().required()).optional(),
  }),
  append: Joi.object({
    headers: Joi.array().items(Joi.string().required()).optional(),
    json: Joi.array().items(Joi.string().required()).optional(),
    json_types: Joi.array().items(Joi.string().required()).optional(),
  }),
});

export default responseTransformerSchema;
