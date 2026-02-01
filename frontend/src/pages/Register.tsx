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
      const resp = await fetch("http://localhost:8081/register", {
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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-blue-200 dark:from-gray-900 dark:to-gray-800 px-4">
      <form className="bg-white dark:bg-gray-900 shadow-2xl rounded-xl p-8 w-full max-w-sm flex flex-col gap-4" onSubmit={handleReg}>
        <h2 className="text-center text-2xl font-bold text-gray-800 dark:text-white">Register</h2>
        <div>
          <input name="username" placeholder="Username" value={username} onChange={e=>setUsername(e.target.value)} className="w-full px-3 py-2 rounded-lg border mt-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-blue-300" />
        </div>
        <div>
          <input type="email" name="email" placeholder="Email" value={email} onChange={e=>setEmail(e.target.value)} className="w-full px-3 py-2 rounded-lg border mt-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-blue-300" />
        </div>
        <div>
          <input type="password" name="password" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} className="w-full px-3 py-2 rounded-lg border mt-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-blue-300" />
        </div>
        {err && <div className="text-red-500 text-sm text-center font-medium">{err}</div>}
        {ok && <div className="text-green-600 text-sm text-center font-semibold">Registered! Redirecting…</div>}
        <button type="submit" className="block bg-blue-600 hover:bg-blue-700 w-full text-white font-bold py-2 mt-1 rounded-lg transition">Register</button>
        <div className="mt-3 text-center text-sm text-gray-500 dark:text-gray-400">
          Already have an account? <Link to="/login" className="underline hover:text-blue-700 font-semibold">Login</Link>
        </div>
      </form>
    </div>
  )
}
