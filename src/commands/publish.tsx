import fetch from 'node-fetch';
import FormData from 'form-data';
import fs from 'fs';
import path from 'path';

export default async function (input: string, options: any) {
  const cwd = process.cwd();

  console.log('Uploading config...');
  try {
    const form = new FormData();
    form.append('configFile', fs.createReadStream(path.resolve(cwd, input)));
    form.append('dryRun', Boolean(options.dryRun).toString());

    const req = await fetch('http://localhost:3000', {
      method: 'PUT',
      body: form,
    });

    if (req.ok) {
      console.log('File uploaded');
    }
  } catch (err) {
    console.error('Upload Failed');
    console.error(err);
  }
}
