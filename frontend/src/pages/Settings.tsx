import React, { useEffect, useState } from "react"
import { fetchPrivacy, setPrivacy, fetchPresence, setPresence } from "../api"

const userId = "demo-user" // replace with auth user

export default function Settings() {
  const [privacy, setPriv] = useState<{[k:string]:boolean}>(null as any)
  const [presence, setPres] = useState<{online:boolean,last_seen:string}|null>(null)
  const [loading, setLoading] = useState(true)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    async function load() {
      setPrivacy(await fetchPrivacy(userId))
      setPres(await fetchPresence(userId))
      setLoading(false)
    }
    load()
  }, [])

  function setPrivacy(p: any) { setPriv(p); setSaved(false) }
  async function save() {
    await setPrivacyApi(userId, privacy)
    setSaved(true)
  }
  async function toggleOnline(val: boolean) {
    await setPresence(userId, val)
    setPres(await fetchPresence(userId))
  }

  async function setPrivacyApi(uid: string, p: any) { await setPrivacy(uid, p) }

  if(loading) return <div className="p-12 text-center text-lg text-gray-400">Loading settings…</div>
  return (
    <div className="max-w-xl mx-auto py-8 px-4 bg-white dark:bg-gray-900 rounded-2xl shadow mt-8">
      <h2 className="text-2xl font-bold mb-4">Settings & Privacy</h2>
      <div className="space-y-6">
        <section className="flex items-center gap-6">
          <div>
            <label className="font-semibold text-gray-700 dark:text-white">Online Status:</label><br/>
            <span className={`inline-block w-3 h-3 rounded-full mr-2 ${presence?.online?'bg-green-500':'bg-gray-400'}`}></span>
            <span className="text-gray-500 text-sm">
              {presence?.online ? "Online" : `Last seen: ${presence?.last_seen?.slice(0,19).replace('T',' ')}`}
            </span>
          </div>
          <button className="px-4 py-1 text-sm font-bold rounded-lg shadow border border-blue-200 dark:border-blue-800 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200 hover:bg-blue-200" onClick={()=>toggleOnline(!presence?.online)}>{presence?.online?'Go Offline':'Go Online'}</button>
        </section>
        <hr className="my-4"/>
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span>Profile Visible</span>
            <label className="flex items-center cursor-pointer">
              <input type="checkbox" checked={privacy?.profile_visible} onChange={e=>setPrivacy({...privacy, profile_visible:e.target.checked})} className="w-5 h-5 text-blue-600 rounded-lg border-gray-200 focus:ring-blue-500" />
            </label>
          </div>
          <div className="flex items-center justify-between">
            <span>Show Last Seen</span>
            <label className="flex items-center cursor-pointer">
              <input type="checkbox" checked={privacy?.last_seen_visible} onChange={e=>setPrivacy({...privacy, last_seen_visible:e.target.checked})} className="w-5 h-5 text-blue-600 rounded-lg border-gray-200 focus:ring-blue-500" />
            </label>
          </div>
          <div className="flex items-center justify-between">
            <span>Read Receipts</span>
            <label className="flex items-center cursor-pointer">
              <input type="checkbox" checked={privacy?.read_receipts} onChange={e=>setPrivacy({...privacy, read_receipts:e.target.checked})} className="w-5 h-5 text-blue-600 rounded-lg border-gray-200 focus:ring-blue-500" />
            </label>
          </div>
        </div>
        <button className="mt-4 w-full py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 font-bold text-lg transition" onClick={save}>Save</button>
        {saved&&<div className="text-green-600 text-sm text-center mt-1">Saved!</div>}
      </div>
    </div>
  )
}
