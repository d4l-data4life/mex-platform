const fs = require('fs');
const de = require('../src/translations/accessDenied-de.json').accessDenied;
const en = require('../src/translations/accessDenied-en.json').accessDenied;

const html = fs
  .readFileSync(__dirname + '/../src/static/access-denied.html', 'utf8')
  .replace('{"de": {}: "en": {}}', JSON.stringify({ de, en }));
fs.writeFileSync(__dirname + '/../www/access-denied.html', html, 'utf8');
