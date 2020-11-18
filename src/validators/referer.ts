import Joi from 'joi';

const refererSchema = Joi.object({
  referers: Joi.string().required(),
});

export default refererSchema;
