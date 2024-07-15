<script lang="ts">
  import { ensureNonEmpty } from '../utils/assertions'
  import ErrorBanner from './error-banner.svelte'
  import { createAccount, logIn } from '../net/requests'

  export let onLogIn: (accessToken: string) => void

  let name = ''
  let errorMessage = ''

  const register = async () => {
    const response = await createAccount(name)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      await login(ensureNonEmpty(response.user_id))
    }
  }

  const login = async (user_id: string) => {
    const response = await logIn(user_id)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      onLogIn(ensureNonEmpty(response.access_token))
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
  <form on:submit|preventDefault={register}>
    <input bind:value={name} type="text" placeholder="Enter your name" required />
    <button type="submit">Create account</button>
  </form>
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
