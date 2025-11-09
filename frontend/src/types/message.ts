export interface Message {
  id: string;
  sender_id: string;
  thread_id: string;
  content: string;
  sent_at: string;
  content_type?: string;
  metadata?: Record<string, any>;
}

export interface Thread {
  thread_id: string;
  other_user_id: string;
  last_message_at: string;
}

export type MessageState = 'sent' | 'delivered' | 'read';