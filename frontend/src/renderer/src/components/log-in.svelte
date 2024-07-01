<script lang="ts">
  import { ensureNonEmpty } from '../utils/assertions'
  import ErrorBanner from './error-banner.svelte'
  import type { RegisterResponse, LogInResponse } from '../models/models'

  let name = ''
  let errorMessage = ''
  export let onLogIn: (accessToken: string) => void

  const register = async () => {
    if (name.length < 5) {
      errorMessage = 'Name must be at least 5 characters long'
      return
    }

    // TODO prevent SQL injection here
    let response = await fetch('http://localhost:8080/users/register', {
      method: 'POST',
      body: JSON.stringify({ name: name })
    })

    if (response.ok) {
      const registerResponse = (await response.json()) as RegisterResponse
      ensureNonEmpty(registerResponse.user_id)
      await login(registerResponse.user_id)
    } else {
      let error = await response.text()
      errorMessage = `Failed to register user (${error})`
    }
  }

  const login = async (user_id: string) => {
    let response = await fetch('http://localhost:8080/users/login', {
      method: 'POST',
      body: JSON.stringify({ user_id: user_id })
    })

    if (response.ok) {
      const loginResponse = (await response.json()) as LogInResponse
      ensureNonEmpty(loginResponse.access_token)
      onLogIn(loginResponse.access_token)
    } else {
      let error = await response.text()
      errorMessage = `Failed to log in (${error})`
    }
  }

  const onDismissError = () => {
    errorMessage = ''
    name = ''
  }
</script>

<ErrorBanner {errorMessage} {onDismissError} />

<div class="card">
  <h1>Tasks</h1>
  <input bind:value={name} type="text" placeholder="Enter your name" />
  <button on:click={() => register()}>Create account</button>
  <div class="footer">&copy; 2024 Jeroen Mols</div>
</div>

<style>
  @import '../assets/base.css';

  .card {
    background: var(--white);
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    padding: 40px;
    width: 350px;
    text-align: center;
  }
  .card h1 {
    font-weight: 600;
    margin-bottom: 20px;
    color: var(--primary);
  }
  .card input[type='text'] {
    width: 100%;
    padding: 10px;
    margin: 10px 0;
    border: 1px solid var(--light-gray);
    border-radius: 4px;
    font-size: 16px;
  }
  .card button {
    background-color: var(--primary);
    color: var(--white);
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
    width: 100%;
  }
  .card button:hover {
    background-color: var(--primary-hover);
  }
  .card .footer {
    margin-top: 20px;
    font-size: 14px;
    color: var(--dark-gray);
  }

  p {
    margin: 20px;
    color: var(--dark-gray);
  }
</style>
