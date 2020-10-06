import * as Joi from 'joi';

const requestTerminationSchema = Joi.object({
  status_code: Joi.number().optional(),
  message: Joi.string().empty('').optional(),
  body: Joi.string().empty('').optional(),
  content_type: Joi.string().empty('').optional(),
});

export default requestTerminationSchema;
