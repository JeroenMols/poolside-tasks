const { test, expect, _electron: electron } = require('@playwright/test')

test('end 2 end test', async () => {
  const electronApp = await electron.launch({ args: ['.'] })
  const isPackaged = await electronApp.evaluate(async ({ app }) => {
    // This runs in Electron's main process, parameter here is always
    // the result of the require('electron') in the main app script.
    return app.isPackaged
  })

  expect(isPackaged).toBe(false)

  // Login
  const window = await electronApp.firstWindow()
  await window.getByRole('textbox').fill('jeroen');
  await window.screenshot({ path: 'screenshots/intro.png' })
  
  await window.getByRole('button').click();

  // Todo lists
  await window.getByText('Choose a list').isVisible();
  await window.screenshot({ path: 'screenshots/lists.png' })
  await window.getByText('create new').click();

  // Todo
  await window.getByText('Nothing to do').isVisible();
  await window.getByRole('textbox').fill('My first todo');
  await window.getByText('Add').click();

  await window.getByText('My first todo').isVisible();
  await window.screenshot({ path: 'screenshots/todos-todo.png' })
  
  await window.getByText('next').click();
  await window.getByText('ongoing').isVisible();
  await window.screenshot({ path: 'screenshots/todos-ongoing.png' })

  await window.getByText('next').click();
  await window.getByText('done').isVisible();
  await window.screenshot({ path: 'screenshots/todos-done.png' })

  await window.getByText('next').click();
  await window.getByText('Failed to update todo ({"error":"invalid status transition from done to ongoing})').isVisible();
  await window.screenshot({ path: 'screenshots/todos-done.png' })

  await electronApp.close()
})