
○ An exception list of user_ids for the "My Contacts Except..." option.
○ A unique username (string).
● Outputs:
○ An updated record in the user's privacy settings on the server.
○ When other users view the profile, they see either the user's data or a default
placeholder, according to the rules.
Non-Functional Requirements :
How well it needs to do it:
● Reliability: These privacy rules must be strictly enforced. If a user sets their picture to
"Nobody," there can be no bugs that would ever allow a stranger to see it.
3.3 Feature: Disappearing Media Controls
Functional Requirements :
● 3.3.1 Purpose of the Function: "As a privacy-conscious user, I want to be notified if
someone screenshots my disappearing photo, and I want to prevent others from
downloading the media I send."
● 3.3.2 Inputs and Outputs:
1. Inputs:
■ An OS-level event indicating a screenshot was taken.
■ The sender's choice to "Block Downloads" (true/false) on a "View
Once" message.
2. Outputs:
■ A system notification message delivered to the sender: "[Recipient's
Name] took a screenshot."
■ The absence of a "Save to gallery" button in the recipient's media viewer
UI.
Non-Functional Requirements :
How well it needs to do it:
● Responsiveness: The screenshot notification must be sent instantly.
● Security: The "Block Downloads" feature must effectively prevent the recipient from
using the app's interface to save the file to their device.
3.4 Feature: Chat Management & Organization
Functional Requirements :
● 3.4.1 Purpose of the Function: "As a user, I want to mute, pin, archive, and star
messages so that I can keep my chat list organized and easily find important things."
● 3.4.2 Inputs and Outputs:
1. Inputs:
■ User commands initiated by long-pressing a chat or message: pinChat,
muteChat, archiveChat, starMessage.
2. Outputs:
■ An immediate visual update to the chat list or message UI.
■ An updated record in the local metadata database.
Non-Functional Requirements :
How well it needs to do it:
● Usability: All these actions should be easy to find, fast to use, and simple to undo.
3.5 Feature: Advanced Message & Presence Control
Functional Requirements :
3.5.1 Purpose of the Function: "As a user, I want to delete messages for everyone and go
into 'stealth mode' so that I have complete control over my conversations and online
presence."
3.5.2 Inputs and Outputs:
● Inputs:
○ User command deleteMessage with scope: 'everyone'.
○ An on/off toggle for "Stealth Mode".
○ A selected date for the "Clear History" function.
● Outputs:
○ A recall_message command broadcast to relevant clients.
○ An updated presence visibility setting on the server.
○ Messages before a certain date are deleted from the local device's storage.
Non-Functional Requirements :
How well it needs to do it:
● Reliability: "Delete for Everyone" must reliably and quickly remove the message from all
recipient devices. When "Stealth Mode" is on, the user's "online" or "last seen" status
must never be shown to anyone.
Module 4: Privacy and User Control
Feature: Voice-to-Text Transcription
Functional Requirements :
● 4.1.1 Purpose of the Function: "As an everyday user, I want to turn voice notes into
text so that I can read a message in a quiet place (like a library or meeting) without
listening to it out loud."
● 4.1.2 Inputs and Outputs:
○ Inputs:
■ User command to transcribeMessage.
■ The audio file of the voice message (e.g., in Opus or AAC format).
■ The message_id of the voice note.
○ Outputs:
■ A text string containing the transcription of the voice message.
■ A state change in the UI, displaying the text below the corresponding
voice message player.
Non-Functional Requirements :
How well it needs to do it:
● Accuracy: The transcription must be highly accurate, correctly identifying at least 95%
of the words in a clear recording.
● Performance: The text transcription should appear quickly, ideally within a few seconds
of the request.
Feature: Smart Reply Suggestions
Functional Requirements :
● 4.2.1 Purpose of the Function: "As a busy user, I want to see smart reply suggestions
so that I can respond to common messages quickly with a single tap."
● 4.2.2 Inputs and Outputs:
1. Inputs:
■ An event trigger: newMessageReceived.
■ The text content of the newly received message and the preceding 1-2
messages for context (an array of strings).
2. Outputs:
■ An array of 3 short strings, representing the suggested replies.
■ Three tappable buttons displayed in the UI above the text input area.
Non-Functional Requirements :
How well it needs to do it:
● Relevance: The suggestions must be contextually appropriate for the conversation.
● Performance: The reply suggestions must appear almost instantly after a new message
is received.
Feature: Smart Media Search
Functional Requirements :
● 4.3.1 Purpose of the Function: "As a user, I want to search my media using keywords
(like 'dog', 'beach', or text from a screenshot) so that I can find a specific photo or
video quickly."
● 4.3.2 Inputs and Outputs:
1. Inputs:
■ An image or video file.
■ A user's search query (string).
2. Outputs:
■ A set of keyword tags associated with a media file in the local database.
■ A grid of media thumbnails that match the user's search query.
Non-Functional Requirements :
How well it needs to do it:
● Accuracy: The object and text recognition must be accurate enough to provide useful
results.
● Privacy: This scanning and indexing process must happen on the user's device only, not
on the company's servers, to protect user privacy.
● Performance: Search results should appear within a few seconds.
Feature: AI Meeting Summaries
Functional Requirements :
● 4.4.1 Purpose of the Function: "As a professional, I want to get an AI-generated
summary of recorded group calls so that I can quickly catch up on meetings I missed
or review key decisions."
● 4.4.2 Inputs and Outputs:
1. Inputs:
■ User command to generateSummary.
■ A multi-channel audio file from a recorded group call.
■ Consent from all participants to record and process the call.
2. Outputs:
■ A structured text summary, formatted with headings like "Key Topics,"
"Decisions," and "Action Items."
Non-Functional Requirements :
How well it needs to do it:
● Accuracy: The summary must accurately reflect the main points of the conversation.
● Clarity: The generated summary must be well-structured and easy to read.
Feature Module 5: Moderation and Reporting
Feature: Reporting & Context
Functional Requirements :
● 5.1.1 Purpose of the Function: "As a user, I want to report specific content or users
with clear reasons and include recent messages so that the moderation team has the
full context to make a fair decision."
● 5.1.3 Inputs and Outputs:
1. Inputs:
■ User command to createReport.
■ target_id: The unique identifier of the message, user, or group being
reported.
■ reason_code: An enum representing the user's selected reason (e.g.,
spam, harassment, hate_speech).
■ include_context: A boolean (true/false) indicating if the user
consented to include recent messages.
2. Outputs:
■ A new ticket created in the moderation system's database.
■ A success message displayed to the user: "Report submitted. Thank you
for helping keep our community safe."
Non-Functional Requirements :
How well it needs to do it:
● Usability: The reporting process should be quick and easy to find, so users are
encouraged to report bad behavior.
● Privacy: The reported messages sent for context must be handled securely and only be
visible to the authorized moderation team.
Feature: User Blocking
Functional Requirements :
5.2.1 Purpose of the Function: "As a user, I want to block another person so that I can
immediately stop them from contacting me or seeing my profile."
5.2.2 Inputs and Outputs:
● Inputs:
○ User command to blockUser or unblockUser.
○ The user_id of the user to be blocked/unblocked.
● Outputs:
○ A new row inserted into (or removed from) the block_list table.
○ An immediate cessation of all interactions between the two users.
○ The blocked user is removed from the blocker's contact list and vice versa.
Non-Functional Requirements :
How well it needs to do it:
● Reliability: The block must be 100% effective and take effect instantly. There can be no
bugs that allow a blocked user's messages to get through.
Feature: Automated Content Filtering
Functional Requirements :
● 5.3.1 Purpose of the Function: "As a user, I want the platform to automatically filter
spam, scams, and hate speech so that my environment is safer and free from harmful
content."
● 5.3.2 Inputs and Outputs:
1. Inputs:
■ Content (text or image file) posted to a public channel.
2. Outputs:
■ A moderation status for the content: approved, quarantined, or
rejected.
■ Quarantined content is hidden from public view but flagged for human
review.
Non-Functional Requirements :
How well it needs to do it:
● Accuracy: The filter should be very good at catching genuinely harmful content while
having a very low rate of incorrectly flagging normal, safe content (false positives).
● Performance: This filtering process must happen in the background and not slow down
the delivery of messages for the user.
Feature: Parental Controls
Functional Requirements :
● 5.4.1 Purpose of the Function: "As a parent, I want to enable parental controls for my
child's account so that I can restrict their exposure to inappropriate content and
ensure their safety."
● 5.4.2 Inputs and Outputs:
1. Inputs:
■ A link request from a parent account to a child account.
■ An explicit acceptance of the link from the child account.
■ The parent's desired safety settings (e.g., contact_policy:
'contacts_only').
2. Outputs:
■ A secure, consent-based link between the parent and child accounts in
the database.
■ The child's account settings are updated to reflect the parent's choices.
Non-Functional Requirements :
How well it needs to do it:
● Security: The process for linking accounts must be secure to prevent unauthorized
access. The child should not be able to easily disable the controls set by the parent.
Feature: Automated Sanctions
Functional Requirements :
● 5.5.1 Purpose of the Function: "As a community member, I want the platform to
automatically issue temporary bans to users who repeatedly break the rules so that
offenders are dealt with quickly and the community remains safe."
● 5.5.2 Inputs and Outputs:
1. Inputs:
■ A add_strike command issued by an authorized human moderator to
a specific user_id.
2. Outputs:
■ An updated account status for the user (e.g., active, suspended,
banned).
■ A suspension_end_time timestamp set for temporary bans.
■ A system-generated notification sent to the user explaining the
sanction, its duration, and the reason.
Non-Functional Requirements :
How well it needs to do it:
● Fairness & Reliability: The automated system must be highly reliable, only issuing bans
based on confirmed violations to avoid unfairly punishing innocent users.
3.6 Core Foundational Features
● Registration and Login (User Accounts)
● Basic Messaging and Rich Media
● File Sharing
● Voice and Video Calls
● Group Chats
● End-to-End Encryption
These are the features without which the app cannot fulfill its primary purpose of enabling
communication between users.
While other functionalities like channels, smart replies, and customization enhance the user
experience, the app would be non-functional and unable to launch without the core
capabilities of sending messages, sharing files, making calls, forming groups, and ensuring
user privacy through encryption.
4. External Interface Requirements
4.1 User Interfaces
The application shall present a clean, intuitive, and consistent user interface across all
platforms. It must be fully responsive and support both light and dark themes, adhering to
platform-specific human interface guidelines.
4.2 Hardware Interfaces
The application will interface with the device's camera, microphone, local storage, and
GPS/location services.
4.3 Software Interfaces
The application will interface with:
● OS APIs: For contact list integration, push notifications (APNS, FCM), and native file
pickers.
● Payment Gateway APIs: For handling channel subscriptions.
● Cloud STT and LLM APIs: For transcription and summarization features.
5. Non-Functional Requirements
5.1 Performance
● Message Latency: 99% of text messages must be delivered in under 1.5 seconds under
nominal network conditions.
● App Launch Time: The application must achieve a cold start to an interactive state in
under 3 seconds.
● Call Quality: Voice and video calls must maintain a high Mean Opinion Score (MOS) of
>4.0 on a stable broadband connection.
5.2 Security
● The implementation of the Signal Protocol shall be subject to a third-party security audit
prior to launch.
● The system must be protected against common vulnerabilities, including MITM attacks,
SQL injection, and XSS.
● All APIs must require authentication and authorization.
5.3 Reliability
● The service shall maintain a minimum uptime of 99.95%.
● The system must guarantee zero message loss during transit.
● The client application must handle network interruptions gracefully and reconnect
automatically.
5.4 Usability
● The new user registration and onboarding process must be completed in under 90
seconds.
● The application must comply with WCAG 2.1 AA accessibility standards.
5.5 Maintainability
● The codebase must be modular, well-documented, and adhere to established coding
standards for each platform.
● The system must have comprehensive logging, monitoring, and alerting to facilitate
proactive maintenance