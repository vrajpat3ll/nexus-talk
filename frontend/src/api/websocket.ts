import { Message } from '../types/message';

export interface WebSocketMessage {
  type: 'new_message' | 'message_delivered' | 'message_read';
  data: Message;
}

interface WebSocketOptions {
  onMessage: (message: Message) => void;
  onConnectionChange?: (isConnected: boolean) => void;
}

export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private userId: string;
  private options: WebSocketOptions;

  constructor(userId: string, options: WebSocketOptions) {
    this.userId = userId;
    this.options = options;
    this.connect();
  }

  private connect() {
    if (this.ws) {
      this.ws.close();
    }

    this.ws = new WebSocket(`ws://localhost:8082/v1/ws?user_id=${this.userId}`);
    
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.options.onConnectionChange?.(true);
      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WebSocketMessage;
        if (message.type === 'new_message') {
          this.options.onMessage(message.data);
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.options.onConnectionChange?.(false);
      // Attempt to reconnect after a delay
      this.reconnectTimer = setTimeout(() => this.connect(), 5000);
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  public disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }
}

export class MessagingService {
  private wsService: WebSocketService | null = null;

  constructor(private baseUrl: string = 'http://localhost:8082') {}

  public connect(userId: string, onMessage: (message: Message) => void) {
    this.wsService = new WebSocketService(userId, {
      onMessage,
      onConnectionChange: (isConnected) => {
        console.log('WebSocket connection status:', isConnected);
      }
    });
  }

  public disconnect() {
    this.wsService?.disconnect();
    this.wsService = null;
  }

  public async sendMessage(threadId: string, content: string, senderId: string): Promise<Message> {
    try {
      const response = await fetch(`${this.baseUrl}/v1/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          sender_id: senderId,
          thread_id: threadId,
          content: content,
        }),
      });
      
      if (!response.ok) {
        throw new Error(`Failed to send message: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to send message:', error);
      throw error;
    }
  }

  public async getMessages(threadId: string, since?: Date): Promise<Message[]> {
    try {
      const url = new URL(`${this.baseUrl}/v1/messages/${threadId}`);
      if (since) {
        url.searchParams.set('since', since.toISOString());
      }

      const response = await fetch(url.toString());
      if (!response.ok) {
        throw new Error(`Failed to fetch messages: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch messages:', error);
      throw error;
    }
  }
}