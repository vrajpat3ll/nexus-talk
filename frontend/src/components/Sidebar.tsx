import { Link, useLocation } from "react-router-dom"

const nav = [
  { label: "Home", icon: (
    <svg viewBox="0 0 24 24" fill="none" className="w-5 h-5"><path d="M3 12l8-8 8 8" stroke="currentColor" strokeWidth="1.5"/><path d="M5 10v10h14V10" stroke="currentColor" strokeWidth="1.5"/></svg>
  ), path: "/inbox" },
  { label: "Groups", icon: (
    <svg viewBox="0 0 24 24" fill="none" className="w-5 h-5"><circle cx="7.5" cy="8.5" r="2.5" stroke="currentColor" strokeWidth="1.5"/><circle cx="16.5" cy="8.5" r="2.5" stroke="currentColor" strokeWidth="1.5"/><rect x="2" y="14" width="20" height="7" rx="3.5" stroke="currentColor" strokeWidth="1.5"/></svg>
  ), path: "/groups" },
  { label: "Faves", icon: (
    <svg viewBox="0 0 24 24" fill="none" className="w-5 h-5"><path d="M12 18l-6.16 3.24 1.18-6.88L2 9.76l6.92-1.01L12 2.5l3.08 6.25L22 9.76l-4.98 4.6 1.18 6.88z" stroke="currentColor" strokeWidth="1.4"/></svg>
  ), path: "/faves" },
  { label: "Events", icon: (
    <svg viewBox="0 0 24 24" fill="none" className="w-5 h-5"><rect x="3" y="5" width="18" height="16" rx="3" stroke="currentColor" strokeWidth="1.5"/><path d="M16 3v4M8 3v4" stroke="currentColor" strokeWidth="1.5"/></svg>
  ), path: "/events" }
]

export default function Sidebar() {
  const loc = useLocation()
  return (
    <aside className="hidden lg:flex flex-col justify-between h-screen bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 w-[80px] py-6 px-2 fixed">
      <div>
        {/* Logo */}
        <div className="flex justify-center items-center mb-8">
          <span className="rounded-lg bg-blue-100 dark:bg-blue-900 w-11 h-11 flex items-center justify-center text-2xl shadow-md">💬</span>
        </div>
        <nav className="flex flex-col gap-5">
          {nav.map(item => (
            <Link key={item.label} to={item.path} title={item.label} className={`flex flex-col items-center py-3 px-2 rounded-lg hover:bg-blue-50 dark:hover:bg-blue-800 transition ${loc.pathname.startsWith(item.path) ? 'bg-blue-100 dark:bg-blue-900' : ''}`}>
              <span>{item.icon}</span>
              <span className="text-xs mt-1 text-center text-gray-800 dark:text-gray-200 font-medium">{item.label}</span>
            </Link>
          ))}
        </nav>
        <div className="flex items-center justify-center mt-8">
          <button className="bg-green-600 hover:bg-green-700 text-white rounded-full px-4 py-2 font-bold shadow">Post</button>
        </div>
      </div>
      <div className="flex flex-col gap-1 items-center text-xs text-gray-500 dark:text-gray-300 mb-4">
        <Link to="/settings" className="hover:underline">Settings</Link>
        <Link to="/help" className="hover:underline">Help</Link>
        <a href="#invite" className="hover:underline">Invite</a>
        <div className="mt-3 mb-1 flex items-center justify-center">
          <span className="rounded-full w-9 h-9 bg-blue-200 dark:bg-blue-800 flex items-center justify-center shadow-md">U</span>
        </div>
      </div>
    </aside>
  )
}
