import { Link } from "react-router-dom"

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-blue-100 dark:from-gray-900 dark:to-gray-800 dark:via-neutral-900 px-6 lg:px-12 py-12 flex items-center">
      <div className="w-full max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
        {/* Hero content */}
        <section className="px-6 py-6 lg:py-0">
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold text-blue-700 dark:text-blue-300 leading-tight">NexusTalk — messaging built for people and communities</h1>
          <p className="mt-6 text-lg text-gray-600 dark:text-gray-300 max-w-xl">Fast, secure, and delightful chat with modern features — channels, groups, voice/video, and moderation tools. Open, extensible, and privacy-forward.</p>

          <div className="mt-8 flex flex-col sm:flex-row gap-4 sm:gap-6">
            <Link to="/login" className="inline-block">
              <span className="inline-flex items-center justify-center px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white text-lg font-semibold rounded-xl shadow-md transition">Get started — Login</span>
            </Link>
            <Link to="/register" className="inline-block">
              <span className="inline-flex items-center justify-center px-6 py-3 bg-white hover:bg-blue-50 text-blue-800 font-semibold rounded-xl shadow-sm border border-gray-200 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-200 transition">Create an account</span>
            </Link>
          </div>

          <ul className="mt-10 grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-md text-sm text-gray-600 dark:text-gray-300">
            <li className="flex items-start gap-3">
              <span className="rounded-md bg-blue-100 text-blue-700 p-2">🔒</span>
              <div>
                <strong className="block text-gray-900 dark:text-gray-100">Secure by default</strong>
                <span className="text-xs">End-to-end encryption ready, and strong auth options.</span>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="rounded-md bg-purple-100 text-purple-700 p-2">⚡</span>
              <div>
                <strong className="block text-gray-900 dark:text-gray-100">Realtime & fast</strong>
                <span className="text-xs">Low-latency messaging and media streaming.</span>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="rounded-md bg-green-100 text-green-700 p-2">🧭</span>
              <div>
                <strong className="block text-gray-900 dark:text-gray-100">Moderation tools</strong>
                <span className="text-xs">Built-in moderation and admin controls.</span>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="rounded-md bg-yellow-100 text-yellow-700 p-2">🔗</span>
              <div>
                <strong className="block text-gray-900 dark:text-gray-100">Extendable</strong>
                <span className="text-xs">Open APIs and microservices-first architecture.</span>
              </div>
            </li>
          </ul>
        </section>

        {/* Illustration / preview card */}
        <aside className="px-6">
          <div className="relative mx-auto max-w-md lg:max-w-none">
            <div className="hero-card bg-white dark:bg-gray-800/60 backdrop-blur-md rounded-2xl shadow-2xl p-6 sm:p-8">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Live preview</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-300">An example conversation and compact channel list.</p>
                </div>
                <div className="text-sm text-gray-400">v1.0</div>
              </div>

              <div className="mt-6 bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800 rounded-lg p-4 border border-gray-100 dark:border-gray-700">
                {/* Placeholder chat mock */}
                <div className="space-y-3">
                  <div className="flex items-start gap-3">
                    <span className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-blue-100 text-blue-700">A</span>
                    <div className="bg-gray-100 dark:bg-gray-900 rounded-lg px-3 py-2 text-sm">Hey — welcome to NexusTalk! This is a preview of the app UI.</div>
                  </div>
                  <div className="flex items-start gap-3 justify-end">
                    <div className="bg-blue-600 text-white rounded-lg px-3 py-2 text-sm">Nice — love the look 😄</div>
                    <span className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-green-100 text-green-700">Y</span>
                  </div>
                  <div className="mt-2 text-xs text-gray-400">Channels • 12 members • Active now</div>
                </div>
              </div>
            </div>
            <div className="absolute -right-8 -bottom-8 w-40 h-40 rounded-2xl bg-gradient-to-tr from-blue-200 to-purple-200 opacity-70 blur-lg transform rotate-12" aria-hidden></div>
          </div>
        </aside>
      </div>
    </main>
  )
}
