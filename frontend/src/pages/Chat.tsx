import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { fetchThread, sendMessage, Message } from "../api";

interface UIMessage {
  id: string;
  sender: string; // "me" or "other"
  text: string;
  sentAt: string;
}

export default function Chat() {
  const { threadId } = useParams();
  const [input, setInput] = useState("");
  const [msgs, setMsgs] = useState<UIMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const userRaw = (typeof sessionStorage !== 'undefined' && sessionStorage.getItem("nt_user")) || localStorage.getItem("nt_user");
  const self = userRaw ? JSON.parse(userRaw) as {id:string, username:string} : null;
  // Attempt to derive peer user for header (prefer sessionStorage, fallback to localStorage)
  const peerRaw = threadId
    ? ((typeof sessionStorage !== 'undefined' && sessionStorage.getItem(`nt_thread_peer_${threadId}`))
        || localStorage.getItem(`nt_thread_peer_${threadId}`))
    : null;
  const peerUser = peerRaw ? JSON.parse(peerRaw) as {id:string, username:string} : null;

  useEffect(() => {
    let ignore = false;
    const load = async () => {
      if (!threadId || !self?.id) return;
      setLoading(true);
      try {
        const data: Message[] = await fetchThread(threadId);
        if (ignore) return;
        const mapped: UIMessage[] = data.filter(m => m.content !== "__handshake__").map(m => ({
          id: m.id,
          sender: self && m.sender_id === self.id ? "me" : "other",
          text: m.content,
          sentAt: m.sent_at,
        }));
        setMsgs(mapped);
      } catch(e) {
        console.error("fetchThread failed", e);
      } finally {
        setLoading(false);
      }
    };

    // Initial load of messages
    load();

    // Set up WebSocket connection
    if (self?.id) {
      const ws = new WebSocket(`ws://172.18.12.251:8082/v1/ws?user_id=${self.id}`);

      ws.onopen = () => {
        console.log('WebSocket connected');
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'new_message' && data.data.thread_id === threadId) {
            // Add new message to the list if it's for this thread
            const newMessage: UIMessage = {
              id: data.data.id,
              sender: data.data.sender_id === self.id ? "me" : "other",
              text: data.data.content,
              sentAt: data.data.sent_at,
            };
            setMsgs(msgs => [...msgs, newMessage]);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        // Attempt to reconnect after a delay if not navigating away
        if (!ignore) {
          setTimeout(() => load(), 3000);
        }
      };

      return () => {
        ignore = true;
        ws.close();
      };
    }
  }, [threadId, self?.id]);

  const send = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || !self) return;
    try {
      const API_HOST = import.meta.env.VITE_API_HOST || "localhost";
      // For direct thread we don't know other participant id here; assume thread already exists.
      // To support reply we need recipient id — prototype omits advanced logic.
      const temp: UIMessage = { id: "tmp"+Date.now(), sender: "me", text: input, sentAt: new Date().toISOString() };
      setMsgs([...msgs, temp]);
      setInput("");
      // Backend expects sender + to_id OR thread_id; we have threadId
      await fetch(`http://${API_HOST}:8082/v1/messages`, {
        method: "POST",
        headers: {"Content-Type":"application/json"},
        body: JSON.stringify({ sender_id: self.id, thread_id: threadId, content: temp.text })
      });
    } catch(err) {
      console.error("send failed", err);
    }
  };

  return (
    <div className="flex flex-col h-full bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      {/* Modern Chat Header */}
      <header className="flex items-center border-b border-slate-700/50 bg-slate-800/50 backdrop-blur-xl px-6 py-4 shadow-lg">
        <div className="w-11 h-11 rounded-full mr-4 flex items-center justify-center bg-gradient-to-br from-blue-500 to-purple-600 text-white font-bold text-lg shadow-lg ring-2 ring-blue-400/30">
          {peerUser ? peerUser.username.charAt(0).toUpperCase() : '?'}
        </div>
        <div className="flex-1">
          <div className="text-lg font-bold text-white">
            {peerUser ? peerUser.username : 'Direct Message'}
          </div>
          <div className="flex items-center gap-2 text-xs text-gray-400">
            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
            <span>{peerUser ? 'Active now' : 'Waiting...'}</span>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <button className="p-2 hover:bg-slate-700/50 rounded-full transition-colors">
            <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
            </svg>
          </button>
          <button className="p-2 hover:bg-slate-700/50 rounded-full transition-colors">
            <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
            </svg>
          </button>
        </div>
      </header>

      {/* Messages Area with Custom Scrollbar */}
      <div className="flex-1 overflow-y-auto px-6 py-6 space-y-4 scrollbar-thin scrollbar-thumb-slate-700 scrollbar-track-transparent">
        {loading && (
          <div className="flex items-center justify-center py-4">
            <div className="flex items-center gap-2 text-sm text-blue-400">
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce"></div>
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{animationDelay: '0.1s'}}></div>
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{animationDelay: '0.2s'}}></div>
            </div>
          </div>
        )}
        
        {msgs.length === 0 && !loading && (
          <div className="flex flex-col items-center justify-center h-full text-center py-12">
            <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-600/20 to-purple-600/20 flex items-center justify-center mb-4">
              <svg className="w-10 h-10 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <p className="text-gray-400 text-lg font-semibold">No messages yet</p>
            <p className="text-gray-500 text-sm mt-1">Send a message to start the conversation</p>
          </div>
        )}

        {msgs.map((m, idx) => (
          <div
            key={idx}
            className={`flex items-end gap-2 ${
              m.sender === "me" ? "justify-end" : "justify-start"
            }`}
          >
            {m.sender !== "me" && (
              <div className="w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center text-white text-xs font-bold shadow-lg mb-1">
                {peerUser ? peerUser.username.charAt(0).toUpperCase() : '?'}
              </div>
            )}
            
            <div className={`flex flex-col ${m.sender === "me" ? "items-end" : "items-start"} max-w-md`}>
              <div
                className={`rounded-2xl px-5 py-3 shadow-lg backdrop-blur-sm ${
                  m.sender === "me"
                    ? "bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-br-md"
                    : "bg-slate-700/70 text-white rounded-bl-md"
                }`}
              >
                <p className="text-base leading-relaxed break-words">{m.text}</p>
              </div>
              <span className="text-xs text-gray-500 mt-1 px-1">
                {new Date(m.sentAt).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}
              </span>
            </div>

            {m.sender === "me" && (
              <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center text-white text-xs font-bold shadow-lg mb-1">
                {self?.username.charAt(0).toUpperCase()}
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Modern Input Area */}
      <form
        onSubmit={send}
        className="border-t border-slate-700/50 bg-slate-800/50 backdrop-blur-xl px-6 py-4 shadow-2xl"
      >
        <div className="flex items-center gap-3 bg-slate-700/50 rounded-2xl px-4 py-2 focus-within:ring-2 focus-within:ring-blue-500/50 transition-all">
          <button
            type="button"
            className="text-gray-400 hover:text-blue-400 transition-colors p-1"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
          
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            className="flex-1 bg-transparent text-white placeholder-gray-400 outline-none text-base py-2"
          />
          
          <button
            type="button"
            className="text-gray-400 hover:text-blue-400 transition-colors p-1"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
            </svg>
          </button>

          <button
            type="submit"
            disabled={!input.trim()}
            className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed text-white rounded-xl p-2.5 shadow-lg hover:shadow-blue-500/50 transition-all duration-200"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
            </svg>
          </button>
        </div>
      </form>
    </div>
  );
}
