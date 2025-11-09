import { Link, useNavigate } from "react-router-dom"
import React, { useState } from "react"

export default function Register() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [email, setEmail] = useState("")
  const [err, setErr] = useState("")
  const [ok, setOk] = useState(false)
  const navigate = useNavigate()

  const handleReg = async (e: React.FormEvent) => {
    e.preventDefault()
    setErr("")
    try {
      const API_HOST = import.meta.env.VITE_API_HOST || "localhost";
      const resp = await fetch(`http://${API_HOST}:8081/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password }),
      })
      if (!resp.ok) throw new Error("Registration failed")
      const js = await resp.json()
      // Store user temporarily so second tab can see them without refresh list endpoint yet
      const listRaw = localStorage.getItem("nt_reg_users")
      const arr = listRaw ? JSON.parse(listRaw) : []
      arr.push(js)
      localStorage.setItem("nt_reg_users", JSON.stringify(arr))
      setOk(true)
      setTimeout(() => navigate("/login"), 1200)
    } catch (e: any) {
      setErr(e.message)
    }
  }

  return (
    <div className="flex h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="flex-1 flex items-center justify-center px-4">
        <form className="bg-slate-800/50 backdrop-blur-xl border border-slate-700/50 shadow-2xl rounded-2xl p-8 w-full max-w-md flex flex-col gap-5" onSubmit={handleReg}>
          <div className="text-center mb-2">
            <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-3xl shadow-lg">
              💬
            </div>
            <h2 className="text-3xl font-bold text-white mb-1">Create Account</h2>
            <p className="text-gray-400 text-sm">Join NexusTalk today</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-300 mb-1 block">Username</label>
            <input 
              name="username" 
              placeholder="Choose a username" 
              value={username} 
              onChange={e=>setUsername(e.target.value)} 
              className="w-full px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition" 
            />
          </div>
          <div>
            <label className="text-sm font-medium text-gray-300 mb-1 block">Email</label>
            <input 
              type="email" 
              name="email" 
              placeholder="your@email.com" 
              value={email} 
              onChange={e=>setEmail(e.target.value)} 
              className="w-full px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition" 
            />
          </div>
          <div>
            <label className="text-sm font-medium text-gray-300 mb-1 block">Password</label>
            <input 
              type="password" 
              name="password" 
              placeholder="Create a password" 
              value={password} 
              onChange={e=>setPassword(e.target.value)} 
              className="w-full px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition" 
            />
          </div>
          {err && <div className="text-red-400 text-sm text-center font-medium bg-red-500/10 border border-red-500/30 rounded-lg py-2">{err}</div>}
          {ok && <div className="text-green-400 text-sm text-center font-semibold bg-green-500/10 border border-green-500/30 rounded-lg py-2">✓ Registered! Redirecting…</div>}
          <button 
            type="submit" 
            className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 w-full text-white font-bold py-3 rounded-xl transition shadow-lg hover:shadow-blue-500/50"
          >
            Create Account
          </button>
          <div className="mt-2 text-center text-sm text-gray-400">
            Already have an account? <Link to="/login" className="text-blue-400 hover:text-blue-300 font-semibold">Sign in</Link>
          </div>
        </form>
      </div>
    </div>
  )
}
