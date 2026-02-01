import React from "react"
import { BrowserRouter, Routes, Route, Link, useLocation } from "react-router-dom"
import Home from "./pages/Home"
import Login from "./pages/Login"
import Register from "./pages/Register"
import Inbox from "./pages/Inbox"
import Chat from "./pages/Chat"
import Settings from "./pages/Settings"
import Sidebar from "./components/Sidebar"

function Layout({ children }: {children: React.ReactNode}) {
  const { pathname } = useLocation()
  const authPage = ["/login", "/register", "/"].includes(pathname)
  return (
    <div className="bg-gray-50 dark:bg-gray-900 min-h-screen flex">
      {!authPage && <Sidebar />}
      <div className="flex-1 min-w-0">
        {children}
      </div>
    </div>
  )
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Home/>}/>
      <Route path="/login" element={<Login/>}/>
      <Route path="/register" element={<Register/>}/>
      <Route path="/inbox" element={<Inbox />}/>
      <Route path="/chat/:threadId" element={<Chat />}/>
      <Route path="/settings" element={<Settings/>}/>
      <Route path="*" element={<Home/>}/>
    </Routes>
  )
}

function App() {
  return (
    <BrowserRouter>
      <Layout>
        <AppRoutes/>
      </Layout>
    </BrowserRouter>
  )
}
export default App
