<script lang="ts">
  import { ensureNonEmpty } from '../utils/assertions'

  let name = ''
  export let onLogIn: (accessToken: string) => void

  type RegisterResponse = {
    account_number: string
  }

  type LogInResponse = {
    access_token: string
  }

  const register = async () => {
    // TODO prevent SQL injection here
    let response = await fetch('http://localhost:8080/users/register', {
      method: 'POST',
      body: JSON.stringify({ name: name })
    })

    if (response.ok) {
      const registerResponse = (await response.json()) as RegisterResponse
      ensureNonEmpty(registerResponse.account_number)
      await login(registerResponse.account_number)
    } else {
      alert(`Failed to register - backend error (${response.status})`)
    }
  }

  const login = async (accountNumber: string) => {
    let response = await fetch('http://localhost:8080/users/login', {
      method: 'POST',
      body: JSON.stringify({ account_number: accountNumber })
    })

    if (response.ok) {
      const loginResponse = (await response.json()) as LogInResponse
      ensureNonEmpty(loginResponse.access_token)
      onLogIn(loginResponse.access_token)
    } else {
      alert(`Failed to login - backend error (${response.status})`)
    }
  }
</script>

<h1>Tasks</h1>

<input bind:value={name} type="text" placeholder="Enter your name" />
<button on:click={() => register()}>Get started</button>
