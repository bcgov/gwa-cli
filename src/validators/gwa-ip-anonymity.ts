import * as Joi from 'joi';

const gwaIpAnonymitySchema = Joi.object({
  ipv4_mask: Joi.number().required(),
  ipv6_mask: Joi.number().required(),
});

export default gwaIpAnonymitySchema;
