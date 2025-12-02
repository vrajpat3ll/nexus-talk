import { Link, useLocation } from "react-router-dom"

const nav = [
  { label: "Chats", icon: (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
      <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
    </svg>
  ), path: "/inbox" },
  { label: "Settings", icon: (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
      <path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94L14.4 2.81c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>
    </svg>
  ), path: "/settings" }
]

export default function Sidebar() {
  const loc = useLocation()
  const userRaw = (typeof sessionStorage !== 'undefined' && sessionStorage.getItem("nt_user")) || localStorage.getItem("nt_user")
  const user = userRaw ? JSON.parse(userRaw) as {id:string, username:string} : null

  return (
    <aside className="flex flex-col h-screen bg-slate-800 border-r border-slate-700/50 w-[72px] flex-shrink-0">
      {/* Logo */}
      <div className="flex justify-center items-center py-6">
        <div className="rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 w-12 h-12 flex items-center justify-center text-2xl shadow-lg">
          💬
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 flex flex-col gap-1 px-2">
        {nav.map(item => {
          const isActive = loc.pathname === item.path || (item.path === '/inbox' && loc.pathname.startsWith('/chat'))
          return (
            <Link 
              key={item.label} 
              to={item.path} 
              title={item.label}
              className={`group relative flex flex-col items-center justify-center py-3 rounded-xl transition-all duration-200 ${
                isActive 
                  ? 'bg-slate-700 text-blue-400' 
                  : 'text-gray-400 hover:bg-slate-700/50 hover:text-white'
              }`}
            >
              <span className={isActive ? 'scale-110' : 'group-hover:scale-105 transition-transform'}>
                {item.icon}
              </span>
              <span className="text-[10px] mt-1 font-medium">{item.label}</span>
              {isActive && (
                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-blue-500 rounded-r-full"></div>
              )}
            </Link>
          )
        })}
      </nav>

      {/* User Profile */}
      <div className="px-2 pb-4 pt-4 border-t border-slate-700/50">
        <div className="flex flex-col items-center gap-2">
          <div className="relative group cursor-pointer">
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-lg ring-2 ring-slate-700 hover:ring-blue-500 transition-all">
              {user ? user.username.charAt(0).toUpperCase() : 'U'}
            </div>
            <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-slate-800"></div>
          </div>
        </div>
      </div>
    </aside>
  )
}
