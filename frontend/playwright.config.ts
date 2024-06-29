import { defineConfig } from '@playwright/test';

export default defineConfig({
  reporter: [
    ['list'],
    ['json', {  outputFile: 'test-results/test-results.json' }],
    ['html', {  outputFolder: 'test-results', open: 'always' }]
  ],
});