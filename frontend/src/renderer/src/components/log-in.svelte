<script lang="ts">
  let accountNumber = ''
  let accessToken = ''
  let name = ''
  export let onLogIn: (accessToken: string) => void

  console.log('App component loaded')

  type RegisterResponse = {
    account_number: string
  }

  type LogInResponse = {
    access_token: string
  }

  const register = async () => {
    let response = await fetch('http://localhost:8080/users/register', {
      method: 'POST',
      body: JSON.stringify({ name: name })
    })

    if (response.ok) {
      const body = (await response.json()) as RegisterResponse
      accountNumber = body.account_number
      await login(accountNumber)
    }
  }

  const login = async (accountNumber: string) => {
    let response = await fetch('http://localhost:8080/users/login', {
      method: 'POST',
      body: JSON.stringify({ account_number: accountNumber })
    })

    if (response.ok) {
      const body = (await response.json()) as LogInResponse
      accessToken = body.access_token
      onLogIn(accessToken)
    }
  }
</script>

<h1>Tasks</h1>

<input bind:value={name} type="text" placeholder="Enter your name" />
<button on:click={() => register()}>Get started</button>
