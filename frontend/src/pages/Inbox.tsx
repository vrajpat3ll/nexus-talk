import React, { useEffect, useState } from 'react'
import { Link, Outlet, useLocation, useParams, useNavigate } from "react-router-dom"
import Sidebar from "../components/Sidebar"
import { fetchOrCreateDirectThread, listThreads, ThreadSummary } from "../api"

export default function Inbox() {
  const { threadId } = useParams()
  const loc = useLocation()
  const nav = useNavigate()
  const userRaw = (typeof sessionStorage !== 'undefined' && sessionStorage.getItem("nt_user")) || localStorage.getItem("nt_user")
  const self = userRaw ? JSON.parse(userRaw) as {id:string, username:string} : null

  const [threads, setThreads] = useState<ThreadSummary[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let ignore = false
    const load = async () => {
      if (!self) return
      setLoading(true)
      try {
        const ts = await listThreads(self.id)
        if (!ignore) setThreads(ts)
      } catch(e) { console.error(e) } finally { setLoading(false) }
    }
    load()
    const int = setInterval(load, 5000)
    return () => { ignore = true; clearInterval(int) }
  }, [self?.id])

  const newDM = async () => {
    const other = prompt("Enter other user's username")
    if (!other) return
    try {
      const r = await fetch(`http://172.18.12.251:8081/profile/by-username?u=${encodeURIComponent(other)}`)
      if (!r.ok) { alert("User not found"); return }
      const ou = await r.json() as {id:string, username:string}
      const tid = await fetchOrCreateDirectThread(self!.id, ou.id)
      // Store peer hint so Chat header can show a name without extra API
      try {
        sessionStorage.setItem(`nt_thread_peer_${tid}`, JSON.stringify(ou))
      } catch {}
      try {
        // Also store in localStorage as a fallback across reloads
        localStorage.setItem(`nt_thread_peer_${tid}`, JSON.stringify(ou))
      } catch {}
      nav(`/chat/${tid}`)
    } catch(e) {
      alert("Failed to start DM: "+(e as any).message)
    }
  }
  return (
    <div className="flex bg-gradient-to-br from-gray-900 via-slate-900 to-gray-900 h-screen overflow-hidden">
      <Sidebar />
      {/* Chat List Panel - REDESIGNED */}
      <section className="flex flex-col w-80 h-full bg-slate-800/50 backdrop-blur-xl border-r border-slate-700/50 shadow-2xl">
        {/* Header */}
        <header className="flex items-center justify-between py-5 px-5 border-b border-slate-700/50 bg-gradient-to-r from-blue-600/10 to-purple-600/10">
          <div>
            <h1 className="font-bold text-white text-xl tracking-tight">Messages</h1>
            <p className="text-xs text-gray-400 mt-0.5">{threads.length} conversation{threads.length !== 1 ? 's' : ''}</p>
          </div>
          <button 
            onClick={newDM} 
            className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white text-sm font-semibold rounded-lg px-4 py-2 shadow-lg hover:shadow-blue-500/50 transition-all duration-200 flex items-center gap-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            New
          </button>
        </header>

        {/* Filter Tabs */}
        <div className="flex gap-2 px-5 py-3 border-b border-slate-700/30">
          <button className="text-xs px-4 py-2 rounded-full font-semibold bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-lg">All</button>
          <button className="text-xs px-4 py-2 rounded-full font-semibold bg-slate-700/50 text-gray-300 hover:bg-slate-700 transition-colors">Unread</button>
          <button className="text-xs px-4 py-2 rounded-full font-semibold bg-slate-700/50 text-gray-300 hover:bg-slate-700 transition-colors">Groups</button>
        </div>

        {/* Status Bar */}
        <div className="px-5 py-3 border-b border-slate-700/30">
          {loading && (
            <div className="flex items-center gap-2 text-xs text-blue-400">
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
              <span>Syncing messages...</span>
            </div>
          )}
          {!loading && threads.length === 0 && (
            <div className="text-center py-8">
              <div className="w-16 h-16 mx-auto mb-3 rounded-full bg-gradient-to-br from-blue-600/20 to-purple-600/20 flex items-center justify-center">
                <svg className="w-8 h-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
              <p className="text-sm text-gray-400">No conversations yet</p>
              <p className="text-xs text-gray-500 mt-1">Start chatting with someone!</p>
            </div>
          )}
        </div>

        {/* Thread List - Scrollable */}
        <div className="flex-1 overflow-y-auto overflow-x-hidden scrollbar-thin scrollbar-thumb-slate-700 scrollbar-track-transparent">
          {threads.map((t: ThreadSummary) => {
            const peerRaw = sessionStorage.getItem(`nt_thread_peer_${t.thread_id}`) || localStorage.getItem(`nt_thread_peer_${t.thread_id}`)
            let peerName = 'Unknown'
            if (peerRaw) { try { peerName = JSON.parse(peerRaw).username } catch {} }
            const active = threadId === t.thread_id
            return (
              <button 
                key={t.thread_id} 
                onClick={()=>nav(`/chat/${t.thread_id}`)}
                className={`w-full flex items-center gap-4 px-5 py-4 text-left transition-all duration-200 border-b border-slate-700/30 ${
                  active 
                    ? 'bg-gradient-to-r from-blue-600/30 to-purple-600/30 border-l-4 border-l-blue-500 shadow-lg' 
                    : 'hover:bg-slate-700/30 border-l-4 border-l-transparent'
                }`}
              > 
                {/* Avatar with gradient */}
                <div className={`relative w-12 h-12 rounded-full flex items-center justify-center text-base font-bold shadow-lg ${
                  active 
                    ? 'bg-gradient-to-br from-blue-500 to-purple-600 text-white ring-2 ring-blue-400' 
                    : 'bg-gradient-to-br from-blue-600 to-blue-800 text-white'
                }`}>
                  {peerName.charAt(0).toUpperCase()}
                  <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-slate-800"></div>
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-1">
                    <h3 className={`truncate text-base font-semibold ${active ? 'text-white' : 'text-gray-200'}`}>
                      {peerName}
                    </h3>
                    <span className={`text-xs flex-shrink-0 ml-2 ${active ? 'text-blue-300' : 'text-gray-500'}`}>
                      {new Date(t.last_message_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
                    </span>
                  </div>
                  <p className={`truncate text-sm ${active ? 'text-blue-200' : 'text-gray-400'}`}>
                    Tap to view conversation
                  </p>
                </div>

                {/* Unread badge (placeholder) */}
                {!active && Math.random() > 0.7 && (
                  <div className="flex-shrink-0 w-5 h-5 bg-blue-600 rounded-full flex items-center justify-center">
                    <span className="text-xs text-white font-bold">2</span>
                  </div>
                )}
              </button>
            )
          })}
        </div>
      </section>

      {/* Main chat window */}
      <section className="flex-1 h-full overflow-hidden flex flex-col">
        {loc.pathname.startsWith("/chat/")
          ? <Outlet/>
          : (
            <div className="h-full flex flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
              <div className="text-center space-y-4 max-w-md px-6">
                <div className="w-24 h-24 mx-auto rounded-full bg-gradient-to-br from-blue-600/20 to-purple-600/20 flex items-center justify-center mb-4">
                  <svg className="w-12 h-12 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                </div>
                <h2 className="text-2xl font-bold text-white">Welcome to NexusTalk</h2>
                <p className="text-gray-400">Select a conversation from the sidebar or start a new one to begin messaging</p>
                <button 
                  onClick={newDM}
                  className="mt-6 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white font-semibold rounded-lg px-6 py-3 shadow-xl hover:shadow-2xl transition-all duration-200"
                >
                  Start New Conversation
                </button>
              </div>
            </div>
          )}
      </section>
    </div>
  )
}
