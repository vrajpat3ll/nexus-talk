import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Inbox from '../pages/Inbox';
import Chat from '../pages/Chat';
import Channel from '../pages/Channel';
import Settings from '../pages/Settings';

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/inbox" element={<Inbox />} />
      <Route path="/chat/:threadId" element={<Chat />} />
      <Route path="/channel/:channelId" element={<Channel />} />
      <Route path="/settings" element={<Settings />} />
      <Route path="*" element={<Login />} />
    </Routes>
  );
}
