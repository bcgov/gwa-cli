import * as Joi from 'joi';

const kongSpecExposeSchema = Joi.object({
  spec_url: Joi.string().uri().required(),
});

export default kongSpecExposeSchema;
