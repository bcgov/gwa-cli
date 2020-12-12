import fetch from 'node-fetch';
import { URLSearchParams } from 'url';

import config from '../config';

async function authenticate(url: string): Promise<string> {
  try {
    const { clientId, clientSecret } = config();
    const body = new URLSearchParams();
    body.append('client_id', clientId);
    body.append('client_secret', clientSecret);
    body.append('grant_type', 'client_credentials');

    const res = await fetch(url, {
      method: 'POST',
      body,
    });

    if (res.ok) {
      const json = await res.json();
      return json.access_token;
    } else {
      throw res.statusText;
    }
  } catch (err) {
    throw new Error(err);
  }
}

export default authenticate;
