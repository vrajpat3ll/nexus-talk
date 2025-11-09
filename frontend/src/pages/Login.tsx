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
    <div className="flex h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="flex-1 flex items-center justify-center px-4">
        <form className="bg-slate-800/50 backdrop-blur-xl border border-slate-700/50 shadow-2xl rounded-2xl p-8 w-full max-w-md flex flex-col gap-5" onSubmit={handleLogin}>
          <div className="text-center mb-2">
            <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-3xl shadow-lg">
              💬
            </div>
            <h2 className="text-3xl font-bold text-white mb-1">Welcome Back</h2>
            <p className="text-gray-400 text-sm">Sign in to continue to NexusTalk</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-300 mb-1 block">Username</label>
            <input 
              autoFocus 
              name="username" 
              autoComplete="username" 
              placeholder="Enter your username" 
              value={username} 
              onChange={e=>setUsername(e.target.value)} 
              className="w-full px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition" 
            />
          </div>
          <div>
            <label className="text-sm font-medium text-gray-300 mb-1 block">Password</label>
            <input 
              type="password" 
              name="password" 
              autoComplete="current-password" 
              placeholder="Enter your password" 
              value={password} 
              onChange={e=>setPassword(e.target.value)} 
              className="w-full px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition" 
            />
          </div>
          {err && <div className="text-red-400 text-sm text-center font-medium bg-red-500/10 border border-red-500/30 rounded-lg py-2">{err}</div>}
          <button 
            type="submit" 
            className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 w-full text-white font-bold py-3 rounded-xl transition shadow-lg hover:shadow-blue-500/50"
          >
            Sign In
          </button>
          <div className="mt-2 text-center text-sm text-gray-400">
            Don't have an account? <Link to="/register" className="text-blue-400 hover:text-blue-300 font-semibold">Create one</Link>
          </div>
        </form>
      </div>
    </div>
  )
}
