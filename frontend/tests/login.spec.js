const { test, expect, _electron: electron } = require('@playwright/test')

test('example test', async () => {
  const electronApp = await electron.launch({ args: ['.'] })
  const isPackaged = await electronApp.evaluate(async ({ app }) => {
    // This runs in Electron's main process, parameter here is always
    // the result of the require('electron') in the main app script.
    return app.isPackaged
  })

  expect(isPackaged).toBe(false)

  // Login
  const window = await electronApp.firstWindow()
  await window.screenshot({ path: 'screenshots/intro.png' })
  await window.getByRole('textbox').fill('jeroen');
  await window.getByRole('button').click();

  // Todo lists
  await window.getByText('TODO lists').isVisible();
  await window.screenshot({ path: 'screenshots/todo.png' })

  await electronApp.close()
})