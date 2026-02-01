const API_HOST = "172.18.12.251";

export async function login(username: string, password: string) {
  const r = await fetch(`http://${API_HOST}:8081/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({username, password})
  });
  if (!r.ok) throw new Error('Login failed');
  return r.json() as Promise<{token:string,user:{id:string,username:string}}>;
}

export async function register(username: string, email: string, password: string) {
  const r = await fetch(`http://${API_HOST}:8081/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({username, email, password})
  });
  if (!r.ok) throw new Error('Registration failed');
  return r.json();
}

export async function getMessages(threadId: string) {
  const r = await fetch(`http://${API_HOST}:8082/messages?thread_id=${threadId}`);
  if (!r.ok) throw new Error('Fetch messages failed');
  return r.json();
}

export async function sendMessage(from_id: string, to_id: string, content: string) {
  // Use unified /v1/messages endpoint which will find/create a direct thread
  const resp = await fetch(`http://${API_HOST}:8082/v1/messages`, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({sender_id: from_id, to_id, content})
  })
  if (!resp.ok) throw new Error('Send failed')
  return resp.json() as Promise<{id:string, thread_id:string, sender_id:string, content:string, sent_at:string}>;
}

export interface Message { id:string; sender_id:string; thread_id:string; content:string; sent_at:string }

export async function fetchThread(thread_id: string): Promise<Message[]> {
  const r = await fetch(`http://${API_HOST}:8082/v1/messages/${thread_id}`)
  if (!r.ok) throw new Error('fetchThread failed')
  return r.json()
}

export async function fetchOrCreateDirectThread(selfId: string, otherId: string): Promise<string> {
  const r = await fetch(`http://${API_HOST}:8082/v1/threads/direct`, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({ user_a: selfId, user_b: otherId })
  })
  if (!r.ok) throw new Error('direct thread failed')
  const js = await r.json() as {thread_id:string}
  return js.thread_id
}

export type ThreadSummary = { thread_id: string, other_user_id: string, last_message_at: string }
export async function listThreads(userId: string): Promise<ThreadSummary[]> {
  const r = await fetch(`http://${API_HOST}:8082/v1/threads?user_id=${userId}`)
  if (!r.ok) throw new Error('listThreads failed')
  const js = await r.json() as Array<{thread_id:string, other_user_id:string, last_message_at:string}>
  return js
}

// Channel/group API
type Channel = { id: string, name: string, owner: string, members: string[] };
export async function createChannel(owner: string, name: string): Promise<Channel> {
  const r = await fetch(`http://${API_HOST}:8083/create`, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({owner, name})
  });
  if(!r.ok) throw new Error("createChannel failed")
  return await r.json();
}
export async function joinChannel(user: string, channelID: string) {
  const r = await fetch(`http://${API_HOST}:8083/join`, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({channelID, user})
  });
  if (!r.ok) throw new Error("joinChannel failed");
  return r.json();
}
export async function listChannels(user_id: string): Promise<Channel[]> {
  const r = await fetch(`http://${API_HOST}:8083/list?user_id=${user_id}`);
  if(!r.ok) throw new Error('listChannels failed');
  return await r.json();
}

// Media/file API
export async function uploadMedia(owner: string, file: File): Promise<{id:number,name:string}> {
  const form = new FormData();
  form.append('owner', owner);
  form.append('file', file);
  const r = await fetch(`http://${API_HOST}:8084/upload`, {
    method: "POST",
    body: form
  });
  if(!r.ok) throw new Error("upload failed");
  return await r.json();
}
export async function downloadMedia(id: number): Promise<Blob> {
  const r = await fetch(`http://${API_HOST}:8084/download?id=${id}`);
  if(!r.ok) throw new Error("not found");
  return await r.blob();
}

// Presence/Privacy
type Presence = { online: boolean, last_seen: string };
export async function fetchPresence(user_id: string): Promise<Presence> {
  const r = await fetch(`http://${API_HOST}:8081/presence?user_id=${user_id}`);
  if(!r.ok) throw new Error("fetchPresence failed");
  return await r.json();
}
export async function setPresence(user_id: string, online: boolean) {
  const r = await fetch(`http://${API_HOST}:8081/presence`, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({ userID: user_id, online })
  });
  if(!r.ok) throw new Error("setPresence failed");
  return await r.json();
}
export async function fetchPrivacy(user_id: string): Promise<Record<string,boolean>> {
  const r = await fetch(`http://${API_HOST}:8081/privacy?user_id=${user_id}`);
  if(!r.ok) throw new Error("fetchPrivacy failed");
  return await r.json();
}
export async function setPrivacy(user_id: string, settings: Record<string,boolean>) {
  const r = await fetch(`http://${API_HOST}:8081/privacy`, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({ userID: user_id, settings })
  });
  if(!r.ok) throw new Error("setPrivacy failed");
  return await r.json();
}
