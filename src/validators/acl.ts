import Joi from 'joi';

interface AclSchema {
  allow: string;
  deny: string;
  hide_groups_header: boolean;
}

const aclSchema = Joi.object<AclSchema>({
  allow: Joi.string().empty(''),
  deny: Joi.string().empty(''),
  hide_groups_header: Joi.boolean().optional(),
}).or('allow', 'deny');

export default aclSchema;
