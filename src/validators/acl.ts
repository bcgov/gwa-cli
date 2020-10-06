import * as Joi from 'joi';

const aclSchema = Joi.object({
  allow: Joi.string().empty(''),
  deny: Joi.string().empty(''),
  hide_groups_header: Joi.boolean().optional(),
}).or('allow', 'deny');

console.log(aclSchema.type);
export default aclSchema;
