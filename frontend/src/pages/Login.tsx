import { Link, useNavigate } from "react-router-dom"
import React, { useState } from "react"
import { login } from "../api"

export default function Login() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [err, setErr] = useState("")
  const navigate = useNavigate()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setErr("")
    try {
  const { user, token } = await login(username, password)
  // Use sessionStorage so two windows can log in as different users without clobbering
  sessionStorage.setItem("nt_user", JSON.stringify(user))
  sessionStorage.setItem("nt_token", token)
      navigate("/inbox")
    } catch (e: any) {
      setErr(e.message)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 to-blue-300 dark:from-gray-900 dark:to-gray-800 px-4">
      <form className="bg-white dark:bg-gray-900 shadow-2xl rounded-xl p-8 w-full max-w-sm flex flex-col gap-4" onSubmit={handleLogin}>
        <h2 className="text-center text-2xl font-bold text-gray-800 dark:text-white">Login</h2>
        <div>
          <input autoFocus name="username" autoComplete="username" placeholder="Username" value={username} onChange={e=>setUsername(e.target.value)} className="w-full px-3 py-2 rounded-lg border mt-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-blue-300" />
        </div>
        <div>
          <input type="password" name="password" autoComplete="current-password" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} className="w-full px-3 py-2 rounded-lg border mt-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-blue-300" />
        </div>
        {err && <div className="text-red-500 text-sm text-center font-medium">{err}</div>}
        <button type="submit" className="block bg-blue-600 hover:bg-blue-700 w-full text-white font-bold py-2 mt-1 rounded-lg transition">Login</button>
        <div className="mt-3 text-center text-sm text-gray-500 dark:text-gray-400">
          New here? <Link to="/register" className="underline hover:text-blue-700 font-semibold">Register</Link>
        </div>
      </form>
    </div>
  )
}
