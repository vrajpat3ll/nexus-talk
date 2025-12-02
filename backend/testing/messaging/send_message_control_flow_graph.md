# Control Flow Graph for messagesHandler (Send Message)

```mermaid
flowchart TD
     %% Node definitions with numbers
     N1([Entry]) --> N2{Is POST?}
     N2 -- No --> N3[Return 405]
     N2 -- Yes --> N4[Decode JSON]
     N4 -- Error --> N5[Return 400]
     N4 -- Ok --> N6{ThreadID Provided?}
     N6 -- Yes --> N7[Create Message]
     N6 -- No --> N8{ToID Provided?}
     N8 -- No --> N9[Return 400]
     N8 -- Yes --> N10[Find/Create Direct Thread]
     N10 -- Error --> N11[Return 500]
     N10 -- Ok --> N7
     N7 --> N12{useMemory?}
     N12 -- Yes --> N13[Store in Memory]
     N13 --> N14[Respond 200]
     N12 -- No --> N15[Persist to DB]
     N15 -- Error --> N16[Return 500]
     N15 -- Ok --> N17[Respond 200]
     N14 --> N18([Exit])
     N17 --> N18
     N3 --> N18
     N5 --> N18
     N9 --> N18
     N11 --> N18
     N16 --> N18

     style N1 fill:#90EE90
     style N18 fill:#FFB6C1
     style N3 fill:#FFD700
     style N5 fill:#FFD700
     style N9 fill:#FFD700
     style N11 fill:#FFA500
     style N16 fill:#FFA500
     style N14 fill:#87CEEB
     style N17 fill:#87CEEB
     style N12 fill:#E6E6FA
     style N2 fill:#E6E6FA
     style N4 fill:#B0E0E6
     style N6 fill:#E6E6FA
     style N8 fill:#E6E6FA
     style N10 fill:#B0E0E6
     style N7 fill:#B0E0E6
     style N13 fill:#B0E0E6
     style N15 fill:#B0E0E6
```

## Node Numbering
| Node | Description                      |
| ---- | -------------------------------- |
| N1   | Entry                            |
| N2   | Is POST?                         |
| N3   | Return 405                       |
| N4   | Decode JSON                      |
| N5   | Return 400 (decode error)        |
| N6   | ThreadID Provided?               |
| N7   | Create Message                   |
| N8   | ToID Provided?                   |
| N9   | Return 400 (missing to_id)       |
| N10  | Find/Create Direct Thread        |
| N11  | Return 500 (direct thread error) |
| N12  | useMemory?                       |
| N13  | Store in Memory                  |
| N14  | Respond 200 (memory)             |
| N15  | Persist to DB                    |
| N16  | Return 500 (DB error)            |
| N17  | Respond 200 (DB)                 |
| N18  | Exit                             |

## Prime Paths (Node Numbers)
- N1 -> N2 -> N3 -> N18
- N1 -> N2 -> N4 -> N5 -> N18
- N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N13 -> N14 -> N18
- N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N15 -> N17 -> N18
- N1 -> N2 -> N4 -> N6 -> N8 -> N9 -> N18
- N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N11 -> N18
- N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N7 -> N12 -> N13 -> N14 -> N18
- N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N7 -> N12 -> N15 -> N17 -> N18
- N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N15 -> N16 -> N18

## Test Paths (Node Numbers)
Below are typical test paths that should be covered in send_message_test.go:

1. Invalid method:
    N1 -> N2 -> N3 -> N18
2. Invalid JSON body:
    N1 -> N2 -> N4 -> N5 -> N18
3. Missing thread_id and to_id:
    N1 -> N2 -> N4 -> N6 -> N8 -> N9 -> N18
4. Direct thread creation error:
    N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N11 -> N18
5. Direct thread creation success, store in memory:
    N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N7 -> N12 -> N13 -> N14 -> N18
6. Direct thread creation success, persist to DB:
    N1 -> N2 -> N4 -> N6 -> N8 -> N10 -> N7 -> N12 -> N15 -> N17 -> N18
7. ThreadID provided, store in memory:
    N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N13 -> N14 -> N18
8. ThreadID provided, persist to DB:
    N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N15 -> N17 -> N18
9. DB error on persist:
    N1 -> N2 -> N4 -> N6 -> N7 -> N12 -> N15 -> N16 -> N18

## Edge Pair Coverage
All transitions between nodes are covered by the above prime paths and the test cases in send_message_test.go.
