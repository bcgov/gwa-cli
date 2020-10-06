import * as Joi from 'joi';

const bgGovGwaEndpointSchema = Joi.object({
  api_owners: Joi.array().items(Joi.string().required()),
});

export default bgGovGwaEndpointSchema;
