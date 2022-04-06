const { argv } = require('yargs/yargs')(process.argv.slice(2))
  .command('stage', 'The deployment stage')
  .demandOption(['stage']);

const fs = require('fs');

const manifest = './.serverless/manifest.json';
const { stage } = argv;

let data;
let json;

try {
  data = fs.readFileSync(manifest, { encoding: 'utf-8' });
  json = JSON.parse(data);
} catch (err) {
  console.error(err, err.error);
  process.exit(1);
}

const stageconf = json[stage];

if (!stageconf) {
  console.log(`Stage ${stage} not found in ${manifest}`);
  process.exit(1);
}

const wsconf = stageconf.outputs.filter((e) => e.OutputKey === 'ServiceEndpointWebsocket');

if (!wsconf) {
  console.log(`Outputkey ServiceEndpointWebsocket not found in ${JSON.stringify(wsconf)}`);
  process.exit(1);
}

const wsurl = wsconf[0].OutputValue;

console.log(wsurl);
