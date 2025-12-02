# Control Flow Graph for wsHandler (WebSocket Handler)

```mermaid
flowchart TD
    W1([Entry]) --> W2{user_id provided?}
    W2 -- No --> W3[Return 400]
    W2 -- Yes --> W4[Upgrade to WebSocket]
    W4 -- Error --> W5[Log upgrade error]
    W4 -- Ok --> W6[Add client]
    W6 --> W7[Defer remove client]
    W7 --> W8{ReadMessage loop}
    W8 -- Error --> W9{Unexpected close?}
    W9 -- Yes --> W10[Log read error]
    W9 -- No --> W11[Break loop]
    W8 -- Ok --> W12[Log received message]
    W12 --> W8
    W10 --> W11
    W11 --> W13([Exit])
    W3 --> W13
    W5 --> W13

    style W1 fill:#90EE90
    style W13 fill:#FFB6C1
    style W3 fill:#FFD700
    style W5 fill:#FFA500
    style W10 fill:#FFA500
    style W11 fill:#E6E6FA
    style W2 fill:#E6E6FA
    style W4 fill:#B0E0E6
    style W6 fill:#B0E0E6
    style W7 fill:#B0E0E6
    style W8 fill:#E6E6FA
    style W9 fill:#E6E6FA
    style W12 fill:#87CEEB
```

## Node Numbering
| Node | Description |
|------|-------------|
| W1   | Entry |
| W2   | user_id provided? |
| W3   | Return 400 |
| W4   | Upgrade to WebSocket |
| W5   | Log upgrade error |
| W6   | Add client |
| W7   | Defer remove client |
| W8   | ReadMessage loop |
| W9   | Unexpected close? |
| W10  | Log read error |
| W11  | Break loop |
| W12  | Log received message |
| W13  | Exit |

## Prime Paths (Node Numbers)
- W1 → W2 → W3 → W13
- W1 → W2 → W4 → W5 → W13
- W1 → W2 → W4 → W6 → W7 → W8 → W9 → W10 → W11 → W13
- W1 → W2 → W4 → W6 → W7 → W8 → W9 → W11 → W13
- W1 → W2 → W4 → W6 → W7 → W8 → W12 → W8 (loop)

## Test Paths (Node Numbers)
1. Missing user_id:
   W1 → W2 → W3 → W13
2. WebSocket upgrade error:
   W1 → W2 → W4 → W5 → W13
3. Successful connection, immediate close (unexpected):
   W1 → W2 → W4 → W6 → W7 → W8 → W9 → W10 → W11 → W13
4. Successful connection, normal close:
   W1 → W2 → W4 → W6 → W7 → W8 → W9 → W11 → W13
5. Successful connection, message received:
   W1 → W2 → W4 → W6 → W7 → W8 → W12 → W8 (loop)

## Edge Pair Coverage
All transitions between nodes are covered by the above prime paths and the test cases in ws_handler_test.go.
