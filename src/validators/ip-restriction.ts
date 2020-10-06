import * as Joi from 'joi';

const ipRestrictionSchema = Joi.object({
  allow: Joi.array().items(Joi.string().ip()),
  deny: Joi.array().items(Joi.string().ip()),
}).or('allow', 'deny');

export default ipRestrictionSchema;
